package routing

import (
	"context"
	"net/http"
)

type Context struct {
	context.Context
	Handler http.HandlerFunc
	Vars    RouteVars
}

type RouteVars struct {
	Keys, Values []string
}

func (rv *RouteVars) Set(key, value string) {
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

func newRouteVars() *RouteVars {
	return &RouteVars{
		[]string{},
		[]string{},
	}
}

func NewCtx(ctx context.Context) *Context {
	if ctx == nil {
		return &Context{
			Context: context.Background(),
		}
	}

	return &Context{
		Context: ctx,
	}
}
