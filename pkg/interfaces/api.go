// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package interfaces

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	pbr "github.com/schollz/progressbar/v3"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/runtime"
	"github.com/vine-io/gpm/pkg/runtime/client"
	vclient "github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/registry"
	ahandler "github.com/vine-io/vine/lib/api/handler"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	arpc "github.com/vine-io/vine/lib/api/handler/rpc"
	"github.com/vine-io/vine/lib/api/resolver"
	"github.com/vine-io/vine/lib/api/resolver/grpc"
	"github.com/vine-io/vine/lib/api/router"
	regRouter "github.com/vine-io/vine/lib/api/router/registry"
	apihttp "github.com/vine-io/vine/lib/api/server"
	httpapi "github.com/vine-io/vine/lib/api/server/http"
	"github.com/vine-io/vine/lib/config"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/namespace"

	_ "github.com/vine-io/vine/lib/api/handler/openapi/statik"
)

const (
	Handler = "rpc"
	Type    = "api"
	APIPath = "/"
)

type RestAPI struct {
	apihttp.Server

	app *fiber.App
}

func (r *RestAPI) Init(opts ...apihttp.Option) error {
	// create the router
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if config.Get("enable", "openapi").Bool(false) {
		openapi.RegisterOpenAPI(app)
	}

	// TODO: more api
	app.Get("/api/v1/endpoints", r.getEndpointsHandler())
	app.Post("/api/v1/push", r.pushHandler())

	// create the namespace resolver
	nsResolver := namespace.NewResolver(Type, runtime.Namespace)
	// resolver options
	ropts := []resolver.Option{
		resolver.WithNamespace(nsResolver.ResolveWithType),
		resolver.WithHandler(Handler),
	}

	log.Infof("Registering API RPC Handler at %s", APIPath)
	rr := grpc.NewResolver(ropts...)
	rt := regRouter.NewRouter(
		router.WithHandler(arpc.Handler),
		router.WithResolver(rr),
		router.WithRegistry(registry.DefaultRegistry),
	)
	rp := arpc.NewHandler(
		ahandler.WithNamespace(runtime.Namespace),
		ahandler.WithRouter(rt),
		ahandler.WithClient(vclient.DefaultClient),
	)

	app.Group(APIPath, rp.Handle)
	api := httpapi.NewServer(config.Get("api", "address").String(""))
	if err := api.Init(opts...); err != nil {
		return err
	}

	api.Handle("/", app)
	r.app = app
	r.Server = api

	if err := r.Server.Init(opts...); err != nil {
		return err
	}

	return nil
}

func (r *RestAPI) getEndpointsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		endpoints := make([]map[string]string, 0)
		keys := make(map[string]struct{}, 0)
		list, _ := registry.GetService(runtime.GpmName)
		for _, item := range list {
			for _, e := range item.Endpoints {
				if _, ok := keys[e.Name]; ok {
					continue
				} else {
					keys[e.Name] = struct{}{}
				}
				if v, ok := e.Metadata["stream"]; ok && v == "true" {
					continue
				}
				if _, ok := e.Metadata["path"]; !ok {
					continue
				}
				endpoints = append(endpoints, map[string]string{
					"name":        e.Name,
					"method":      e.Metadata["method"],
					"path":        e.Metadata["path"],
					"description": e.Metadata["description"],
				})
			}
		}

		return c.JSON(fiber.Map{
			"data": endpoints,
		})
	}
}

func (r *RestAPI) pushHandler() fiber.Handler {

	return func(c *fiber.Ctx) error {
		mf, err := c.MultipartForm()
		if err != nil {
			return fiber.NewError(http.StatusBadRequest, err.Error())
		}

		fs, ok := mf.File["file"]
		if !ok {
			return fiber.NewError(http.StatusBadRequest, "missing file")
		}
		vv := mf.Value["dst"]
		if !ok {
			return fiber.NewError(http.StatusBadRequest, "missing dst")
		}
		file := fs[0]
		dst := vv[0]

		fd, err := file.Open()
		if err != nil {
			return fiber.NewError(http.StatusInternalServerError, "open file: "+err.Error())
		}
		defer fd.Close()

		in := &gpmv1.PushIn{
			Name:  file.Filename,
			Dst:   dst,
			Total: file.Size,
		}

		ctx := c.Context()
		cc := client.New()
		opts := []vclient.CallOption{
			vclient.WithDialTimeout(time.Hour * 2),
			vclient.WithStreamTimeout(time.Hour * 2),
		}

		stream, err := cc.Push(ctx, opts...)
		if err != nil {
			return fiber.NewError(http.StatusBadGateway, "connect to gpm server: "+err.Error())
		}

		outE := log.DefaultLogger.Options().Out
		pb := pbr.NewOptions(int(file.Size),
			pbr.OptionSetWriter(outE),
			pbr.OptionShowBytes(true),
			pbr.OptionEnableColorCodes(true),
			pbr.OptionOnCompletion(func() {
				fmt.Fprintf(outE, "\n")
			}),
		)

		buf := make([]byte, 1024*32)
		for {
			n, e := fd.Read(buf)
			if e != nil && e != io.EOF {
				return fiber.NewError(http.StatusInternalServerError, e.Error())
			}

			if n > 0 {
				_ = pb.Add(n)
				in.Length = int64(n)
				in.Chunk = buf[0:n]
				err = stream.Send(in)
				if err != nil {
					return fiber.NewError(http.StatusInternalServerError, err.Error())
				}
			}

			if e == io.EOF {
				break
			}
		}

		if err = stream.Wait(); err != nil {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(fiber.Map{
			"result": "OK",
		})
	}
}