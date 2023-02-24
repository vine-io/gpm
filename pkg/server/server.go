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
	"fmt"
	"os"
	"path/filepath"
	gruntime "runtime"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/config"
	"github.com/vine-io/gpm/pkg/internal/store"
	"github.com/vine-io/gpm/pkg/internal/wrap"
	"github.com/vine-io/gpm/pkg/service"
	"github.com/vine-io/pkg/release"
	"github.com/vine-io/plugins/logger/zap"
	"github.com/vine-io/vine"
	"github.com/vine-io/vine/core/registry/mdns"
	grpcServer "github.com/vine-io/vine/core/server/grpc"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	log "github.com/vine-io/vine/lib/logger"
	uc "github.com/vine-io/vine/util/config"
)

var Flag = pflag.NewFlagSet("gpm", pflag.ContinueOnError)

func init() {

	uc.SetConfigName("gpm.yml")
	uc.SetConfigType("yaml")
	uc.AddConfigPath(".")
	uc.AddConfigPath("deploy")
	uc.AddConfigPath("config")

	mdns.DefaultMdnsDomain = "gpm"
	Flag.String("gpm.root", config.DefaultRoot, "Sets the Base directory for gpm service")
}

type GpmApp struct {
	s vine.Service
}

func New(opts ...vine.Option) (*GpmApp, error) {
	s := vine.NewService(opts...)

	app := &GpmApp{
		s: vine.NewService(opts...),
	}

	opts = append(opts,
		vine.Name(internal.GpmName),
		vine.ID(internal.GpmId),
		vine.Version(internal.GetVersion()),
		vine.Metadata(map[string]string{
			"namespace": internal.Namespace,
		}),
		vine.WrapHandler(wrap.NewLoggerWrapper()),
		vine.Action(func(c *cobra.Command, args []string) error {

			root := uc.GetString("gpm.root")
			if root == "" {
				if gruntime.GOOS == "windows" {
					root = "C:\\opt\\gpm"
				} else {
					root = "/opt/gpm"
				}
				_ = os.MkdirAll(filepath.Join(root, "logs"), 0o777)
				_ = os.MkdirAll(filepath.Join(root, "services"), 0o777)
				_ = os.MkdirAll(filepath.Join(root, "packages"), 0o777)
			}

			lopts := []log.Option{zap.WithJSONEncode()}
			filename := uc.GetString("logger.zap.filename")
			if filename != "" {
				writer := zap.FileWriter{
					FileName:   filename,
					MaxSize:    1,
					MaxBackups: 5,
					MaxAge:     30,
					Compress:   false,
				}

				if v := uc.GetInt("logger.zap.max-size"); v != 0 {
					writer.MaxSize = v
				}
				if v := uc.GetInt("logger.zap.max-backups"); v != 0 {
					writer.MaxBackups = v
				}
				if v := uc.GetInt("logger.zap.max-age"); v != 0 {
					writer.MaxAge = v
				}
				if v := cast.ToBool(uc.Get("logger.zap.compress")); v {
					writer.Compress = true
				}

				lopts = append(lopts, zap.WithFileWriter(writer))
			}

			l, err := zap.New(lopts...)
			if err != nil {
				return err
			}
			log.DefaultLogger = l

			return nil
		}),
	)

	// vine service 初始化，解析命令行参数
	if err := app.s.Init(opts...); err != nil {
		return nil, err
	}

	if err := uc.UnmarshalKey(&config.DefaultConfig, "gpm"); err != nil {
		return nil, fmt.Errorf("unmarshal config file: %v", err)
	}
	_ = uc.UnmarshalKey(&config.DefaultAddress, "server.address")

	ctx := app.s.Options().Context
	reg := app.s.Options().Registry
	client := app.s.Options().Client
	server := app.s.Options().Server

	db := new(store.DB)
	manager, err := service.NewManagerService(ctx, server, db)
	if err != nil {
		return nil, err
	}

	sftp, err := service.NewSFtpService(ctx, server)
	if err != nil {
		return nil, err
	}

	if err = RegistryGpmRpcServer(ctx, server, manager, sftp); err != nil {
		return nil, err
	}

	if err = openapi.RegisterOpenAPIHandler(server); err != nil {
		return nil, err
	}

	handler, err := RegistryGpmAPIServer(ctx, reg, client)
	if err != nil {
		return nil, err
	}

	if err = s.Server().Init(grpcServer.HttpHandler(handler)); err != nil {
		return nil, err
	}

	or, _ := release.Get()
	log.Infof("system information: %s", or)

	return app, nil
}

func (app *GpmApp) Run() error {
	if err := app.s.Run(); err != nil {
		return err
	}

	return nil
}
