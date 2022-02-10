package h

import "fmt"

func init() {

	html := E("test",
		A{"test": true},
		E("div", A{}),
	)

	fmt.Println(html)
}

func Render(data ...NodeData) string {
	return ""
}

type Node struct {
	tag        string
	attributes A
	content    []Renderer
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

func (n Node) Render() string {
	result := "<" + n.tag
	for name, attribute := range n.attributes {
		result += parseAttribute(name, attribute)
	}
	if len(n.content) == 0 {
		return result + " />"
	}
	for _, renderer := range n.content {
		result += renderer.Render()
	}
	return result
}

func (n Node) addToNode(node *Node) {
	node.content = append(node.content, n)
}

type Renderer interface {
	Render() string
}

// Fragment
type NodeData interface {
	addToNode(*Node)
}

// Attributes
type A map[string]interface{}

func (a A) addToNode(node *Node) {
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

type Comment struct {
	content string
}

func (c Comment) Render() string {
	return "<!-- " + c.content + " -->"
}

func (c Comment) addToNode(node *Node) {
	node.content = append(node.content, c)
}

// Comment
func C(content string) NodeData {
	return Comment{content: content}
}

// Element
func E(selector string, data ...NodeData) NodeData {
	node := parseSelector(selector)

	for _, datum := range data {
		datum.addToNode(&node)
	}

	return node
}

func parseSelector(selector string) Node {
	n := Node{tag: selector}
	attrs := A{}
	attrs.addToNode(&n)

	return n
}
