package zephyr

import (
	"context"
	"net/http"
)

// nil bad.
func GetCtx(ctx context.Context) *Context {
	realCtx, ok := ctx.Value("special-context").(*Context)
	if !ok {
		panic("failed to get real context bro")
	}

	return realCtx
}

func newContext() *Context {
	return &Context{}
}

type Context struct {
	Handler http.HandlerFunc

	Mw []http.Handler

	Vars *RouteVars

	Method uint8 // we use a map to convert method to an index

	Routes *node
}

func (ctx *Context) configure(r *http.Request) {
	ctx.Method = methodToIndexMap[r.Method]

	isRoot := r.URL.Path == "/"
	route := cleanRouteSegs(r.URL.Path)

	last := ctx.Routes.traverse(ctx, route, isRoot)
	if last != nil { // wasn't able to traverse the full route
		ctx.Handler = http.NotFound
	}

	if ctx.Handler == nil { // no handler registered for method of r
		ctx.Handler = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (c *Context) reset() {
	c.Handler = nil

	c.Mw = c.Mw[:0]

	c.Vars.Keys = c.Vars.Keys[:0]

	c.Vars.Values = c.Vars.Values[:0]

	c.Method = 99
}

type RouteVars struct {
	Keys, Values []string
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
