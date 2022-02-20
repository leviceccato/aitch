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

The primary building blocks of an Aitch template are Nodes. You can create them using the `h.E()` (element) function. Nodes can be HTML elements or fragments. A fragment is a Node that represents a list of HTML elements and does not itself get rendered. To create an element you must pass the tag of an HTML element as the first argument.

```go
h.E("div")
```

To create a fragment you must pass an empty string.

```go
h.E("")
```