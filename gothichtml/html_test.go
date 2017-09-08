package gothichtml

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHtml(t *testing.T) {
	div := NewTag("div", "id", "test")
	assert.NotNil(t, div)
	expected := `<div id="test"></div>`
	assert.Equal(t, expected, String(div))

	var asContainer ContainerNode
	asContainer = div
	assert.NotNil(t, asContainer)

	var asTag TagNode
	asTag = div
	assert.NotNil(t, asTag)
}

func TestFormat(t *testing.T) {
	li1 := NewTag("li")
	li1.AddChildren(NewText("Item 1"))
	li2 := NewTag("li")
	li2.AddChildren(NewText("Item 2"))
	li3 := NewTag("li")
	li3.AddChildren(NewText("Item 3"))
	ul := NewTag("ul")
	ul.AddChildren(li1, li2, li3)
	p1 := NewTag("p")
	p1.AddChildren(NewText("This is a test"))
	p2 := NewTag("p")
	p2.AddChildren(NewText("This is another test"))
	p3 := NewTag("p")
	p3.AddChildren(NewText("That's right, 3 paragraphs."))
	div := NewTag("div")
	div.AddChildren(p1, ul, p2)
	f := NewFragment(div, p3)
	expected := `<div>
  <p>This is a test</p>
  <ul>
    <li>Item 1</li>
    <li>Item 2</li>
    <li>Item 3</li>
  </ul>
  <p>This is another test</p>
</div>
<p>That's right, 3 paragraphs.</p>`
	assert.Equal(t, expected, String(f))

	assert.True(t, li1 == li1)
}

func TestClone(t *testing.T) {
	dt := NewDoctype("html")
	div := NewTag("div", "id", "testing")
	div.AddChildren(NewText("this is a test"))
	hr := NewVoidTag("hr")
	f := NewFragment(dt, div, hr)
	cl := String(Clone(f))
	assert.Equal(t, String(f), cl)
	div.AddAttributes("id", "not-equal")
	assert.NotEqual(t, String(f), cl)
}

func TestRemoveChild(t *testing.T) {
	p1 := NewTag("p", "id", "p1")
	p2 := NewTag("p", "id", "p2")
	p3 := NewTag("p", "id", "p3")
	p4 := NewTag("p", "id", "p4")
	p5 := NewTag("p", "id", "p5")
	div := NewTag("div", "id", "testing")
	div.AddChildren(p1, p2, p3, p4, p5)

	assert.Len(t, div.children, 5)
	div.RemoveChild(3)
	assert.Len(t, div.children, 4)
	id, _ := div.Child(2).(TagNode).Attribute("id")
	assert.Equal(t, "p3", id)
	id, _ = div.Child(3).(TagNode).Attribute("id")
	assert.Equal(t, "p5", id)

	div.RemoveChild(3)
	assert.Len(t, div.children, 3)
}

func TestRemoveFragment(t *testing.T) {
	doctype := NewDoctype("html")

	title := NewTag("title")
	title.AddChildren(NewText("Test Doc"))
	head := NewTag("head")

	content := NewFragment()
	scripts := NewFragment()
	body := NewTag("body")
	body.AddChildren(content, scripts)

	root := NewFragment(doctype, head, body)

	p := NewTag("p")
	p.AddChildren(NewText("here's some text"))
	content.AddChildren(p)

	js1 := NewTag("script", "src", "ui.js")
	scripts.AddChildren(
		js1,
		NewTag("script", "src", "logic.js"),
	)

	assert.Equal(t, scripts, js1.Parent())

	before := String(root)
	RemoveFragments(root)
	after := String(root)

	assert.Equal(t, before, after)
	assert.Equal(t, body, js1.Parent())
}

func TestFormatWrap(t *testing.T) {
	li1 := NewTag("li")
	li1.AddChildren(NewText("Item 1"))
	li2 := NewTag("li")
	li2.AddChildren(NewText("Item 2"))
	li3 := NewTag("li")
	li3.AddChildren(NewText("Item 3"))
	ul := NewTag("ul")
	ul.AddChildren(li1, li2, li3)
	p1 := NewTag("p")
	p1.AddChildren(NewText("This is a test"))
	p2 := NewTag("p")
	p2.AddChildren(NewText("This is a very long  paragraph that should be wrapped. It is longer than 80 characters so it will need to be on at least 2 lines, each indented to the same depth. This should be long enough, I hope."))
	p3 := NewTag("p")
	p3.AddChildren(NewText("That's right, 3 paragraphs."))
	div := NewTag("div")
	div.AddChildren(p1, ul, p2)
	f := NewFragment(div, p3)
	expected := `<div>
  <p>This is a test</p>
  <ul>
    <li>Item 1</li>
    <li>Item 2</li>
    <li>Item 3</li>
  </ul>
  <p>
    This is a very long paragraph that should be wrapped. It is longer than 80
    characters so it will need to be on at least 2 lines, each indented to the
    same depth. This should be long enough, I hope.
  </p>
</div>
<p>That's right, 3 paragraphs.</p>`
	assert.Equal(t, expected, String(f))

	assert.True(t, li1 == li1)
}
