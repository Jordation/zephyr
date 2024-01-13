package zephyr

import (
	"strings"
)

type Zephyr struct {
	*mux
}

func New() *Zephyr {
	z := Zephyr{
		mux: NewMux(),
	}
	return &z
}

func cleanRouteSegs(route string) []string {
	return strings.Split(strings.Trim(route, "/"), "/")
}
