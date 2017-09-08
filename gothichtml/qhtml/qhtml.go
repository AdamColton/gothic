package qhtml

import (
	"fmt"
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar/regexgram"
	"github.com/adamcolton/parlex/lexer/stacklexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
	"strings"
)

const (
	lexerRules = `
== Main ==
closeTag /<\/[^>]*>/
openTag  /<(\w+)/ (1) Tag
text     /[^<\s][^<]*[^<\s]+/
ws       /\s+/ -
== Tag ==
voidClose /\/>/ ^
close     />/ ^
class     /\.(\w+)/ (1)
id        /#(\w+)/ (1)
set        /=/ TagData
ws        /[ \t]+/ TagData -
== TagData ==
voidClose /\/>/ ^^
close     />/ ^^
word      /[\w\$]([\w\/\.]*[\w])?/
eq        /=/
string    /"((?:[^"\\]|\\.)*)"/ (1)
`

	grammarRules = `
  Doc         -> Node*
  Node        -> Tag|VoidTag|text
  Tag         -> openTag Attributes close Contents closeTag
  Contents    -> Node*
  VoidTag     -> openTag Attributes voidClose
  Attributes  -> id? class* Primary? Attribute*
  Primary     -> set word|string
  Attribute   -> word eq word|string
`
)

var lxr = parlex.MustLexer(stacklexer.New(lexerRules))
var grmr, grmrRdcr = regexgram.Must(grammarRules)
var prsr = packrat.New(grmr)

var rdcr = tree.Merge(grmrRdcr, tree.Reducer{
	"Node":      tree.PromoteSingleChild,
	"Tag":       tree.PromoteChildValue(0).RemoveChildren(1, -1),
	"VoidTag":   tree.RemoveChildren(-1).PromoteChildValue(0),
	"Primary":   tree.PromoteChildValue(1).RemoveChildren(0),
	"Attribute": tree.RemoveChildren(1),
})

var runner = parlex.New(lxr, prsr, rdcr)

func ParseDocument(str string) (gothichtml.Node, error) {
	pn, err := runner.Run(str)
	if err != nil {
		return nil, err
	}

	op := &evalOp{
		header:  gothichtml.NewTag("header"),
		body:    gothichtml.NewTag("body"),
		scripts: gothichtml.NewFragment(),
	}
	op.cur = op.body
	op.eval(pn.(*tree.PN))
	op.body.AddChildren(op.scripts)

	doc := gothichtml.NewFragment(
		gothichtml.NewDoctype("html"),
		op.header,
		op.body,
	)
	return doc, nil
}

func ParseFragment(str string) (gothichtml.Node, error) {
	pn, err := runner.Run(str)
	if err != nil {
		return nil, err
	}

	op := &evalOp{
		body: gothichtml.NewFragment(),
	}
	op.cur = op.body
	op.eval(pn.(*tree.PN))

	if op.body.Children() == 1 {
		return op.body.Child(0), nil
	}

	return op.body, nil
}

func MustFragment(str string) gothichtml.Node {
	f, err := ParseFragment(str)
	if err != nil {
		panic(err)
	}
	return f
}

func Tag(str string) *gothichtml.Tag {
	return MustFragment(str).(*gothichtml.Tag)
}

func VoidTag(str string) *gothichtml.VoidTag {
	return MustFragment(str).(*gothichtml.VoidTag)
}

type evalOp struct {
	header  *gothichtml.Tag
	body    gothichtml.ContainerNode
	scripts *gothichtml.Fragment
	cur     gothichtml.ContainerNode
	err     error
}

func (op *evalOp) eval(pn *tree.PN) {
	switch pn.Kind().String() {
	case "Doc", "Contents":
		for _, c := range pn.C {
			op.eval(c)
			if op.err != nil {
				break
			}
		}
	case "VoidTag":
		op.voidtag(pn)
	case "Tag":
		op.tag(pn)
	case "text":
		op.cur.AddChildren(gothichtml.NewText(pn.Value()))
	}
}

var voidHandlers = map[string]func(*evalOp, *tree.PN){
	"title":   title,
	"css":     css,
	"scripts": scripts,
}

var primaries = map[string]string{
	"label": "for",
	"a":     "href",
}

func (op *evalOp) voidtag(pn *tree.PN) {
	if handler, ok := voidHandlers[pn.Value()]; ok {
		handler(op, pn)
		return
	}
	tag := gothichtml.NewVoidTag(pn.Value())
	addAttributes(pn, tag)
	op.cur.AddChildren(tag)
}

func (op *evalOp) tag(pn *tree.PN) {
	tag := gothichtml.NewTag(pn.Value())
	addAttributes(pn, tag)
	op.cur.AddChildren(tag)
	prev := op.cur
	op.cur = tag
	for _, c := range pn.C[1].C {
		op.eval(c)
	}
	op.cur = prev
}

func addAttributes(pn *tree.PN, tag gothichtml.TagNode) {
	for _, c := range pn.C[0].C {
		addAttribute(c, tag)
	}
}

func addAttribute(attr *tree.PN, tag gothichtml.TagNode) {
	switch attr.Kind().String() {
	case "Primary":
		key, ok := primaries[tag.Name()]
		if !ok {
			key = "value"
		}
		tag.AddAttributes(key, attr.Value())
	case "id":
		tag.AddAttributes("id", attr.Value())
	case "Attribute":
		key := attr.C[0].Value()
		if cur, ok := tag.Attribute(key); ok {
			tag.AddAttributes(key, fmt.Sprintf("%s %s", cur, attr.C[1].Value()))
		} else {
			tag.AddAttributes(key, attr.C[1].Value())
		}
	case "class":
		if cur, ok := tag.Attribute("class"); ok {
			tag.AddAttributes("class", fmt.Sprintf("%s %s", cur, attr.Value()))
		} else {
			tag.AddAttributes("class", attr.Value())
		}
	}
}

func title(op *evalOp, pn *tree.PN) {
	if op.header == nil {
		op.err = fmt.Errorf("Cannot use <title />in fragment")
		return
	}
	tag := gothichtml.NewTag("title")
	addAttributes(pn, tag)
	titleText, _ := tag.Attribute("value")
	tag.Remove("value")
	tag.AddChildren(gothichtml.NewText(titleText))
	op.header.AddChildren(tag)
}

func css(op *evalOp, pn *tree.PN) {
	if op.header == nil {
		op.err = fmt.Errorf("Cannot use <css /> in fragment")
		return
	}
	dummy := gothichtml.NewVoidTag("css")
	addAttributes(pn, dummy)
	hrefs, _ := dummy.Attribute("value")
	for _, href := range strings.Split(hrefs, " ") {
		if len(href) == 0 {
			continue
		}
		tag := gothichtml.NewVoidTag("link", "rel", "stylesheet", "href", href)
		op.header.AddChildren(tag)
	}
}

func scripts(op *evalOp, pn *tree.PN) {
	if op.header == nil {
		op.err = fmt.Errorf("Cannot use <scripts /> in fragment")
		return
	}
	dummy := gothichtml.NewVoidTag("scripts")
	addAttributes(pn, dummy)
	srcs, _ := dummy.Attribute("value")
	for _, src := range strings.Split(srcs, " ") {
		if len(src) == 0 {
			continue
		}
		tag := gothichtml.NewTag("script", "src", src)
		op.scripts.AddChildren(tag)
	}
}
