package zephyr

import (
	"io/fs"
	"strings"
)

type Zephyr struct {
	muxer *muxer
	fs    fs.FS
}

func New() *Zephyr {
	z := Zephyr{
		muxer: newMuxer(),
	}
	return &z
}

// Run blocks.
func (z *Zephyr) Run(addr string) error {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	return z.muxer.Server.ListenAndServe()
}

func cleanRouteSegs(route string) []string {
	return strings.Split(strings.Trim(route, "/"), "/")
}
