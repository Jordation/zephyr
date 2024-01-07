package zephyr

import "net/http"

func (z *Zephyr) GET(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodGet], f)
}
func (z *Zephyr) HEAD(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodHead], f)
}
func (z *Zephyr) POST(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodPost], f)
}
func (z *Zephyr) PUT(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodPut], f)
}
func (z *Zephyr) PATCH(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodPatch], f)
}
func (z *Zephyr) DELETE(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodDelete], f)
}
func (z *Zephyr) CONNECT(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodConnect], f)
}
func (z *Zephyr) OPTIONS(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodOptions], f)
}
func (z *Zephyr) TRACE(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.muxer.Root.addRoute(tokens, MethodToIndexMap[http.MethodTrace], f)
}
