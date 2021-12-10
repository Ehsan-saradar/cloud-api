package auth

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"api.cloud.io/pkg/errors"
	"github.com/semrush/zenrpc"
)

type contextKey string

const (
	tokenContextKey     contextKey = "token"
	authScopeContextKey contextKey = "scope"
)

// Authorizer check the scopes of methods in JWT.
func Authorizer(methodsScopes map[string]map[string][]string) zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			namespace := zenrpc.NamespaceFromContext(ctx)
			scopes := methodsScopes[namespace][method]
			if len(scopes) == 0 {
				return h(ctx, method, params)
			}

			token := ContextToken(ctx)
			if token == nil {
				res := zenrpc.Response{}
				res.Set(nil, errors.ErrUnauthorized)

				return res
			}

			// Sort scopes if needed
			if !sort.StringsAreSorted(token.Scopes) {
				sort.Strings(token.Scopes)
			}

			for _, s := range scopes {
				if sort.SearchStrings(token.Scopes, s) < len(token.Scopes) {
					ctx = context.WithValue(ctx, authScopeContextKey, s)

					return h(ctx, method, params)
				}
			}

			res := zenrpc.Response{}
			res.Set(nil, errors.ErrAccessDenied)

			return res
		}
	}
}

// ContextScope get the selected Scope of request.
func ContextScope(ctx context.Context) string {
	if ctx.Value(authScopeContextKey) != nil {
		return ctx.Value(authScopeContextKey).(string)
	}

	return ""
}

// JWTMiddleware fetch JWT from Authorization header of HTTP request.
func JWTMiddleware() zenrpc.MiddlewareFunc {
	return func(h zenrpc.InvokeFunc) zenrpc.InvokeFunc {
		return func(ctx context.Context, method string, params json.RawMessage) zenrpc.Response {
			req, ok := zenrpc.RequestFromContext(ctx)
			if !ok {
				return h(ctx, method, params)
			}

			authHeader := req.Header.Get("Authorization")
			segments := strings.Split(authHeader, " ")
			if len(segments) != 2 {
				return h(ctx, method, params)
			}
			if segments[0] != "Bearer" {
				return h(ctx, method, params)
			}

			token, err := ParseSignedToken(segments[1])
			if err != nil {
				res := zenrpc.Response{}
				res.Set(nil, err)

				return res
			}

			err = token.Validate()
			if err != nil {
				res := zenrpc.Response{}
				res.Set(nil, err)

				return res
			}

			ctx = context.WithValue(ctx, tokenContextKey, token)

			return h(ctx, method, params)
		}
	}
}
