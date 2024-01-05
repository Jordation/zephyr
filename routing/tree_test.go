package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		}
		root = newNode(RouteToken{Type: Root, Value: "/"}, nil)
	)

	for _, route := range routes {
		root.insert(GetRouteTokens(route))
	}

	gotRoutes := root.walk(true)

	for _, insertedRoute := range gotRoutes {
		assert.Contains(t, routes, insertedRoute)
	}
}

func TestHT(t *testing.T) {
	sl := []int{1, 2, 3, 4, 5}
	h, tail := ht(sl)
	assert.Equal(t, h, 1)
	assert.Equal(t, tail, []int{2, 3, 4, 5})
}
