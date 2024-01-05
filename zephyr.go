package zephyr

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"syreclabs.com/go/faker"
)

type View interface {
	Render(ctx context.Context, w io.Writer) error
	http.Handler
}

func defaultHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1>Hello!</h1>"))
		panic("recover me or else!")
	}
}

type Zephyr struct {
	Pages  map[string]View
	pool   *sync.Pool
	router *chi.Mux
	server *http.Server
	fs     fs.FS
}

type zctx struct {
	specialSauce string
}

type zctxKey struct {
	key string
}

var specialCtxKey *zctxKey = &zctxKey{"specialkey"}

func New() *Zephyr {
	r := chi.NewRouter()

	pool := &sync.Pool{
		New: func() any {
			return &zctx{specialSauce: faker.Lorem().Word()}
		}}

	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		sauce := r.Context().Value(specialCtxKey).(*zctx)
		w.Write([]byte(fmt.Sprintf(`{"hello": "%v"}`, sauce.specialSauce)))
		pool.Put(sauce)
	})

	z := Zephyr{
		router: r,
		pool:   pool,
		server: &http.Server{
			Handler:      r,
			Addr:         ":3000",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			ConnContext: func(ctx context.Context, c net.Conn) context.Context {
				return context.WithValue(ctx, specialCtxKey, pool.Get())
			},
		},
	}

	return &z
}

func (z *Zephyr) Run(addr string) error {
	if !strings.Contains(addr, ":") {
		addr = ":" + addr
	}

	z.registerViews()

	return z.server.ListenAndServe()
}

func (z *Zephyr) RegisterFileServe(urlPattern, dir string) error {
	z.fs = os.DirFS(dir)

	z.router.Get(urlPattern, func(w http.ResponseWriter, r *http.Request) {
		r.Header["Accept-Origin"] = []string{"*"}
		r.Header["Content-Type"] = []string{"text/javsacript"}

		uriSegments := strings.Split(strings.TrimPrefix(r.RequestURI, "/"), "/")
		filePath := strings.Join(uriSegments[1:], "/")

		file, err := z.fs.Open(filePath)
		if err != nil {
			logrus.Error("open err", err)
			w.WriteHeader(404)
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			logrus.Error("io err", err)
			w.WriteHeader(404)
			return
		}

		w.Write(data)

		logrus.Debugf("served %v", filePath)
	})

	return nil
}

func (z *Zephyr) AddViews(views map[string]View) error {
	for path, view := range views {
		z.Pages[path] = view
	}
	return nil
}

func (z *Zephyr) registerViews() error {
	for path, view := range z.Pages {
		z.router.Get(path, view.ServeHTTP)
	}
	return nil
}

func authenticateFsRequest(r *http.Request) error {

	return nil
}
