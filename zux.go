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
	ctx := m.pool.Get().(*zCtx)

	ctx.reset()
	defer func() {
		m.pool.Put(ctx)
	}()

	// configure
	ctx.routes = m.Root
	ctx.configure(r)
	r = r.WithContext(context.WithValue(r.Context(), "special-context", ctx))

	// serve.
	ctx.ServeHTTP(w, r)
}

func Vars(ctx context.Context, key string) string {
	realCtx, ok := ctx.Value("special-context").(*zCtx)
	if !ok {
		return ""
	}

	return realCtx.vars.Get(key)
}
