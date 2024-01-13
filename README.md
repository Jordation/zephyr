## Potential Modules

- [ ] Minimal db interface (sqlx wrapper or similar?)
- [x] Cache
- [ ] Build / bundle JS files 
- [ ] Hot reload backend (go files)
- [x] Static FS
- [x] Register Endpoints
- [x] Register Middleware




# What it do

## Basic

```go
func main(){
    z := zephyr.New()
    z.GET("/", func(w http.ResponseWriter, r *http.Request) {
	    w.Write([]byte("Hello, world!"))
    })

    // blocks
    log.Fatal(z.Run())
}
```
Routes are registered on the method level
`405: Method Not Allowed` is written to the client if a registered route is matched with an unregistered method

## Routing Options

```go
    // POST: host/captured
    z.POST("/{route_param_between_braces}", func(w http.ResponseWriter, r *http.Request){
        value := zephyr.Vars(r.Context(), "route_param_between_braces")
        // ^ "caputred"

        value2 := zephyr.Vars(r.Context(), "unknown")
        // ^  ""
    })

    z.TRACE("/hello/~[A-Z][0-9]{2}.", func(w http.ResponseWriter, r *http.Request){})
    //              ^ indicate regex with squigly

    z.PATCH("/hello/*", func(w http.ResponseWriter, r *http.Request){ })
    //              ^ catch-all 
```

Routes are matched in order of priority, the order can be found in `route_tokenizer.go`. However it is as follows:

1. Exact match on route path
2. Regex match
3. Capture Param 
4. Wildcard

This means that if you register `/hello/{param}` and `/hello/*`, the latter will never be used to handle an incoming request.
If a match fails however, it will look to the next in priority such that routes `hello/{param}` `hello/~regex` could both handle requests, if the regex fails the param registered route will be the backup

 # TODO 
Sub-Routers are handled differently when added to a parent router and will nuke any other children attached to that route, i.e. if a `Root` type route (full path == "/") exists as a child, it will be the only child of that node

```go
func main(){
    z := zephyr.New()
    z.GET("/hello/world", func(w http.ResponseWriter, r *http.Request) {
	    w.Write([]byte("Hello, world!"))
    })

    // Get a new mux handler with a new root node
    sub := zephyr.NewMux()
    z.Use("/hello", sub.ServeHTTP)

    // Routes on the sub-router will match their parent router i.e. hello/{name}
    sub.GET("/{name}", func(w http.ResponseWriter, r *http.Request) {
        name := zephyr.Vars(r.Context(), "name")
	    w.Write([]byte(fmt.Sprintf("Hello, %v!", name)))
    })

    subby2 := zephyr.NewMux()
    sub.Get("/todos" subby2.ServeHTTP)

    // hello/{name}/todos/{id}
    subby2.GET("/{id}", {handler}) 
    // etc..
}
```
# !TODO

# Modules

Some helpful basic modules are provided 
... cache
... fileserve
... build api wrapper
... hot reload 