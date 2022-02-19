package h

import (
	"bytes"
	"fmt"
	"html"
	"strings"
	"unicode"
)

// Node
type Node struct {
	tag     string
	attrs   A
	content []fmt.Stringer
}

func (n Node) Else(data ...D) Node {
	return n.ElseIf(true, data...)
}

func (n Node) ElseIf(cond bool, data ...D) Node {
	if !cond && (len(n.content) > 0) {
		return E("")
	}
	return E("", data...)
}

func (n Node) String() string {
	isElement := n.tag != ""

	b := bytes.Buffer{}

	if isElement {
		b.WriteString("<" + n.tag + n.attrs.String())

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

func (n Node) addToNode(node *Node) {
	node.content = append(node.content, n)
}

// Data
type D interface {
	addToNode(*Node)
}

// Attributes
type A map[string]any

func (a A) String() string {
	b := bytes.Buffer{}

	for name, attr := range a {
		switch name {
		case "class":
			b.WriteString(" class=\"")

			for name, isActive := range attr.(A) {
				if !isActive.(bool) {
					continue
				}
				b.WriteString(" " + name)
			}

			b.WriteString("\"")

		case "style":
			b.WriteString(" style=\"")

			for prop, value := range attr.(A) {
				b.WriteString(" " + prop + ": " + value.(string) + ";")
			}

			b.WriteString("\"")
		}

		switch attr.(type) {
		case bool:
			if !attr.(bool) {
				continue
			}
			b.WriteString(" " + name)

		case string:
			b.WriteString(" " + name + "=\"" + attr.(string) + "\"")
		}
	}

	return b.String()
}

func (a A) addToNode(n *Node) {
	if n.attrs == nil {
		n.attrs = A{}
	}

	for key, value := range a {
		if value == nil {
			continue
		}

		switch key {
		case "class":
			_, ok := n.attrs["class"]
			if !ok {
				n.attrs["class"] = A{}
			}

			switch value.(type) {
			case A:
				for name, isActive := range value.(A) {
					n.attrs["class"].(A)[name] = isActive.(bool)
				}

			case string:
				names := strings.FieldsFunc(value.(string), unicode.IsSpace)
				for _, name := range names {
					n.attrs["class"].(A)[name] = true
				}
			}

			continue

		case "style":
			_, ok := n.attrs["style"]
			if !ok {
				n.attrs["style"] = A{}
			}

			switch value.(type) {
			case A:
				for prop, propVal := range value.(A) {
					n.attrs["style"].(A)[prop] = propVal.(string)
				}

			case string:
				styles := strings.FieldsFunc(value.(string), func(char rune) bool {
					return char == ';'
				})

				for _, style := range styles {
					propAndValue := strings.FieldsFunc(style, func(char rune) bool {
						return char == ':'
					})

					// Not formatted correctly, abandon
					if len(propAndValue) < 2 {
						break
					}

					n.attrs["class"].(A)[propAndValue[0]] = propAndValue[1]
				}
			}

			continue
		}

		switch value.(type) {
		case bool:
			v := value.(bool)
			if v {
				n.attrs[key] = v
				continue
			}
			continue

		case string:
			n.attrs[key] = value.(string)
			continue
		}

		n.attrs[key] = fmt.Sprintf("%v", value)
	}
}

// Text
type T string

func (t T) String() string {
	return html.EscapeString(string(t))
}

func (t T) addToNode(node *Node) {
	node.content = append(node.content, t)
}

// Raw HTML
type R string

func (r R) String() string {
	return string(r)
}

func (r R) addToNode(node *Node) {
	node.content = append(node.content, r)
}

// Comment
type C string

func (c C) String() string {
	return "<!-- " + string(c) + " -->"
}

func (c C) addToNode(node *Node) {
	node.content = append(node.content, c)
}

// Element
func E(selector string, data ...D) Node {
	node := newNode(selector)

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

func newNode(selector string) Node {
	tag, attrs := parseSelector(selector)

	node := Node{tag: tag}
	for _, attr := range attrs {
		attr.addToNode(&node)
	}

	return node
}

func parseSelector(selector string) (string, []A) {
	attrs := []A{}
	if selector == "" {
		return "", attrs
	}

	chars := []rune(selector)
	charsLength := len(chars)

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
			attrs = append(attrs, A{"id": compactStr(string(chars[i+1 : segmentEnd]))})
			segmentEnd = i
			continue
		}

		// Close out class
		if char == '.' {
			attrs = append(attrs, A{"class": compactStr(string(chars[i+1 : segmentEnd]))})
			segmentEnd = i
			continue
		}

		// Starts with a tag, use it for the node
		if i == 0 {
			tag = compactStr(string(chars[i:segmentEnd]))
		}
	}

	return tag, attrs
}

func compactStr(str string) string {
	b := bytes.Buffer{}
	b.Grow(len(str))

	for _, char := range str {
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

	name := compactStr(entries[0])

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

// Control flow

func If(cond bool, data ...D) Node {
	if !cond {
		return E("")
	}
	return E("", data...)
}

func For[I any](items []I, fn func(index int, item I) D) D {
	node := E("")

	for index, item := range items {
		n := fn(index, item)
		n.addToNode(&node)
	}

	return node
}
