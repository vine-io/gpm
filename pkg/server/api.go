package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pbr "github.com/schollz/progressbar/v3"
	gpmv1 "github.com/vine-io/gpm/api/types/gpm/v1"
	"github.com/vine-io/gpm/pkg/internal"
	"github.com/vine-io/gpm/pkg/internal/client"
	vclient "github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/registry"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	log "github.com/vine-io/vine/lib/logger"
	uapi "github.com/vine-io/vine/util/api"
)

type GpmHttpServer struct {
	*gin.Engine

	register registry.Registry
	client   vclient.Client
}

func RegistryGpmAPIServer(ctx context.Context, reg registry.Registry, client vclient.Client) (http.Handler, error) {

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())

	s := &GpmHttpServer{
		register: reg,
		client:   client,
	}

	openapi.RegisterOpenAPI(client, reg, app)

	s.GET("/metrics", gin.WrapH(promhttp.Handler()))
	s.GET("/api/v1/endpoints", s.getEndpointsHandle)
	s.POST("/api/v1/push", s.pushHandle)

	prefix := "/debug/pprof"
	group := s.Group(prefix)
	{
		group.GET("/", gin.WrapF(pprof.Index))
		group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		group.GET("/profile", gin.WrapF(pprof.Profile))
		group.POST("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/trace", gin.WrapF(pprof.Trace))
		group.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		group.GET("/block", gin.WrapH(pprof.Handler("block")))
		group.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		group.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		group.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		group.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	ns := internal.Namespace
	uapi.PrimpHandler(app, reg, client, ns)
	s.Engine = app

	return s, nil
}

func (s *GpmHttpServer) getEndpointsHandle(ctx *gin.Context) {
	endpoints := make([]map[string]string, 0)
	keys := make(map[string]struct{}, 0)
	list, _ := registry.GetService(ctx, internal.GpmName)
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

	ctx.JSON(200, gin.H{
		"data": endpoints,
	})
}

func (s *GpmHttpServer) pushHandle(ctx *gin.Context) {
	mf, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	fs, ok := mf.File["file"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "missing file")
		return
	}
	vv := mf.Value["dst"]
	if !ok {
		ctx.JSON(http.StatusBadRequest, "missing dst")
		return
	}
	file := fs[0]
	dst := vv[0]

	fd, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "open file: "+err.Error())
		return
	}
	defer fd.Close()

	in := &gpmv1.PushIn{
		Name:  file.Filename,
		Dst:   dst,
		Total: file.Size,
	}

	cc := client.New()
	opts := []vclient.CallOption{
		vclient.WithDialTimeout(time.Hour * 2),
		vclient.WithStreamTimeout(time.Hour * 2),
	}

	stream, err := cc.Push(ctx, opts...)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, "connect to gpm server: "+err.Error())
		return
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
			ctx.JSON(http.StatusInternalServerError, e.Error())
			return
		}

		if n > 0 {
			_ = pb.Add(n)
			in.Length = int64(n)
			in.Chunk = buf[0:n]
			err = stream.Send(in)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		if e == io.EOF {
			break
		}
	}

	if err = stream.Wait(); err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "OK",
	})
}
