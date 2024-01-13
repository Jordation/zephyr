package zephyr

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testScenario struct {
	name          string
	inputRoute    string
	usageRoute    string
	handler       http.HandlerFunc
	mw            []http.Handler
	expectVars    *RouteVars
	expectBody    string
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
			expectBody: "/hello",
			mw: []http.Handler{
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("X-test", "yipee")
				}),
			},
			expectHeaders: http.Header{
				http.CanonicalHeaderKey("X-test"): []string{"yipee"},
			},
		},
		{
			inputRoute: "/hello/world",
			usageRoute: "/hello/world",
			handler:    defaultHandler("/hello/world"),
			expectBody: "/hello/world",
		},
		{
			inputRoute: "/hello/~[0-9]{4}/world",
			usageRoute: "/hello/1234/world",
			handler:    defaultHandler("/hello/~[0-9]{4}/world"),
			expectBody: "/hello/~[0-9]{4}/world",
		},
		{
			inputRoute: "/hello/{paramKey}",
			usageRoute: "/hello/paramValue",
			handler:    defaultHandler("/hello/{paramKey}"),
			expectBody: "/hello/{paramKey}",
			expectVars: &RouteVars{
				Keys:   []string{"paramKey"},
				Values: []string{"paramValue"},
			},
		},
		{
			inputRoute: "/hello/*",
			usageRoute: "/hello/hadfasdhfalshfd",
			expectBody: "/hello/{paramKey}",
		},
		{
			inputRoute: "/@/*",
			usageRoute: "/@/yolo.js",
			handler:    DefaultFsHandler("/").ServeHTTP,
			expectBody: "error opening file:open yolo.js: The system cannot find the file specified.",
		},
	}

	z := New()
	for _, s := range tests {
		z.GET(s.inputRoute, s.handler)
		z.Use(s.inputRoute, false, s.mw...)
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

			if s.expectHeaders != nil {
				assert.Equal(t, s.expectHeaders.Get("X-test"), res.Header.Get("X-test"))
			}

			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, s.expectBody, string(body), "mismatch body for test on route %v, expected: %v, got: %v", s.usageRoute, string(s.expectBody), string(body))

			if s.expectVars != nil {
				//assert.Equal(t, Vars())
			}
		})

	}

	// test middlewares are applied
	// test routes hit in expected fashion
	// check for all route params
}
