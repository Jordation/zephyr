package zephyr

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func newMuxer() *muxer {
	m := &muxer{
		Root: &node{
			routeType: Root,
			value:     "/",
			handlers:  newHandlers(),
			leaf:      true,
		},

		pool: &sync.Pool{New: func() any { return newContext() }},
	}

	m.Server = &http.Server{
		Handler:      m,
		Addr:         ":3000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return m
}

type muxer struct {
	Root *node

	Server *http.Server

	pool *sync.Pool
}

func (m *muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := newContext()

	// configure the ctx
	ctx.Routes = m.Root
	ctx.configure(r)

	// serve with the ctx
	r = r.WithContext(context.WithValue(r.Context(), "special-context", ctx))
	ctx.Handler.ServeHTTP(w, r)
}

func Vars(ctx context.Context, key string) string {
	realCtx, ok := ctx.Value("special-context").(*Context)
	if !ok {
		return ""
	}

	return realCtx.Vars.Get(key)
}
