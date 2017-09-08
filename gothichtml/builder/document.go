package builder

import (
	"github.com/adamcolton/gothic/gothichtml"
	"io"
)

type Document struct {
	header  *gothichtml.Tag
	body    *gothichtml.Fragment
	scripts *gothichtml.Fragment
	doc     *gothichtml.Fragment
}

func NewDocument(title string, attrs ...string) *Document {
	doc := &Document{
		doc: gothichtml.NewFragment(gothichtml.NewDoctype("html")),
	}
	doc.header = gothichtml.NewTag("header")
	if title != "" {
		titleTag := gothichtml.NewTag("title")
		titleTag.AddChildren(gothichtml.NewText(title))
		doc.header.AddChildren(titleTag)
	}

	body := gothichtml.NewTag("body", attrs...)
	doc.body = gothichtml.NewFragment()
	doc.scripts = gothichtml.NewFragment()
	body.AddChildren(doc.body, doc.scripts)

	doc.doc.AddChildren(doc.header, body)
	return doc
}

func (d *Document) AddChildren(children ...gothichtml.Node) *Document {
	d.body.AddChildren(children...)
	return d
}

func (d *Document) Build() *Builder {
	return &Builder{
		cur: d.body,
	}
}

func (d *Document) Write(w io.Writer) {
	d.doc.Write(w)
}

func (d *Document) String() string {
	return gothichtml.String(d.doc)
}

func (d *Document) CSSLinks(hrefs ...string) *Document {
	for _, href := range hrefs {
		d.header.AddChildren(gothichtml.NewVoidTag("link", "href", href, "rel", "stylesheet", "type", "text/css"))
	}
	return d
}

func (d *Document) ScriptLinks(srcs ...string) *Document {
	for _, src := range srcs {
		d.scripts.AddChildren(gothichtml.NewTag("script", "src", src))
	}
	return d
}
