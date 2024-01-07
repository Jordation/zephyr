package zephyr

import (
	"fmt"
	"testing"
)

func TestInsertRoutes(t *testing.T) {
	var (
		routes = []string{
			"/hello/world",
			"/hello",
			"/hello/newpath/world",
			"/hello/world/{ID}/*/~[0-9]",
			"/hello/world/{ID}/*/helper",
			"/hello/world/{ID}/newmate",
			"/newpath",
			"/~anothernew",
			"/*",
			"/{parameterpath}",
			"/",
		}
		root = newNode(RouteToken{Type: Root, Value: "/"})
	)

	for _, route := range routes {
		root.addRoute(GetTokensFromRoute(nil, route), 0, nil)
	}

	fmt.Println("router")
}
