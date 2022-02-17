# Aitch

_This library is a work in progress and not currently ready for production._

Aitch (or Haitch) is a HTML templating library for Go. It takes inspiration from the JavaScript framework Mithril. Unlike traditional templating engines, Aitch produces HTML with nested functions in a style popularized by hyperscript.

This library is designed to output HTML strings and is therefore not recommended for use with GopherJS.

## Features
- Polymorphic and variadic `E` function that can produce HTML with great flexibility
- Control flow with `If`, `IfElse` and `For`
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
    return h.F(
        h.E("!DOCTYPE[html]"),
        h.E("html",
            h.E("head",
                h.E("title",
                    h.IfElse(title == "",
                        "Page",
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