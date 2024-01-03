package routing

type RouteType int8

const (
	Separator RouteType = iota
	Path
	Param
	Regex
	WildCard
	EndRoute
	Root
)

type RouteToken struct {
	Type  RouteType
	Value string
}

func GetRouteTokens(route string) []RouteToken {
	tokens, _, _, _, _ := getRouteTokensRecursively(nil, route, 0, 0, 'x')
	return tokens
}

func getRouteTokensRecursively(tokes []RouteToken, route string, pos, rPos int, ch byte) ([]RouteToken, string, int, int, byte) {
	pos, rPos, ch = readChar(route, pos, rPos, ch)

	if ch == 0 {
		return tokes, "", 0, 0, '0'
	}

	var identifier string

	switch ch {
	case '/':
		break
	case '{':
		pos, rPos, ch = readChar(route, pos, rPos, ch)
		identifier, pos, rPos, ch = readTill(route, pos, rPos, ch, '}')
		tokes = append(tokes, newToken(Param, identifier))
	case '*':
		tokes = append(tokes, newToken(WildCard, ""))
	case '~':
		pos, rPos, ch = readChar(route, pos, rPos, ch)
		identifier, pos, rPos, ch = readTill(route, pos, rPos, ch, '/')
		tokes = append(tokes, newToken(Regex, identifier))
	default:
		identifier, pos, rPos, ch = readTill(route, pos, rPos, ch, '/')
		tokes = append(tokes, newToken(Path, identifier))
	}

	return getRouteTokensRecursively(tokes, route, pos, rPos, ch)
}

func readChar(route string, pos, rPos int, ch byte) (int, int, byte) {
	if rPos >= len(route) {
		ch = 0
	} else {
		ch = route[rPos]
	}
	pos = rPos
	rPos += 1
	return pos, rPos, ch
}

func readTill(route string, pos, rPos int, ch byte, mark byte) (string, int, int, byte) {
	initialPos := pos
	for ch != mark && ch != 0 {
		pos, rPos, ch = readChar(route, pos, rPos, ch)
	}
	ident := route[initialPos:pos]
	return ident, pos, rPos, ch
}

func newToken(tokenType RouteType, value string) RouteToken {
	return RouteToken{Type: tokenType, Value: value}
}
