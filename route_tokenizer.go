package zephyr

import (
	"strings"
)

type RouteType int8

const (
	Separator RouteType = iota
	Path
	Regex
	Param
	WildCard
	EndRoute
	Root
)

type RouteToken struct {
	Type  RouteType
	Value string
}

func GetTokensFromRoute(tokes []RouteToken, route string) []RouteToken {
	if len(route) == 0 {
		return tokes
	}

	var end int

outer:
	for i, ch := range route {
		switch ch {
		case '/', '}':
			if tokes == nil { // so we only create the root on first run
				tokes = append(tokes, newToken(Root, "/"))
			}

			end = 1
			break outer

		case '{':
			end = strings.IndexRune(route, '}')
			if end == -1 {
				panic("Couldn't find } to capture param")
			}
			tokes = append(tokes, newToken(Param, route[i+1:end]))
			break outer

		case '*':
			end = strings.IndexRune(route, '/')
			if end == -1 {
				end = len(route)
			}
			tokes = append(tokes, newToken(WildCard, "*"))
			break outer

		case '~':
			end = strings.IndexRune(route, '/')
			if end == -1 {
				end = len(route)
			}
			tokes = append(tokes, newToken(Regex, route[i+1:end]))
			break outer

		default:
			end = strings.IndexRune(route, '/')
			if end == -1 {
				end = len(route)
			}
			tokes = append(tokes, newToken(Path, route[i:end]))
			break outer
		}
	}

	return GetTokensFromRoute(tokes, route[end:])
}

func newToken(tokenType RouteType, value string) RouteToken {
	return RouteToken{Type: tokenType, Value: value}
}
