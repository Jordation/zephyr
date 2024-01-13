package zephyr

import (
	"net/http"
	"regexp"
	"slices"
	"sort"

	"github.com/sirupsen/logrus"
)

var methodToIndexMap = map[string]uint8{
	http.MethodGet:     0,
	http.MethodPost:    1,
	http.MethodPatch:   2,
	http.MethodPut:     3,
	http.MethodDelete:  4,
	http.MethodTrace:   5,
	http.MethodOptions: 6,
	http.MethodConnect: 7,
	http.MethodHead:    8,
}

type node struct {
	routeType RouteType

	rgx *regexp.Regexp

	value string

	leaf bool

	handlers  []http.HandlerFunc
	isHandler bool

	mw      []http.Handler
	cascade bool

	children
}

type children []*node

// traverse returns nil upon a successful walk to handler and otherwise, the last node it got to
func (n *node) traverse(ctx *zCtx, routeSegs []string) *node {
	if len(routeSegs) == 0 {
		ctx.handler = n.handlers[ctx.method]
		ctx.mw = n.mw
		return nil
	}

	if n.cascade {
		ctx.mw = append(ctx.mw, n.mw...)
	}

	// anything else i.e. {/}hello
	head, tail := ht(routeSegs)

	next := n.findMatchingChildWithCtx(ctx, head)
	if next == nil {
		return n
	}

	return next.traverse(ctx, tail)
}

func (n *node) insert(segments []RouteToken, method uint8, hf http.HandlerFunc, mw []http.Handler, cascade bool) {
	defer func() {
		n.leaf = len(n.children) == 0
	}()

	head, tail := ht(segments)
	if !n.matches(head) {
		logrus.Errorf("node.insert: could not traverse %v:%v with %v:%v", n.routeType, n.value, head.Type, head.Value)
		return
	}

	if len(tail) == 0 {
		if hf != nil {
			n.handlers[method] = hf
			n.isHandler = true
			logrus.Infof("node.insert: assigned handler %v to %v:%v", method, n.routeType, n.value)
		}

		if len(mw) != 0 {
			n.cascade = cascade
			n.mw = append(n.mw, mw...)
			logrus.Infof("node.insert: assigned %v mw to %v:%v", len(mw), n.routeType, n.value)
		}

		return
	}

	// at this stage, we know that the route is not fully consumed
	// and that we are in the right place

	curr := tail[0]

	next := n.findMatchingChild(curr)
	if next == nil {
		next = newNode(curr)
		n.addChild(next)
	}

	next.insert(tail, method, hf, mw, cascade)
}

func (n *node) findMatchingChild(toke RouteToken) *node {
	for _, child := range n.children {
		if child.matches(toke) {
			return child
		}
	}

	return nil
}

func (n *node) findMatchingChildWithCtx(ctx *zCtx, route string) *node {
	for _, c := range n.children {
		switch c.routeType {
		case Path:
			if c.value == route {
				return c
			}
		case Regex:
			if c.rgx.Match([]byte(route)) {
				return c
			}
		case Param:
			ctx.vars.Set(c.value, route)
			return c
		case WildCard:
			return c
		}
	}

	return nil
}

func (n *node) matches(toke RouteToken) bool {
	return n.value == toke.Value && n.routeType == toke.Type
}

func (n *node) addChild(child *node) {
	if len(n.children) == 0 {
		n.children = children{child}
		return
	}

	i := sort.Search(len(n.children), func(i int) bool {
		return n.children[i].routeType < n.routeType
	})

	n.children = slices.Insert(n.children, i, child)
}

func newNode(token RouteToken) *node {
	n := &node{
		routeType: token.Type,
		value:     token.Value,
		children:  children{},
		handlers:  newHandlers(),
		leaf:      true,
	}

	switch token.Type {
	case Path:
	case Param:
	case WildCard:
	case Regex:
		n.rgx = regexp.MustCompile(token.Value)
	}

	return n
}

func ht[S any](s []S) (S, []S) {
	if len(s) == 1 {
		return s[0], nil
	} else {
		return s[0], s[1:]
	}
}

func newHandlers() []http.HandlerFunc {
	hfs := make([]http.HandlerFunc, 9)
	return hfs
}
