package routing

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type muxer struct {
	Root *Node
}

func Vars(ctx context.Context, key string) string {
	realCtx, ok := ctx.Value("special-context").(*Context)
	if !ok {
		return ""
	}

	return realCtx.Vars.Get(key)
}

func (m *muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value("special-context").(*Context)
	m.SetHandler(ctx, r.RequestURI[1:])
	ctx.Handler.ServeHTTP(w, r)
}

func (m *muxer) SetHandler(ctx *Context, route string) {
	segments := strings.Split(strings.Trim(route, "/"), "/")
	fmt.Println(segments)
	if len(segments) == 0 {
		ctx.Handler = m.Root.Handler
		return
	}

	m.Root.configureCtx(ctx, segments)
}

func runMuxer() error {
	m := &muxer{
		Root: &Node{Type: Root, Value: "/"},
	}

	server := &http.Server{
		Handler:      m,
		Addr:         ":8080",
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return context.WithValue(ctx, "special-context", NewCtx(ctx))
		},
	}

	tokes1 := GetRouteTokens("/hello/world/{thatsIncredible}/*")
	tokes2 := GetRouteTokens("/hello/world/~regex/*/wow")
	tokes3 := GetRouteTokens("/hello/world/{whattheheckisup}")

	m.Root.insert(tokes1)
	m.Root.insert(tokes2)
	m.Root.insert(tokes3)

	return server.ListenAndServe()
}
