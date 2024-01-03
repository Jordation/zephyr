package zephyr

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type defaultFS struct {
	root        string
	once        *sync.Once
	middlewares chi.Middlewares
}

func newDefaultFs() *defaultFS {
	return &defaultFS{}
}

func (dfs *defaultFS) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	dfs.once.Do(func() {
		if dfs.root == "" {
			dfs.root = "/static"
		}
	})

}
