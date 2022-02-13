package h

import (
	"fmt"
	"testing"
)

func TestH(t *testing.T) {
	tests := []struct {
		name string
		got  fmt.Stringer
		want string
	}{
		// General rendering
		{
			"renders elements",
			E("span"),
			`<span />`,
		},
		{
			"renders text",
			T{"aitch best templating library"},
			"aitch best templating library",
		},
		{
			"renders fragments",
			F(E("span"), T{"woohoo!!"}),
			"<span />woohoo!!",
		},

		// Attributes
		{
			"renders attributes",
			E("img", A{"src": "/path/image.png"}),
			`<img src="/path/image.png" />`,
		},
		{
			"renders true boolean attributes without value",
			E("div", A{"hidden": true}),
			`<div hidden />`,
		},
		{
			"does not render false boolean attributes",
			E("div", A{"hidden": false}),
			`<div />`,
		},
		{
			"renders multiple attributes from separate arguments",
			E("a", A{"aria-hidden": "true"}, A{"href": "#some-heading"}),
			`<a aria-hidden="true" href="#some-heading" />`,
		},
		{
			"merges consecutive class attributes",
			E("div", A{"class": "big"}, A{"class": "green"}),
			`<div class="big green" />`,
		},

		// Nesting
		{
			"renders child elements",
			E("div", E("div", E("div"))),
			`<div><div><div /></div></div>`,
		},
		{
			"renders child text",
			E("div", T{"Hello!!"}),
			`<div>Hello!!</div>`,
		},
		{
			"allows passing empty text to prevent element self-closing",
			E("div", T{}),
			`<div></div>`,
		},
		{
			"allows passing a fragment to prevent element self-closing",
			E("div", F()),
			`<div></div>`,
		},

		// Conditional rendering
		{
			"renders if true",
			E("div", If(true, func() N { return E("span") })),
			"<div><span /></div>",
		},
		{
			"renders if true (IfElse)",
			E("div", IfElse(true,
				func() N { return E("span") },
				func() N { return E("div") })),
			"<div><span /></div>",
		},
		{
			"does not render if false",
			E("div", If(false, func() N {
				return E("span")
			})),
			"<div></div>",
		},
		{
			"does not render if false (IfElse)",
			E("div", IfElse(false,
				func() N { return E("span") },
				func() N { return E("div") })),
			"<div><div /></div>",
		},
	}

	for _, test := range tests {
		got := Render(test.got)
		if got != test.want {
			t.Errorf("%s: got '%s', want '%s'", test.name, got, test.want)
		}
	}
}
