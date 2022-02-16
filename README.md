# Aitch

_This library is a work in progress and not currently ready for production._

Aitch (or Haitch depending on your region) is a pragmatic HTML templating library for Go. It takes inspiration from the JavaScript framework Mithril which in turn takes inspiration from hyperscript. Unlike traditional templating engines, Aitch and other hyperscript-like libraries produce HTML through nested functions.

This library is designed to output HTML strings and is therefore not recommended for use with GopherJS.

## Features
- Polymorphic and variadic `h.E` function that can produce HTML with great flexibility
- Control flow with `h.If`, `h.IfElse` and `h.For`
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

func Page(title string, users []int) string {
    return h.F(
        h.E("!DOCTYPE[html]"),
        h.E("html",
            h.E("head",
                h.E("title",
                    h.IfElse(title == "",
                        func() h.D { return "Page" },
                        func() h.D { return h.T{title} },
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