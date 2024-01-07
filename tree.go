package zephyr

import (
	"net/http"
	"regexp"
	"slices"
	"sort"

	"github.com/sirupsen/logrus"
)

type Node struct {
	Type RouteType

	Rgx *regexp.Regexp

	Value string

	Leaf bool

	IsHandler bool

	Handlers []http.HandlerFunc

	Children
}

type Children []*Node

var MethodToIndexMap = map[string]uint8{
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

func newHandlers() []http.HandlerFunc {
	hfs := make([]http.HandlerFunc, 9)
	for i := range hfs {
		hfs[i] = func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	return hfs
}

func (n *Node) configureCtx(ctx *Context, r *http.Request) {
	path := r.URL.Path
}

func (n *Node) addRoute(segments []RouteToken, methodIndex uint8, hf http.HandlerFunc) {
	head, tail := ht(segments)
	if !n.matches(head) {
		logrus.Errorf("node.insertRoute: could not traverse %#v:%v with %v:%v", n.Type, n.Value, head.Type, head.Value)
		return
	}

	if len(tail) == 0 {
		n.Handlers[methodIndex] = hf
		logrus.Infof("node.insertRoute: assigned handler %v to %v:%v", methodIndex, n.Type, n.Value)
		return
	}

	// at this stage, we know that the route is not fully consumed
	// and that we are in the right place

	next := tail[0]

	child := n.findMatchingChild(next)
	if child == nil {
		child = newNode(next)
		n.addChild(child)
	}

	child.addRoute(tail, methodIndex, hf)

	n.Leaf = len(n.Children) == 0
}

func (n *Node) findMatchingChild(toke RouteToken) *Node {
	for _, child := range n.Children {
		if child.matches(toke) {
			return child
		}
	}

	return nil
}

func (n *Node) matches(toke RouteToken) bool {
	return n.Value == toke.Value && n.Type == toke.Type
}

func (n *Node) addChild(child *Node) {
	if len(n.Children) == 0 {
		n.Children = Children{child}
		return
	}

	i := sort.Search(len(n.Children), func(i int) bool {
		return n.Children[i].Type >= n.Type
	})

	n.Children = slices.Insert(n.Children, i, child)
}

func newNode(token RouteToken) *Node {
	n := &Node{
		Type:     token.Type,
		Value:    token.Value,
		Children: Children{},
		Handlers: newHandlers(),
		Leaf:     true,
	}

	switch token.Type {
	case Path:
	case Param:
	case WildCard:
	case Regex:
		n.Rgx = regexp.MustCompile(token.Value)
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
