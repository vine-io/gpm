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
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	gruntime "runtime"
	"time"

	"github.com/gpm2/gpm/pkg/api"
	"github.com/gpm2/gpm/pkg/dao"
	"github.com/gpm2/gpm/pkg/runtime"
	"github.com/gpm2/gpm/pkg/runtime/config"
	"github.com/gpm2/gpm/pkg/runtime/inject"
	"github.com/gpm2/gpm/pkg/runtime/ssl"
	"github.com/gpm2/gpm/pkg/service"
	pb "github.com/gpm2/gpm/proto/service/gpm/v1"
	"github.com/lack-io/cli"
	"github.com/lack-io/plugins/logger/zap"
	"github.com/lack-io/vine"
	grpcClient "github.com/lack-io/vine/core/client/grpc"
	vserver "github.com/lack-io/vine/core/server"
	grpcServer "github.com/lack-io/vine/core/server/grpc"
	apihttp "github.com/lack-io/vine/lib/api/server"
	log "github.com/lack-io/vine/lib/logger"
	"github.com/lack-io/vine/util/helper"
	"google.golang.org/grpc/peer"
)

var (
	APIAddress    = ":7800"
	EnableOpenAPI = true
	Address       = ":7700"

	flags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "enable-log",
			Usage:   "write log to file",
			EnvVars: []string{"VINE_LOG"},
		},
		&cli.StringFlag{
			Name:        "api-address",
			Usage:       "The specify for api address",
			EnvVars:     []string{"VINE_API_ADDRESS"},
			Value:       APIAddress,
			Destination: &APIAddress,
		},
		&cli.BoolFlag{
			Name:        "enable-openapi",
			Usage:       "Enable OpenAPI3",
			EnvVars:     []string{"VINE_ENABLE_OPENAPI"},
			Destination: &EnableOpenAPI,
			Value:       EnableOpenAPI,
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

	API *api.RestAPI `inject:""`

	H service.Gpm `inject:""`
}

func (s *server) Init() error {
	var err error

	// Init API
	var aopts []apihttp.Option

	cfg := &config.Config{}
	if cfg.Root == "" {
		if gruntime.GOOS == "windows" {
			cfg.Root = "C:\\opt\\gpm"
		} else {
			cfg.Root = "/opt/gpm"
		}
		_ = os.MkdirAll(filepath.Join(cfg.Root, "logs"), 0777)
		_ = os.MkdirAll(filepath.Join(cfg.Root, "services"), 0777)
		_ = os.MkdirAll(filepath.Join(cfg.Root, "packages"), 0777)
	}

	gh, err := ssl.GetSSL(cfg.Root)
	if err != nil {
		return fmt.Errorf("load server tls: %v", err)
	}
	tls, err := ssl.GetTLS()
	if err != nil {
		return fmt.Errorf("load client tls: %v", err)
	}

	ghTLSOption := func() vine.Option { return func(o *vine.Options) { _ = o.Server.Init(grpcServer.GrpcToHttp(gh)) } }
	cliTLSOption := func() vine.Option { return func(o *vine.Options) { _ = o.Client.Init(grpcClient.AuthTLS(tls)) } }

	opts := []vine.Option{
		vine.Name(runtime.GpmName),
		vine.Id(runtime.GpmId),
		vine.Version(runtime.GetVersion()),
		vine.Address(Address),
		vine.Metadata(map[string]string{
			"api-address": APIAddress,
			"namespace":   runtime.Namespace,
		}),
		ghTLSOption(),
		cliTLSOption(),
		vine.Flags(flags...),
		vine.WrapHandler(newLoggerWrapper()),
		vine.Action(func(c *cli.Context) error {

			if c.Bool("enable-tls") {
				cfg, err := helper.TLSConfig(c)
				if err != nil {
					log.Errorf(err.Error())
					return err
				}

				aopts = append(aopts, apihttp.EnableTLS(true))
				aopts = append(aopts, apihttp.TLSConfig(cfg))
			}

			Address = c.String("server-address")
			cfg.Address = Address
			cfg.APIAddress = APIAddress
			cfg.EnableOpenAPI = EnableOpenAPI

			if c.Bool("enable-log") {
				cfg.EnableLog = true

				l, err := zap.New(zap.WithFileWriter(zap.FileWriter{
					FileName:   filepath.Join(cfg.Root, "logs", "gpmd.log"),
					MaxSize:    1,
					MaxBackups: 5,
					MaxAge:     30,
					Compress:   false,
				}))
				if err != nil {
					return err
				}
				defer l.Sync()
				log.DefaultLogger = l
			}

			return nil
		}),
	}

	s.Service.Init(opts...)

	aopts = append(aopts, apihttp.EnableCORS(true))

	if err = inject.Provide(cfg); err != nil {
		return err
	}

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

	if err = s.API.Init(aopts...); err != nil {
		return err
	}

	if err = pb.RegisterGpmServiceHandler(s.Service.Server(), s); err != nil {
		return err
	}

	return err
}

func (s *server) Run() error {
	if err := s.API.Start(); err != nil {
		return err
	}

	if err := s.Service.Run(); err != nil {
		return err
	}

	if err := s.API.Stop(); err != nil {
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

func newLoggerWrapper() vserver.HandlerWrapper {
	return func(fn vserver.HandlerFunc) vserver.HandlerFunc {
		return func(ctx context.Context, req vserver.Request, rsp interface{}) error {
			buf := bytes.NewBuffer([]byte(""))
			now := time.Now()
			err := fn(ctx, req, rsp)
			buf.WriteString(fmt.Sprintf("[%s] ", time.Now().Sub(now).String()))
			buf.WriteString("[" + req.ContentType() + "] ")
			pr, ok := peer.FromContext(ctx)
			if ok {
				buf.WriteString(pr.Addr.String() + " -> ")
			}
			buf.WriteString(req.Service() + "-" + req.Endpoint())
			if err != nil {
				buf.WriteString(": " + err.Error())
				log.Error(buf.String())
			} else {
				log.Info(buf.String())
			}
			return err
		}
	}
}
