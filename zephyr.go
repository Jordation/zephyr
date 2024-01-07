package zephyr

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type View interface {
	Render(ctx context.Context, w io.Writer) error
	http.Handler
}

type Zephyr struct {
	Pages map[string]View
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

	z.registerViews()

	return z.muxer.Server.ListenAndServe()
}

func (z *Zephyr) RegisterFileServe(urlPattern, dir string) error {
	z.fs = os.DirFS(dir)

	z.GET(urlPattern, func(w http.ResponseWriter, r *http.Request) {
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
		z.GET(path, view.ServeHTTP)
	}
	return nil
}

func cleanRouteSegs(route string) []string {
	return strings.Split(strings.Trim(route, "/"), "/")
}
