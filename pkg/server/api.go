package server

import (
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
	"github.com/vine-io/vine"
	vclient "github.com/vine-io/vine/core/client"
	"github.com/vine-io/vine/core/registry"
	ahandler "github.com/vine-io/vine/lib/api/handler"
	"github.com/vine-io/vine/lib/api/handler/openapi"
	arpc "github.com/vine-io/vine/lib/api/handler/rpc"
	"github.com/vine-io/vine/lib/api/resolver"
	"github.com/vine-io/vine/lib/api/resolver/grpc"
	"github.com/vine-io/vine/lib/api/router"
	regRouter "github.com/vine-io/vine/lib/api/router/registry"
	log "github.com/vine-io/vine/lib/logger"
	"github.com/vine-io/vine/util/namespace"
)

const (
	HandlerType = "rpc"
	Type        = "api"
	APIPath     = "/"
)

func newAPIServer(s vine.Service, client vclient.Client) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.Use(gin.Recovery())

	app.GET("/metrics", gin.WrapH(promhttp.Handler()))
	app.GET("/api/v1/endpoints", getEndpointsHandle)
	app.POST("/api/v1/push", pushHandle)

	prefix := "/debug/pprof"
	prefixRouter := app.Group(prefix)
	{
		prefixRouter.GET("/", gin.WrapF(pprof.Index))
		prefixRouter.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		prefixRouter.GET("/profile", gin.WrapF(pprof.Profile))
		prefixRouter.POST("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/symbol", gin.WrapF(pprof.Symbol))
		prefixRouter.GET("/trace", gin.WrapF(pprof.Trace))
		prefixRouter.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		prefixRouter.GET("/block", gin.WrapH(pprof.Handler("block")))
		prefixRouter.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		prefixRouter.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		prefixRouter.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		prefixRouter.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	openapi.RegisterOpenAPI(app)
	// create the namespace resolver
	nsResolver := namespace.NewResolver(Type, internal.Namespace)
	// resolver options
	rops := []resolver.Option{
		resolver.WithNamespace(nsResolver.ResolveWithType),
		resolver.WithHandler(HandlerType),
	}

	log.Infof("Registering API RPC Handler at %s", APIPath)
	rr := grpc.NewResolver(rops...)
	rt := regRouter.NewRouter(
		router.WithHandler(arpc.Handler),
		router.WithResolver(rr),
		router.WithRegistry(s.Options().Registry),
	)

	rp := arpc.NewHandler(
		ahandler.WithNamespace(internal.Namespace),
		ahandler.WithRouter(rt),
		ahandler.WithClient(client),
	)

	app.Use(rp.Handle)

	return app
}

func getEndpointsHandle(ctx *gin.Context) {
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

func pushHandle(ctx *gin.Context) {
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
