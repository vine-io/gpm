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

package server

import (
	"mime"
	gruntime "runtime"

	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	"github.com/gpm2/gpm/pkg/service"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lack-io/cli"
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
	"github.com/lack-io/vine/util/helper"
	"github.com/lack-io/vine/util/namespace"
	"github.com/rakyll/statik/fs"

	_ "github.com/lack-io/vine/lib/api/handler/openapi/statik"
)

var (
	Address       = ":7800"
	Handler       = "rpc"
	Type          = "api"
	APIPath       = "/"
	enableOpenAPI = false

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "root",
			Usage:   "gpmd root directory",
			EnvVars: []string{"GPMD_ROOT"},
		},
		&cli.StringFlag{
			Name:        "api-address",
			Usage:       "The specify for api address",
			EnvVars:     []string{"VINE_API_ADDRESS"},
			Value:       Address,
			Destination: &Address,
		},
		&cli.BoolFlag{
			Name:    "enable-openapi",
			Usage:   "Enable OpenAPI3",
			EnvVars: []string{"VINE_ENABLE_OPENAPI"},
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "enable-cors",
			Usage:   "Enable CORS, allowing the API to be called by frontend applications",
			EnvVars: []string{"VINE_API_ENABLE_CORS"},
			Value:   true,
		},
	}
)

type server struct {
	vine.Service

	api apihttp.Server

	H service.Gpm `inject:""`
}

func (s *server) Init() error {
	var err error

	// Init API
	var aopts []apihttp.Option

	opts := []vine.Option{
		vine.Name(runtime.GpmName),
		vine.Id(runtime.GpmId),
		vine.Version(runtime.GetVersion()),
		vine.Metadata(map[string]string{
			"api-address": Address,
			"namespace":   runtime.Namespace,
		}),
		vine.Flags(flags...),
		vine.Action(func(c *cli.Context) error {
			cfg := &config.Config{}
			cfg.Root = c.String("root")
			if cfg.Root == "" {
				if gruntime.GOOS == "windows" {
					cfg.Root = "C:\\opt\\gpm"
				} else {
					cfg.Root = "/opt/gpm"
				}
			}

			enableOpenAPI = c.Bool("enable-openapi")

			if c.Bool("enable-tls") {
				cfg, err := helper.TLSConfig(c)
				if err != nil {
					log.Errorf(err.Error())
					return err
				}

				aopts = append(aopts, apihttp.EnableTLS(true))
				aopts = append(aopts, apihttp.TLSConfig(cfg))
			}
			return inject.Provide(cfg)
		}),
	}

	s.Service.Init(opts...)

	aopts = append(aopts, apihttp.EnableCORS(true))

	// create the router
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	if enableOpenAPI {
		openAPI := openapi.New(s.Service)
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
		router.WithRegistry(s.Options().Registry),
	)
	rp := arpc.NewHandler(
		ahandler.WithNamespace(runtime.Namespace),
		ahandler.WithRouter(rt),
		ahandler.WithClient(s.Client()),
	)
	app.Group(APIPath, rp.Handle)

	api := httpapi.NewServer(Address)
	if err = api.Init(aopts...); err != nil {
		return err
	}
	api.Handle("/", app)
	s.api = api

	db := new(dao.DB)
	if err = inject.Provide(s.Service, s.Client(), s, db); err != nil {
		return err
	}

	if err = inject.Populate(); err != nil {
		return err
	}

	if err = s.H.Init(); err != nil {
		return err
	}

	if err = pb.RegisterGpmServiceHandler(s.Service.Server(), s); err != nil {
		return err
	}

	return err
}

func (s *server) Run() error {
	if err := s.api.Start(); err != nil {
		return err
	}

	if err := s.Service.Run(); err != nil {
		return err
	}

	if err := s.api.Stop(); err != nil {
		return err
	}

	return nil
}

func New(opts ...vine.Option) *server {
	srv := vine.NewService(opts...)
	return &server{
		Service: srv,
	}
}
