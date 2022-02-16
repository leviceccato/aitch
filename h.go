package h

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"unicode"
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

func renderAttribute(name string, attribute any) string {
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
			b.WriteString(renderAttribute(name, attribute))
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
type D interface {
	addToNode(*N)
}

// Attributes
type A map[string]any

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

// May be instantiated like [T]{"string"}
type wrappedString struct {
	content string
}

// Text
type T wrappedString

func (t T) String() string {
	return html.EscapeString(t.content)
}

func (t T) addToNode(node *N) {
	node.content = append(node.content, t)
}

// Raw HTML
type R wrappedString

func (r R) String() string {
	return r.content
}

func (r R) addToNode(node *N) {
	node.content = append(node.content, r)
}

type Comment wrappedString

func (c Comment) String() string {
	return "<!-- " + c.content + " -->"
}

func (c Comment) addToNode(node *N) {
	node.content = append(node.content, c)
}

// Comment
func C(content string) D {
	return Comment{content}
}

// Element
func E(selector string, data ...D) N {
	node := newNode(selector)

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

// Fragment
func F(data ...D) N {
	node := N{}

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

func newNode(selector string) N {
	tag, attrs := parseSelector(selector)

	node := N{tag: tag}
	for _, attr := range attrs {
		attr.addToNode(&node)
	}

	return node
}

func parseSelector(selector string) (string, []A) {
	chars := []rune(selector)
	charsLength := len(chars)

	attrs := []A{}
	tag := "div"
	segmentEnd := charsLength
	isCustom := false

	for i := charsLength - 1; i >= 0; i-- {
		char := chars[i]

		// Close out custom attribute
		if char == '[' {
			attrs = append(attrs, parseAttribute(string(chars[i+1:segmentEnd])))
			isCustom = false
			segmentEnd = i
			continue
		}

		// Ignore if currently inside a custom attribute
		if isCustom {
			continue
		}

		// Start custom attribute
		if char == ']' {
			segmentEnd = i
			isCustom = true
			continue
		}

		// Close out id
		if char == '#' {
			attrs = append(attrs, A{"id": compactSelector(string(chars[i+1 : segmentEnd]))})
			segmentEnd = i
			continue
		}

		// Close out class
		if char == '.' {
			attrs = append(attrs, A{"class": compactSelector(string(chars[i+1 : segmentEnd]))})
			segmentEnd = i
			continue
		}

		// Starts with a tag, use it for the node
		if i == 0 {
			tag = compactSelector(string(chars[i:segmentEnd]))
		}
	}

	return tag, attrs
}

func compactSelector(selector string) string {
	var b bytes.Buffer
	b.Grow(len(selector))

	for _, char := range selector {
		if unicode.IsSpace(char) {
			continue
		}
		b.WriteRune(char)
	}

	return b.String()
}

func parseAttribute(attrPair string) A {
	entries := strings.FieldsFunc(attrPair, func(char rune) bool {
		if char == '=' {
			return true
		}
		return false
	})

	name := compactSelector(entries[0])

	if len(entries) == 1 {
		return A{name: true}
	}

	value := strings.TrimFunc(entries[1], func(char rune) bool {
		if char == '"' || char == '\'' {
			return true
		}
		return false
	})

	return A{name: value}
}

func If(cond bool, fn func() D) D {
	if !cond {
		return F()
	}
	return fn()
}

func IfElse(cond bool, fnIf func() D, fnElse func() D) D {
	if !cond {
		return fnElse()
	}
	return fnIf()
}

func For[I any](items []I, fn func(index int, item I) D) D {
	node := F()

	for index, item := range items {
		n := fn(index, item)
		n.addToNode(&node)
	}

	return node
}
