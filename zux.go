package zephyr

import (
	"context"
	"net/http"
	"sync"
)

var ctxPool = &sync.Pool{
	New: func() any { return newContext() },
}

func NewMux() *mux {
	m := &mux{
		root: &node{
			routeType: Root,
			value:     "/",
			handlers:  newHandlers(),
		},

		pool: ctxPool,
	}
	return m
}

type mux struct {
	root *node

	pool *sync.Pool
}

// Run is blocking
func (m *mux) Run(addr string) error {
	return http.ListenAndServe(addr, m)
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := m.pool.Get().(*zCtx)

	ctx.reset()
	defer func() {
		m.pool.Put(ctx)
	}()

	// configure
	ctx.routes = m.root
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
