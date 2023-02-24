package wrap

import (
	"bytes"
	"context"
	"fmt"
	"time"

	vserver "github.com/vine-io/vine/core/server"
	log "github.com/vine-io/vine/lib/logger"
	"google.golang.org/grpc/peer"
)

func NewLoggerWrapper() vserver.HandlerWrapper {
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
