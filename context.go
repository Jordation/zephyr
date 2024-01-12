package zephyr

import (
	"context"
	"fmt"
	"net/http"
)

func GetCtx(ctx context.Context) *zCtx {
	realCtx, ok := ctx.Value("special-context").(*zCtx)
	if !ok {
		panic("failed to get real context bro")
	}

	return realCtx
}

func newContext() *zCtx {
	return &zCtx{}
}

type zCtx struct {
	handler http.HandlerFunc

	mw []http.Handler

	vars RouteVars

	method uint8 // we use a map to convert method to an index

	routes *node

	recover bool
}

func (ctx *zCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ctx.recover {
		defer func() {
			err := recover()
			fmt.Println(err)
		}()
	}

	for _, handler := range ctx.mw {
		handler.ServeHTTP(w, r)
	}

	ctx.handler.ServeHTTP(w, r)
}

func (ctx *zCtx) configure(r *http.Request) {
	ctx.method = methodToIndexMap[r.Method]
	isRoot := r.URL.Path == "/"
	path := cleanRouteSegs(r.URL.Path)

	if isRoot {
		path = nil // will set the handler to the first node traverse is called on
	}

	last := ctx.routes.traverse(ctx, path)
	if last != nil { // wasn't able to traverse the full route
		ctx.handler = http.NotFound
	}

	if ctx.handler == nil { // no handler registered for method of r
		ctx.handler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (c *zCtx) reset() {
	c.handler = nil

	c.mw = c.mw[:0]

	c.vars.Reset()

	c.method = 99

	c.routes = nil

	c.recover = false

}

type RouteVars struct {
	Keys, Values []string
}

func (rv *RouteVars) Reset() {
	if rv != nil {
		rv.Keys = rv.Keys[:0]
		rv.Values = rv.Values[:0]
	}
}

func (rv *RouteVars) Set(key, value string) {
	if rv == nil {
		rv = &RouteVars{[]string{}, []string{}}
	}

	rv.Keys = append(rv.Keys, key)
	rv.Values = append(rv.Values, value)
}

func (rv *RouteVars) Get(key string) string {
	for i, v := range rv.Keys {
		if v == key {
			return rv.Values[i]
		}
	}

	return ""
}
