package nodeutil

import (
	"golang.org/x/net/html"
)

//FindOne looks up nodes until the end of the node tree and return the first match
func FindOne(n *html.Node, f func(*html.Node) bool) (*html.Node, bool) {
	if f(n) {
		return n, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, ok := FindOne(c, f); ok {
			return found, ok
		}
	}
	return nil, false
}

//Find looks up nodes until the end of the node tree
func Find(n *html.Node, f func(*html.Node) bool) []*html.Node {
	found := []*html.Node{}
	if f(n) {
		found = append(found, n)

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found = append(found, Find(c, f)...)
	}
	return found
}

//FindOneLimited Looks up a single node up to a depth of depth
func FindOneLimited(n *html.Node, f func(*html.Node) bool, depth int) (*html.Node, bool) {
	if depth <= 0 {
		return nil, false
	}

	if f(n) {
		return n, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found, ok := FindOneLimited(c, f, depth-1); ok {
			return found, ok
		}
	}
	return nil, false
}

//FindLimited Looks up node up to a depth of depth
func FindLimited(n *html.Node, f func(*html.Node) bool, depth int) []*html.Node {
	found := []*html.Node{}
	if depth <= 0 {
		return found
	}

	if f(n) {
		found = append(found, n)

	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found = append(found, FindLimited(c, f, depth-1)...)
	}
	return found
}

//HasAttr returns true when the key value pair is found in the attribute.
func HasAttr(attrs []html.Attribute, key, val string) bool {
	for _, a := range attrs {
		if a.Key == key && a.Val == val {
			return true
		}
	}
	return false
}

// type Term struct {
// 	Type  html.NodeType
// 	Data  string
// 	Attrs []html.Attribute
// }

// func SearchPath() []*html.Node {

// }
