package zephyr

import "net/http"

// cascade will override the target nodes value
func (z *Zephyr) Use(path string, cascade bool, handlers ...http.Handler) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, 0, nil, handlers, cascade)
}
func (z *Zephyr) GET(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodGet], f, nil, false)
}
func (z *Zephyr) HEAD(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodHead], f, nil, false)
}
func (z *Zephyr) POST(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodPost], f, nil, false)
}
func (z *Zephyr) PUT(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodPut], f, nil, false)
}
func (z *Zephyr) PATCH(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodPatch], f, nil, false)
}
func (z *Zephyr) DELETE(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodDelete], f, nil, false)
}
func (z *Zephyr) CONNECT(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodConnect], f, nil, false)
}
func (z *Zephyr) OPTIONS(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodOptions], f, nil, false)
}
func (z *Zephyr) TRACE(path string, f http.HandlerFunc) {
	tokens := GetTokensFromRoute(nil, path)
	z.mux.root.insert(tokens, methodToIndexMap[http.MethodTrace], f, nil, false)
}
