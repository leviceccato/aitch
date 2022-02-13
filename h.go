package h

import (
	"bytes"
	"fmt"
)

func Render(stringer fmt.Stringer) string {
	return stringer.String()
}

// Node
type N struct {
	tag        string
	attributes A
	content    []fmt.Stringer
}

func parseAttribute(name string, attribute interface{}) string {
	switch a := attribute.(type) {
	case bool:
		if a {
			return " " + name
		}
		return ""
	default:
		return fmt.Sprintf(" %v=\"%v\"", name, a)
	}
}

func (n N) String() string {
	isElement := n.tag != ""
	var b bytes.Buffer

	if isElement {
		b.WriteString("<" + n.tag)

		for name, attribute := range n.attributes {
			b.WriteString(parseAttribute(name, attribute))
		}

		if len(n.content) == 0 {
			b.WriteString(" />")
			return b.String()
		}
		b.WriteString(">")
	}

	for _, renderer := range n.content {
		b.WriteString(renderer.String())
	}

	if isElement {
		b.WriteString("</" + n.tag + ">")
	}

	return b.String()
}

func (n N) addToNode(node *N) {
	node.content = append(node.content, n)
}

// Fragment
type NodeData interface {
	addToNode(*N)
}

// Attributes
type A map[string]interface{}

func (a A) addToNode(node *N) {
	if node.attributes == nil {
		node.attributes = A{}
	}

	for key, value := range a {
		if value == nil {
			continue
		}

		if key == "class" {
			class, ok := node.attributes["class"]
			if !ok {
				node.attributes["class"] = fmt.Sprintf("%v", value)
				continue
			}
			node.attributes["class"] = fmt.Sprintf("%v %v", class, value)
			continue
		}

		node.attributes[key] = value
	}
}

type T [1]string

func (t T) String() string {
	return t[0]
}

func (t T) addToNode(node *N) {
	node.content = append(node.content, t)
}

type Comment [1]string

func (c Comment) String() string {
	return "<!-- " + c[0] + " -->"
}

func (c Comment) addToNode(node *N) {
	node.content = append(node.content, c)
}

// Comment
func C(content string) NodeData {
	return Comment{content}
}

// Element
func E(selector string, data ...NodeData) N {
	node := parseSelector(selector)

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

// Fragment
func F(data ...NodeData) N {
	node := N{}

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

func parseSelector(selector string) N {
	node := N{tag: selector}

	attrs := A{}
	attrs.addToNode(&node)

	return node
}

type CondFunc func() N

func If(cond bool, fn CondFunc) N {
	if !cond {
		return F()
	}
	return fn()
}

func IfElse(cond bool, fnIf CondFunc, fnElse CondFunc) N {
	if !cond {
		return fnElse()
	}
	return fnIf()
}

// Make this generic in 1.18
func For(items []interface{}, fn func(item interface{}, index int) N) N {
	node := F()
	for index, item := range items {
		n := fn(item, index)
		node.addToNode(&n)
	}
	return node
}
