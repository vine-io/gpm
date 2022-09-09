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

package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	gruntime "runtime"
	"time"

	pb "github.com/vine-io/gpm/api/service/gpm/v1"
	"github.com/vine-io/gpm/pkg/biz"
	"github.com/vine-io/gpm/pkg/infra/repo"
	"github.com/vine-io/gpm/pkg/runtime"
	"github.com/vine-io/gpm/pkg/runtime/inject"
	"github.com/vine-io/gpm/pkg/runtime/ssl"
	"github.com/vine-io/pkg/release"

	"github.com/vine-io/cli"
	"github.com/vine-io/plugins/logger/zap"
	"github.com/vine-io/vine"
	grpcClient "github.com/vine-io/vine/core/client/grpc"
	"github.com/vine-io/vine/core/registry/memory"
	vserver "github.com/vine-io/vine/core/server"
	grpcServer "github.com/vine-io/vine/core/server/grpc"
	apihttp "github.com/vine-io/vine/lib/api/server"
	"github.com/vine-io/vine/lib/config"
	"github.com/vine-io/vine/lib/config/source"
	ccli "github.com/vine-io/vine/lib/config/source/cli"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/helper"
	"google.golang.org/grpc/peer"
)

var (
	APIAddress    = ":7800"
	EnableOpenAPI = true
	EnableLog     = false
	Address       = ":7700"
	ROOT          = "/opt/gpm"

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "root",
			Usage:       "the root directory of gpmd",
			EnvVars:     []string{"VINE_ROOT"},
			Destination: &ROOT,
			Value:       ROOT,
		},
		&cli.BoolFlag{
			Name:        "enable-log",
			Usage:       "write log to file",
			EnvVars:     []string{"VINE_LOG"},
			Destination: &EnableLog,
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

type GpmAPI struct {
	vine.Service

	API *RestAPI

	G biz.Manager `inject:""`
	T biz.FTP     `inject:""`
}

func (s *GpmAPI) Init() error {
	var err error

	// Init API
	var aopts []apihttp.Option

	if ROOT == "" {
		if gruntime.GOOS == "windows" {
			ROOT = "C:\\opt\\gpm"
		} else {
			ROOT = "/opt/gpm"
		}
		_ = os.MkdirAll(filepath.Join(ROOT, "logs"), 0o777)
		_ = os.MkdirAll(filepath.Join(ROOT, "services"), 0o777)
		_ = os.MkdirAll(filepath.Join(ROOT, "packages"), 0o777)
	}

	gh, err := ssl.GetSSL(ROOT)
	if err != nil {
		return fmt.Errorf("load server tls: %v", err)
	}
	tls, err := ssl.GetTLS()
	if err != nil {
		return fmt.Errorf("load client tls: %v", err)
	}

	ghTLSOption := func() vine.Option { return func(o *vine.Options) { _ = o.Server.Init(grpcServer.GrpcToHttp(gh)) } }
	_ = ghTLSOption()
	cliTLSOption := func() vine.Option { return func(o *vine.Options) { _ = o.Client.Init(grpcClient.AuthTLS(tls)) } }
	_ = cliTLSOption()

	var clisrc source.Source

	opts := []vine.Option{
		vine.Registry(memory.NewRegistry()),
		vine.Name(runtime.GpmName),
		vine.Id(runtime.GpmId),
		vine.Version(runtime.GetVersion()),
		vine.Address(Address),
		vine.Metadata(map[string]string{
			"api-address": APIAddress,
			"namespace":   runtime.Namespace,
		}),
		//ghTLSOption(),
		//cliTLSOption(),
		vine.Flags(flags...),
		vine.WrapHandler(newLoggerWrapper()),
		vine.Action(func(c *cli.Context) error {
			clisrc = ccli.NewSource(ccli.Context(c))

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

			if c.Bool("enable-log") {
				EnableLog = true

				l, err := zap.New(zap.WithFileWriter(zap.FileWriter{
					FileName:   filepath.Join(ROOT, "logs", "gpmd.log"),
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

	or, _ := release.Get()
	log.Infof("system information: %s", or)

	if err = config.Load(clisrc); err != nil {
		log.Fatal(err)
	}

	db := new(repo.DB)
	if err = inject.Provide(s.Service, s.Client(), s, db); err != nil {
		return err
	}

	if err = inject.Populate(); err != nil {
		return err
	}

	if err = s.G.Init(); err != nil {
		return err
	}

	s.API = &RestAPI{}
	if err = s.API.Init(aopts...); err != nil {
		return err
	}

	if err = pb.RegisterGpmServiceHandler(s.Service.Server(), s); err != nil {
		return err
	}

	return err
}

func (s *GpmAPI) Run() error {
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

func New(opts ...vine.Option) *GpmAPI {
	srv := vine.NewService(opts...)
	return &GpmAPI{
		Service: srv,
	}
}

func newLoggerWrapper() vserver.HandlerWrapper {
	return func(fn vserver.HandlerFunc) vserver.HandlerFunc {
		return func(ctx context.Context, req vserver.Request, rsp interface{}) error {
			buf := bytes.NewBuffer([]byte(""))
			buf.WriteString("[" + req.ContentType() + "] ")
			if req.Stream() {
				buf.WriteString("[stream] ")
			}
			now := time.Now()
			err := fn(ctx, req, rsp)
			buf.WriteString(fmt.Sprintf("[%s] ", time.Now().Sub(now).String()))
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
