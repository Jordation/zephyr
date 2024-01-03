package routing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTokens(t *testing.T) {
	routes := []string{
		"something/simple/like/this",
		"something/{with_param}/~ortheRegey[0-5]",
		"maybe/*/wcthat/route",
	}

	tests := []struct {
		expected []RouteType
	}{
		{[]RouteType{Path, Path, Path, Path}},
		{[]RouteType{Path, Param, Regex}},
		{[]RouteType{Path, WildCard, Path, Path}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test:%v:", i+1), func(t *testing.T) {
			tokes := GetRouteTokens(routes[i])
			for j, toke := range tokes {
				assert.Equal(t, test.expected[j], toke.Type, fmt.Sprintf("failed test %v.%v!", i+1, j+1))
			}
		})
	}

}
