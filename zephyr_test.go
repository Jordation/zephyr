package zephyr

import (
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListen(t *testing.T) {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		panic("listen err: " + err.Error())
	}
	defer ln.Close()

	conn, err := ln.Accept()
	if err != nil {
		panic("accept err: " + err.Error())
	}

	buff := []byte("hello world!")
	_, err = conn.Write(buff)
	if err != nil {
		panic("write errs: " + err.Error())
	}
}

type testScenario struct {
	name          string
	inputRoute    string
	usageRoute    string
	handler       http.HandlerFunc
	mw            []http.Handler
	expectVars    *RouteVars
	expectBody    []byte
	expectHeaders http.Header
}

func defaultHandler(body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}
}

func Test_Zephman(t *testing.T) {
	tests := []testScenario{
		{
			inputRoute: "/hello",
			usageRoute: "/hello",
			handler:    defaultHandler("/hello"),
			expectBody: []byte("/hello"),
			mw:         []http.Handler{},
		},
		{
			inputRoute: "/hello/world",
			usageRoute: "/hello/world",
			handler:    defaultHandler("/hello/world"),
			expectBody: []byte("/hello/world"),
		},
		{
			inputRoute: "/hello/~[0-9]{4}/world",
			usageRoute: "/hello/1234/world",
			handler:    defaultHandler("/hello/~[0-9]{4}/world"),
			expectBody: []byte("/hello/~[0-9]{4}/world"),
		},
		{
			inputRoute: "/hello/{paramKey}",
			usageRoute: "/hello/paramValue",
			handler:    defaultHandler("/hello/{paramKey}"),
			expectBody: []byte("/hello/{paramKey}"),
			expectVars: &RouteVars{
				Keys:   []string{"paramKey"},
				Values: []string{"paramValue"},
			},
		},
		{
			inputRoute: "/hello/*",
			usageRoute: "/hello/hadfasdhfalshfd",
			expectBody: []byte("/hello/{paramKey}"),
		},
	}

	z := New()
	for _, s := range tests {
		z.GET(s.inputRoute, s.handler)
	}

	go func() {
		z.Run(":3000")
	}()

	for _, s := range tests {
		t.Run(s.name, func(t *testing.T) {
			res, err := http.Get("http://localhost:3000" + s.usageRoute)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, s.expectBody, body, "mismatch body for test on route %v, expected: %v, got: %v", s.usageRoute, string(s.expectBody), string(body))
		})

	}

	// test middlewares are applied
	// test routes hit in expected fashion
	// check for all route params
}
