package routing

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"sort"
)

type Node struct {
	Type RouteType

	Rgx *regexp.Regexp

	Value string

	Leaf bool

	IsHandler bool

	Handler http.HandlerFunc

	Parent *Node
	Children
}

type Children []*Node

func (c Children) findForInsert(n *Node) *Node {
	for _, child := range c {
		if n.identical(child) {
			return child
		}
	}

	return nil
}

func (c Children) append(n *Node) Children {
	if len(c) == 0 {
		return Children{n}
	}
	i := sort.Search(len(c), func(i int) bool {
		if c[i].Type >= n.Type {
			return true
		}
		return false
	})

	return slices.Insert(c, i, n)
}

func (parent *Node) insert(segments []RouteToken) {
	depth := len(segments)
	if depth == 0 {
		return
	}

	// if we get here, we're no longer going to be a leaf
	parent.Leaf = false

	head, tail := ht(segments)

	nn := newNode(head, parent)

	match := parent.Children.findForInsert(nn)
	if match == nil { // without a match, we can just add the new node
		parent.Children = parent.Children.append(nn)
		nn.insert(tail)
	} else { // we can try insert the rest of the route
		match.Leaf = false
		match.insert(tail)
	}
}

func (n *Node) identical(other *Node) bool {
	return n.Type == other.Type && n.Value == other.Value
}

func (n *Node) walk(print bool) []string {
	routes := []string{}
	for _, child := range n.Children {
		if child.Type == WildCard {
		}
		child.walk(print)
	}

	if n.Leaf || n.IsHandler {
		fullPath := n.pathFromRoot("")
		if print {
			fmt.Printf("Got: %v\n", fullPath)
		}
		routes = append(routes, fullPath)
	}

	return routes
}

func (n *Node) pathFromRoot(path string) string {
	var addition string
	switch n.Type {
	case WildCard:
		addition = "/*"
	case Regex:
		addition = fmt.Sprintf("/~%v", n.Value)
	case Param:
		addition = fmt.Sprintf("/{%v}", n.Value)
	case Path:
		addition = fmt.Sprintf("/%v", n.Value)
	}

	addition += path
	if n.Parent != nil {
		path = n.Parent.pathFromRoot(addition)
	}
	return path
}

func (n *Node) configureCtx(ctx *Context, routeSegs []string) ([]string, *Node) {
	if len(routeSegs) == 0 || n.Leaf {
		ctx.Handler = n.Handler
		return nil, nil
	}

	var (
		foundNode    *Node
		currentRoute = routeSegs[0]
	)

	for _, child := range n.Children {
		if foundNode != nil {
			break
		}

		switch child.Type {
		case Path:
			if child.Value == currentRoute {
				foundNode = child
			}
		case WildCard:
			foundNode = child
		case Regex:
			if child.Rgx.Match([]byte(currentRoute)) {
				foundNode = child
			}
		case Param:
			ctx.Vars.Set(child.Value, currentRoute)
			foundNode = child
		}
	}

	return foundNode.configureCtx(ctx, routeSegs[1:])
}

func newNode(token RouteToken, parent *Node) *Node {
	n := &Node{
		Type:     token.Type,
		Value:    token.Value,
		Children: Children{},
		Leaf:     true,
		Parent:   parent,
	}

	switch token.Type {
	case Path:
	case Param:
	case WildCard:
	case Regex:
		n.Rgx = regexp.MustCompile(token.Value)
	}

	n.Handler = func(w http.ResponseWriter, r *http.Request) {
		if n.Type == Param {
			w.Write([]byte(fmt.Sprintf("here is the var! %v\n", Vars(r.Context(), n.Value))))
		}
		w.Write([]byte(fmt.Sprintf("hello from %v", n.pathFromRoot(""))))
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
