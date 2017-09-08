package query

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWalker(t *testing.T) {
	p1 := gothichtml.NewTag("p")
	p1.AddChildren(gothichtml.NewText("paragraph 1"))
	p2 := gothichtml.NewTag("p")
	p2.AddChildren(gothichtml.NewText("paragraph 2"))
	p3 := gothichtml.NewTag("p")
	p3.AddChildren(gothichtml.NewText("paragraph 3"))
	div := gothichtml.NewTag("div")
	div.AddChildren(p1, gothichtml.NewFragment(p2, p3))
	body := gothichtml.NewTag("body")
	body.AddChildren(div, gothichtml.NewVoidTag("hr"))

	title := gothichtml.NewTag("title")
	title.AddChildren(gothichtml.NewText("This is a test"))
	head := gothichtml.NewTag("head")
	head.AddChildren(title)

	html := gothichtml.NewTag("html")
	html.AddChildren(head, body)
	root := gothichtml.NewFragment(gothichtml.NewDoctype("html"), html)

	assert.NotNil(t, root)

	expected := []struct {
		node gothichtml.Node
		loc  *Location
	}{
		{
			node: root,
			loc:  &Location{},
		},
		{
			node: root.Child(0),
			loc: &Location{
				Path: Path(0),
				Tag:  []int{0},
				Node: []int{0},
			},
		},
		{
			node: html,
			loc: &Location{
				Path: Path(1),
				Tag:  []int{0},
				Node: []int{1},
			},
		},
		{
			node: head,
			loc: &Location{
				Path: Path(1, 0),
				Tag:  []int{0, 0},
				Node: []int{1, 0},
			},
		},
		{
			node: title,
			loc: &Location{
				Path: Path(1, 0, 0),
				Tag:  []int{0, 0, 0},
				Node: []int{1, 0, 0},
			},
		},
		{
			node: title.Child(0),
			loc: &Location{
				Path: Path(1, 0, 0, 0),
				Tag:  []int{0, 0, 0, 0},
				Node: []int{1, 0, 0, 0},
			},
		},
		{
			node: body,
			loc: &Location{
				Path: Path(1, 1),
				Tag:  []int{0, 1},
				Node: []int{1, 1},
			},
		},
		{
			node: div,
			loc: &Location{
				Path: Path(1, 1, 0),
				Tag:  []int{0, 1, 0},
				Node: []int{1, 1, 0},
			},
		},
		{
			node: p1,
			loc: &Location{
				Path: Path(1, 1, 0, 0),
				Tag:  []int{0, 1, 0, 0},
				Node: []int{1, 1, 0, 0},
			},
		},
		{
			node: p1.Child(0),
			loc: &Location{
				Path: Path(1, 1, 0, 0, 0),
				Tag:  []int{0, 1, 0, 0, 0},
				Node: []int{1, 1, 0, 0, 0},
			},
		},
		{
			node: div.Child(1),
			loc: &Location{
				Path: Path(1, 1, 0, 1),
				Tag:  []int{0, 1, 0, 1},
				Node: []int{1, 1, 0, 1},
			},
		},
		{
			node: p2,
			loc: &Location{
				Path: Path(1, 1, 0, 1, 0),
				Tag:  []int{0, 1, 0, 1},
				Node: []int{1, 1, 0, 1},
			},
		},
		{
			node: p2.Child(0),
			loc: &Location{
				Path: Path(1, 1, 0, 1, 0, 0),
				Tag:  []int{0, 1, 0, 1, 0},
				Node: []int{1, 1, 0, 1, 0},
			},
		},
		{
			node: p3,
			loc: &Location{
				Path: Path(1, 1, 0, 1, 1),
				Tag:  []int{0, 1, 0, 2},
				Node: []int{1, 1, 0, 2},
			},
		},
		{
			node: p3.Child(0),
			loc: &Location{
				Path: Path(1, 1, 0, 1, 1, 0),
				Tag:  []int{0, 1, 0, 2, 0},
				Node: []int{1, 1, 0, 2, 0},
			},
		},
		{
			node: body.Child(1),
			loc: &Location{
				Path: Path(1, 1, 1),
				Tag:  []int{0, 1, 1},
				Node: []int{1, 1, 1},
			},
		},
	}
	var i int
	visiter := func(node gothichtml.Node, location *Location) {
		if i >= len(expected) {
			t.Error("Too many nodes")
			return
		}
		e := expected[i]
		i++
		assert.Equal(t, e.node, node)
		assert.EqualValues(t, e.loc, location)
	}

	Walk(root, visiter)
}
