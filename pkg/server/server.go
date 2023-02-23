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

	"github.com/spf13/cobra"
	pb "github.com/vine-io/gpm/api/service/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/inject"
	"github.com/vine-io/gpm/pkg/internal/store"
	"github.com/vine-io/gpm/pkg/service"
	"github.com/vine-io/pkg/release"
	"github.com/vine-io/vine/lib/api/handler/openapi"

	"github.com/vine-io/plugins/logger/zap"
	"github.com/vine-io/vine"
	vserver "github.com/vine-io/vine/core/server"
	grpcServer "github.com/vine-io/vine/core/server/grpc"
	log "github.com/vine-io/vine/lib/logger"
	uc "github.com/vine-io/vine/util/config"
	"google.golang.org/grpc/peer"
)

var (
	Address = ":7700"
	ROOT    = "/opt/gpm"
)

type GpmAPI struct {
	vine.Service

	G service.Manager `inject:""`
	T service.FTP     `inject:""`
}

func (s *GpmAPI) Init() error {
	var err error

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

	opts := []vine.Option{
		vine.Name(internal.GpmName),
		vine.ID(internal.GpmId),
		vine.Version(internal.GetVersion()),
		vine.Address(Address),
		vine.Metadata(map[string]string{
			"namespace": internal.Namespace,
		}),
		vine.WrapHandler(newLoggerWrapper()),
		vine.Action(func(c *cobra.Command, args []string) error {

			Address = uc.GetString("server-address")

			zap.WithFileWriter(zap.FileWriter{
				FileName:   filepath.Join(ROOT, "logs", "gpmd.log"),
				MaxSize:    1,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   false,
			})

			l, err := zap.New(zap.WithJSONEncode())
			if err != nil {
				return err
			}
			log.DefaultLogger = l

			return nil
		}),
	}

	s.Service.Init(opts...)
	app := newAPIServer(s.Service)
	if err = s.Server().Init(grpcServer.HttpHandler(app)); err != nil {
		log.Fatal(err)
	}

	or, _ := release.Get()
	log.Infof("system information: %s", or)

	db := new(store.DB)
	if err = inject.Provide(s.Service, s.Client(), s, db); err != nil {
		return err
	}

	if err = inject.Populate(); err != nil {
		return err
	}

	if err = s.G.Init(); err != nil {
		return err
	}

	if err = openapi.RegisterOpenAPIHandler(s.Service.Server()); err != nil {
		return err
	}
	if err = pb.RegisterGpmServiceHandler(s.Service.Server(), s); err != nil {
		return err
	}

	return err
}

func (s *GpmAPI) Run() error {
	if err := s.Service.Run(); err != nil {
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
