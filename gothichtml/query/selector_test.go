package query

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/gothic/gothichtml/parsehtml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatches(t *testing.T) {
	tag := gothichtml.NewTag("div", "class", "foo glorp", "id", "testing")
	assert.True(t, checkTag("div").check(tag))
	assert.False(t, checkTag("ul").check(tag))
	assert.True(t, checkClass("foo").check(tag))
	assert.True(t, checkClass("glorp").check(tag))
	assert.False(t, checkClass("bar").check(tag))
	assert.True(t, checkID("testing").check(tag))
	assert.False(t, checkID("test").check(tag))
}

func TestSelect(t *testing.T) {
	i := parsehtml.Must("<i>bar</i>")
	p1 := parsehtml.Must("<p>paragraph 1</p>")
	p2 := parsehtml.Must(`<p class="foo">paragraph 2</p>`)
	p3 := parsehtml.Must(`<p class="foo">paragraph 3</p>`)
	div := gothichtml.NewTag("div")
	div.AddChildren(i, p1, gothichtml.NewText("interrupt"), gothichtml.NewFragment(p2, p3))
	p4 := gothichtml.NewTag("p", "class", "foo")
	p4.AddChildren(gothichtml.NewText("paragraph 4"))
	body := gothichtml.NewTag("body")
	body.AddChildren(div, gothichtml.NewVoidTag("hr"), p4)

	title := parsehtml.Must(`<title>This is a test</title>`)
	head := gothichtml.NewTag("head")
	head.AddChildren(title)

	html := gothichtml.NewTag("html")
	html.AddChildren(head, body)
	root := gothichtml.NewFragment(gothichtml.NewDoctype("html"), html)

	s := selectors{
		&selector{
			checkers: []nodeChecker{checkTag("div")},
		},
	}
	matches := s.QueryAll(root)
	if assert.Len(t, matches, 1) {
		assert.Equal(t, div, matches[0])
	}

	s = selectors{
		&selector{
			checkers: []nodeChecker{checkTag("p"), checkClass("foo")},
		},
	}
	matches = s.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p2, matches[0])
		assert.Equal(t, p3, matches[1])
		assert.Equal(t, p4, matches[2])
	}

	s = selectors{
		&selector{
			checkers: []nodeChecker{checkTag("div")},
			next: &selector{
				checkers: []nodeChecker{checkTag("p"), checkClass("foo")},
			},
			nextLoc: newDescendant,
		},
	}

	sel, err := Selector("div p.foo, title")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, title, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}

	sel, err = Selector("div>p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p1, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}

	sel, err = Selector("i+p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 1) {
		assert.Equal(t, p1, matches[0])
	}

	sel, err = Selector("i~p")
	assert.NoError(t, err)
	assert.NotNil(t, sel)
	matches = sel.QueryAll(root)
	if assert.Len(t, matches, 3) {
		assert.Equal(t, p1, matches[0])
		assert.Equal(t, p2, matches[1])
		assert.Equal(t, p3, matches[2])
	}
}
