# Aitch

_This library is a work in progress and not currently ready for production._

Aitch (or Haitch) is a HTML templating library for Go. It takes inspiration from the JavaScript framework Mithril. Unlike traditional templating engines, Aitch produces HTML with nested functions in a style popularized by hyperscript.

This library is designed to output HTML strings and is therefore not recommended for use with GopherJS.

## Features
- Polymorphic and variadic `E` function that can produce HTML with great flexibility
- Control flow with `If`, `Else`, `ElseIf` and `For`
- CSS style attribute definitions
- Smart handling of classes and boolean attributes
- Fragments, raw HTML, text and comments

## Requirements

- Go 1.18

## Installation

```
go get github.com/leviceccato/aitch
```

## Example

```go
package component

import (
    "github.com/leviceccato/aitch"
)

func Page(title string, users []string) string {
    return h.E("",
        h.E("!DOCTYPE[html]"),
        h.E("html",
            h.E("head",
                h.E("title",
                    h.If(title == "",
                        h.T{"Page"},
                    ).Else(
                        h.T{title},
                    ),
                ),
                h.E("meta", h.A{"charset": "utf-8"}),
            ),
            h.E("body#body",
                h.E(".container",
                    h.For(users, func(_ int, user string) h.D {
                        return h.E("p", h.T{"User: " + user})
                    }),
                ),
            ),
        ),
    ).String()
}
```

## Usage

Most of Aitch's functions are shortened to be a single letter. This helps the templates remain terse and readable as they grow larger and the functions are repeated.

### Nodes

The primary building blocks of an Aitch template are nodes. You can create them using the `h.E()` (element) function. Nodes can be HTML elements or fragments. A fragment is a node that represents a list of HTML elements and does not itself get rendered. To create an element you must pass the tag of an HTML element as the first argument, this can be any string and is not limited to valid HTML elements.

```go
h.E("div")
```

To create a fragment you must pass an empty string.

```go
h.E("")
```

Nodes can have any number of child nodes, they may be passed as extra arguments to the `h.E()` function.

```go
h.E("header",
    h.E("img"),
    h.E("nav"),
)
```

### Attributes

Attributes can be added to elements using `h.A{}` structs. They can be passed as extra arguments to the `h.E()` function in the same way as nodes. Attributes passed to fragments will have no effect. Note that attributes are not required to come before node arguments.

```go
h.E("img", h.A{
    "src": "/images/aitch-logo.png",
    "alt": "Aitch",
})
```

Class and style attributes keys have special handling. Classes may be passed as a string and they will behave as expected, but you can also pass a nested `h.A{}` struct with keys as classes and values as booleans to indicate if that class is active. Consecutive class definitions will be merged with the previous.

```go
h.E("div", h.A{
    "class": h.A{
        "big": false,
        "red": true,
    },
}, h.A{
    "class": h.A{
        "another": true,
    },
})
```

Styles may be passed as a `h.A{}` struct with keys as CSS properties and values as CSS values. Note that these values must be strings.

```go
h.E("div", h.A{
    "style": h.A{
        "color": "red",
        "margin": "0",
    }
})
```

### Complex selectors

The `h.E()` function can be passed plain element tags, but it also accepts complex CSS-style selectors. Classes can be added by appending a period followed by the class name. IDs can be added by appending a number sign followed by the ID. Both of these syntaxes can be chained.

```go
h.E("div.card.big#card-1")
```

If no element is provided at the start of the string then a DIV element is assumed.

```go
h.E(".container")
```

Attributes besides class and id can be added as well. They must be surrounded by square brackets in the same way as CSS selectors.

```go
h.E(`
    a.link
    [href="/"]
    [target="_blank"]
`)
```
