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

package api

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/client"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	gpmv1 "github.com/gpm2/gpm/proto/apis/gpm/v1"
	"github.com/lack-io/vine"
	ahandler "github.com/lack-io/vine/lib/api/handler"
	"github.com/lack-io/vine/lib/api/handler/openapi"
	arpc "github.com/lack-io/vine/lib/api/handler/rpc"
	"github.com/lack-io/vine/lib/api/resolver"
	"github.com/lack-io/vine/lib/api/resolver/grpc"
	"github.com/lack-io/vine/lib/api/router"
	regRouter "github.com/lack-io/vine/lib/api/router/registry"
	apihttp "github.com/lack-io/vine/lib/api/server"
	httpapi "github.com/lack-io/vine/lib/api/server/http"
	log "github.com/lack-io/vine/lib/logger"
	verrs "github.com/lack-io/vine/proto/apis/errors"
	"github.com/lack-io/vine/util/namespace"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc/status"
)

func init() {
	inject.ProvidePanic(new(RestAPI))
}

const (
	Handler = "rpc"
	Type    = "api"
	APIPath = "/"
)

type RestAPI struct {
	S vine.Service `inject:""`

	Cfg *config.Config `inject:""`

	apihttp.Server

	app *fiber.App
}

func (r *RestAPI) Init(opts ...apihttp.Option) error {
	// create the router
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if r.Cfg.EnableOpenAPI {
		openAPI := openapi.New(r.S)
		_ = mime.AddExtensionType(".svg", "image/svg+xml")
		sfs, err := fs.New()
		if err != nil {
			log.Fatalf("Starting OpenAPI: %v", err)
		}
		prefix := "/openapi-ui/"
		app.All(prefix, openAPI.OpenAPIHandler)
		app.Use(prefix, filesystem.New(filesystem.Config{Root: sfs}))
		app.Get("/openapi.json", openAPI.OpenAPIJOSNHandler)
		app.Get("/services", openAPI.OpenAPIServiceHandler)
		log.Infof("Starting OpenAPI at %v", prefix)
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
		router.WithRegistry(r.S.Options().Registry),
	)
	rp := arpc.NewHandler(
		ahandler.WithNamespace(runtime.Namespace),
		ahandler.WithRouter(rt),
		ahandler.WithClient(r.S.Client()),
	)

	app.Group(APIPath, rp.Handle)
	api := httpapi.NewServer(r.Cfg.APIAddress)
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
		list, _ := r.S.Options().Registry.GetService(runtime.GpmName)
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
		stream, err := cc.Push(ctx)
		if err != nil {
			return fiber.NewError(http.StatusBadGateway, "connect to gpm server: "+err.Error())
		}
		defer stream.Close()

		done := make(chan struct{}, 1)
		ech := make(chan error, 1)

		go func() {
			buf := make([]byte, 1024)
			for {
				n, e := fd.Read(buf)
				if e != nil && e != io.EOF {
					ech <- fmt.Errorf("read file: %v", err)
					return
				}
				if e == io.EOF {
					in.IsOk = true
				}
				in.Length = int64(n)
				in.Chunk = buf[0:n]
				e = stream.Send(in)
				if e != nil {
					ech <- fmt.Errorf("save file: %v", err)
					return
				}
				if e == io.EOF {
					break
				}
			}
		}()

		go func() {
			for {
				select {
				case <-stream.Context().Done():
					goto EXIT
				default:
				}
				b, e := stream.Recv()
				if e != nil {
					ech <- errors.New(status.Convert(e).Message())
					return
				}
				if b.Error != "" {
					ech <- errors.New(verrs.Parse(b.Error).Detail)
					return
				}
				if b.IsOk {
					break
				}
			}
		EXIT:
			done <- struct{}{}
		}()

		select {
		case e := <-ech:
			return fiber.NewError(http.StatusInternalServerError, e.Error())
		case <-done:
		}

		return c.JSON(fiber.Map{
			"result": "OK",
		})
	}
}
