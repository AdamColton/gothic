// Package parsehtml will parse an html string into a gothichtml representation.
package parsehtml

import (
	"fmt"
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar/regexgram"
	"github.com/adamcolton/parlex/lexer/stacklexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
)

const (
	lexerRules = `
== Main ==
doctype   /<!DOCTYPE/ Tag
closeTag /<\/\s*([^>]*)\s*>/ (1)
openTag  /<([!\w]+)/ (1) Tag
text     /[^<\s][^<]*[^<\s]+/
ws       /\s+/ -
== Tag ==
voidClose /\/>/ ^
close     />/ ^
word      /[\w\$]([\w\/\.]*[\w])?/
eq        /=/
string    /"((?:[^"\\]|\\.)*)"/ (1)
ws        /\s+/ -
`

	grammarRules = `
Doc         -> DocType? Node*
DocType     -> doctype word close
Node        -> Tag|VoidTag|text
Tag         -> openTag Attributes close Contents closeTag
Contents    -> Node*
VoidTag     -> openTag Attributes voidClose
Attributes  -> Attribute*
Attribute   -> word eq word|string
`
)

var lxr = parlex.MustLexer(stacklexer.New(lexerRules))
var grmr, grmrRdcr = regexgram.Must(grammarRules)
var prsr = packrat.New(grmr)

var rdcr = tree.Merge(grmrRdcr, tree.Reducer{
	"Node":      tree.PromoteSingleChild,
	"Tag":       tree.RemoveChildren(-3).PromoteChildValue(0),
	"Attribute": tree.RemoveChildren(1),
	"VoidTag":   tree.RemoveChildren(-1).PromoteChildValue(0),
	"DocType":   tree.PromoteChildValue(1).RemoveChildren(0, 0),
})

var runner = parlex.New(lxr, prsr, rdcr)

// Parse takes an HTML string and returns the root Node representing that html
// as a gothichtml tree.
func Parse(str string) (gothichtml.Node, error) {
	pn, err := runner.Run(str)
	if err != nil {
		return nil, err
	}

	op := &evalOp{
		doc: gothichtml.NewFragment(),
	}
	op.cur = op.doc
	op.eval(pn.(*tree.PN))

	if op.err != nil {
		return nil, op.err
	}

	return op.doc, nil
}

// Must is the same as Parse but will panic instead of returning an error
func Must(str string) gothichtml.Node {
	n, err := Parse(str)
	if err != nil {
		panic(err)
	}
	return n
}

type evalOp struct {
	doc *gothichtml.Fragment
	cur gothichtml.ContainerNode
	err error
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
	case "DocType":
		op.cur.AddChildren(gothichtml.NewDoctype(pn.Value()))
	case "VoidTag":
		op.voidtag(pn)
	case "Tag":
		op.tag(pn)
	case "text":
		op.cur.AddChildren(gothichtml.NewText(pn.Value()))
	}
}

var primaries = map[string]string{
	"label": "for",
	"a":     "href",
}

func (op *evalOp) voidtag(pn *tree.PN) {
	tag := gothichtml.NewVoidTag(pn.Value())
	addAttributes(tag, pn.C[0])
	op.cur.AddChildren(tag)
}

func (op *evalOp) tag(pn *tree.PN) {
	if clt := pn.Child(-1).(*tree.PN); pn.Value() != clt.Value() {
		oc, ol := pn.Lexeme.Pos()
		cc, cl := clt.Lexeme.Pos()
		op.err = fmt.Errorf("Opening tag '%s' at (%d,%d) does not match closing tag '%s' at (%d,%d)", pn.Value(), oc, ol, clt.Value(), cc, cl)
	}
	tag := gothichtml.NewTag(pn.Value())
	addAttributes(tag, pn.C[0])

	op.cur.AddChildren(tag)
	prev := op.cur
	op.cur = tag
	for _, c := range pn.C[1].C {
		op.eval(c)
		if op.err != nil {
			return
		}
	}
	op.cur = prev
}

func addAttributes(tag gothichtml.TagNode, attrs *tree.PN) {
	for _, c := range attrs.C {
		tag.AddAttributes(c.C[0].Value(), c.C[1].Value())
	}
}
