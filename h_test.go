package h

import (
	"fmt"
	"testing"
)

// Make tests less verbose
type w []string

// Render function testing utility
func hTest(t *testing.T, node fmt.Stringer, wants w) {
	got := node.String()

	hasMatch := false

	for _, want := range wants {
		if got == want {
			hasMatch = true
		}
	}

	if !hasMatch {
		t.Errorf("got '%s', want one of %v", got, wants)
	}
}

// General rendering

func TestRenderElement(t *testing.T) {
	hTest(t,
		E("span"),
		w{`<span />`},
	)
}

func TestRenderRawHTML(t *testing.T) {
	hTest(t,
		R("aitch best templating <div />"),
		w{"aitch best templating <div />"},
	)
}

func TestRenderHTMLEscapedText(t *testing.T) {
	hTest(t,
		T("hi! <div>Hmmm</div>"),
		w{"hi! &lt;div&gt;Hmmm&lt;/div&gt;"},
	)
}

func TestRenderFragment(t *testing.T) {
	hTest(t,
		F(E("span"), T("woohoo!!")),
		w{"<span />woohoo!!"},
	)
}

func TestRenderComment(t *testing.T) {
	hTest(t,
		E("span", C("woohoo!!")),
		w{"<span><!-- woohoo!! --></span>"},
	)
}

// Attributes

func TestRenderAttribute(t *testing.T) {
	hTest(t,
		E("img", A{"src": "/path/image.png"}),
		w{`<img src="/path/image.png" />`},
	)
}

func TestRenderTrueAttribute(t *testing.T) {
	hTest(t,
		E("div", A{"hidden": true}),
		w{`<div hidden />`},
	)
}

func TestRenderFalseAttribute(t *testing.T) {
	hTest(t,
		E("div", A{"hidden": false}),
		w{`<div />`},
	)
}

func TestRenderMultipleAttributesSeparateArguments(t *testing.T) {
	hTest(t,
		E("a", A{"aria-hidden": "true"}, A{"href": "#some-heading"}),
		w{
			`<a aria-hidden="true" href="#some-heading" />`,
			`<a href="#some-heading" aria-hidden="true" />`,
		},
	)
}

func TestMergeConsecutiveClassAttributes(t *testing.T) {
	hTest(t,
		E("div", A{"class": "big"}, A{"class": "green"}),
		w{`<div class="big green" />`},
	)
}

// Selectors

func TestClassSelector(t *testing.T) {
	hTest(t,
		E("div.test"),
		w{`<div class="test" />`},
	)
}

func TestIdSelector(t *testing.T) {
	hTest(t,
		E("div#test"),
		w{`<div id="test" />`},
	)
}

func TestDefaultSelectorTag(t *testing.T) {
	hTest(t,
		E(".test"),
		w{`<div class="test" />`},
	)
}

func TestOnlyAddFirstId(t *testing.T) {
	hTest(t,
		E("div#test#test2"),
		w{`<div id="test" />`},
	)
}

func TestIdAndClassSelector(t *testing.T) {
	hTest(t,
		E("div#test.test2"),
		w{
			`<div id="test" class="test2" />`,
			`<div class="test2" id="test" />`,
		},
	)
}

func TestTagFromComplexSelector(t *testing.T) {
	hTest(t,
		E("span.test2"),
		w{`<span class="test2" />`},
	)
}

func TestMultipleClassesSelector(t *testing.T) {
	hTest(t,
		E("div.test.test2"),
		w{
			`<div class="test test2" />`,
			`<div class="test2 test" />`,
		},
	)
}

func TestClassesSelectorAndClassAttribute(t *testing.T) {
	hTest(t,
		E("div.test", A{"class": "test2"}),
		w{
			`<div class="test test2" />`,
			`<div class="test2 test" />`,
		},
	)
}

func TestRenderCustomAttribute(t *testing.T) {
	hTest(t,
		E(`div[test="thing"]`),
		w{`<div test="thing" />`},
	)
}

func TestRenderMultipleCustomAttribute(t *testing.T) {
	hTest(t,
		E(`div[test="thing"][test2="thing"]`),
		w{
			`<div test="thing" test2="thing" />`,
			`<div test2="thing" test="thing" />`,
		},
	)
}

func TestCompactsSelector(t *testing.T) {
	hTest(t,
		E(`div
			.test
			[test2="thing"]
		`),
		w{
			`<div class="test" test2="thing" />`,
			`<div test2="thing" class="test" />`,
		},
	)
}

func TestRenderCustomBooleanAttribute(t *testing.T) {
	hTest(t,
		E(`div[test]`),
		w{`<div test />`},
	)
}

// Nesting

func TestRenderChildElements(t *testing.T) {
	hTest(t,
		E("div", E("div", E("div"))),
		w{`<div><div><div /></div></div>`},
	)
}

func TestRenderChildText(t *testing.T) {
	hTest(t,
		E("div", T("Hello!!")),
		w{`<div>Hello!!</div>`},
	)
}

func TestPreventSelfClosingWithText(t *testing.T) {
	hTest(t,
		E("div", T("")),
		w{`<div></div>`},
	)
}

func TestPreventSelfClosingWithFragment(t *testing.T) {
	hTest(t,
		E("div", F()),
		w{`<div></div>`},
	)
}

// Conditional rendering

func TestRenderIfTrue(t *testing.T) {
	hTest(t,
		E("div", If(true, E("span"))),
		w{`<div><span /></div>`},
	)
}

func TestIfElseRenderIfTrue(t *testing.T) {
	hTest(t,
		E("div", IfElse(true,
			E("span"),
			E("div"),
		)),
		w{`<div><span /></div>`},
	)
}

func TestDontRenderIfFalse(t *testing.T) {
	hTest(t,
		E("div", If(false, E("span"))),
		w{`<div></div>`},
	)
}

func TestIfElseDontRenderIfFalse(t *testing.T) {
	hTest(t,
		E("div", IfElse(false,
			E("span"),
			E("div"),
		)),
		w{`<div><div /></div>`},
	)
}

func TestFor(t *testing.T) {
	names := []string{"Jon", "Lawrie", "Jade"}

	hTest(t,
		E("div", For(names, func(_ int, name string) D {
			return E("span", T(name))
		})),
		w{`<div><span>Jon</span><span>Lawrie</span><span>Jade</span></div>`},
	)
}

func TestExampleHTML(t *testing.T) {
	imgSrc := "/logo.png"

	hTest(t,
		F(
			E("!DOCTYPE[html]"),
			E(`html[lang="en"]`,
				E("head",
					E("title", T("Example HTML")),
				),
				E("body",
					E("header",
						E(".container",
							E("img", A{
								"src": imgSrc,
								"alt": "",
							}),
						),
					),
					E("main",
						E("div",
							E("h1#title",
								T("Title goes here"),
							),
						),
					),
					E("footer",
						E(`a[href="/"]`,
							T("Home"),
						),
					),
				),
			),
		),
		w{
			`<!DOCTYPE html /><html lang="en"><head><title>Example HTML</title></head><body><header><div class="container"><img src="/logo.png" alt="" /></div></header><main><div><h1 id="title">Title goes here</h1></div></main><footer><a href="/">Home</a></footer></body></html>`,
			`<!DOCTYPE html /><html lang="en"><head><title>Example HTML</title></head><body><header><div class="container"><img alt="" src="/logo.png" /></div></header><main><div><h1 id="title">Title goes here</h1></div></main><footer><a href="/">Home</a></footer></body></html>`,
		},
	)
}
