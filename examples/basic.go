package examples

import (
	"net/http"

	"github.com/Jordation/zephyr"
)

func example() {
	z := zephyr.New()
	z.GET("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	sub := zephyr.New()
	z.GET("/sub/route", sub.ServeHTTP)

	// registers to /sub/route/hello
	sub.GET("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("im from the subby!"))
	})

}
