package zephyr

import (
	"context"
	"net/http"
	"sync"
	"time"
)

func newMuxer() *muxer {
	m := &muxer{
		Root: &Node{
			Type:     Root,
			Value:    "/",
			Handlers: newHandlers(),
			Leaf:     true,
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
	Root *Node

	Server *http.Server

	pool *sync.Pool
}

func (m *muxer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ztx := newContext()

	// configure the ctx
	ztx.Routes = m.Root
	ztx.configure(r)

	// serve with the ctx
	r = r.WithContext(context.WithValue(r.Context(), "special-context", ztx))
	ztx.Handler.ServeHTTP(w, r)
}

func Vars(ctx context.Context, key string) string {
	realCtx, ok := ctx.Value("special-context").(*Context)
	if !ok {
		return ""
	}

	return realCtx.Vars.Get(key)
}
