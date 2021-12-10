package log

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"api.cloud.io/pkg/security/auth"
	"github.com/getsentry/raven-go"
	"github.com/semrush/zenrpc"
)

// ErrorLogger save internal errors of request.
func ErrorLogger() zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			r := h(ctx, method, params)
			if r.Error != nil && r.Error.Code == 50001 {
				ns := zenrpc.NamespaceFromContext(ctx)
				token := auth.ContextToken(ctx)
				req, _ := zenrpc.RequestFromContext(ctx)
				saveRequestErr(ns, method, token, req, r.Error.Err)
			}

			return r
		}
	}
}

func saveRequestErr(namespace, method string, token *auth.Token, req *http.Request, err error) {
	tags := map[string]string{
		"method": namespace + "." + method,
	}
	extra := map[string]interface{}{
		"token": token,
	}
	inters := []raven.Interface{
		raven.NewHttp(req),
	}

	packet := raven.NewPacketWithExtra(err.Error(), extra, inters...)
	eventID, ch := raven.Capture(packet, tags)
	if eventID != "" {
		sentryErr := <-ch
		if sentryErr != nil {
			fmt.Printf("could not capture (%s) request error: %s\n", err, sentryErr)
		}
	}
}
