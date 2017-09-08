package query

import (
	"github.com/adamcolton/gothic/gothichtml"
	"github.com/adamcolton/parlex"
	"github.com/adamcolton/parlex/grammar/regexgram"
	"github.com/adamcolton/parlex/lexer/simplelexer"
	"github.com/adamcolton/parlex/parser/packrat"
	"github.com/adamcolton/parlex/tree"
)

const (
	lexerRules = `
class   /\.[\w\-]+/
id      /#[\w\-]+/
tag     /[\w\-]+/
comma   /[ \t]*,[ \t]*/
child   /[ \t]*>[ \t]*/
nextSib /\+/
after   /\~/
descd   /[ \t]+/
`

	grammarRules = `
Selectors -> (Selector comma)* Selector
Selector  -> Element (child|after|descd|nextSib Element)*
Element   -> tag class* id? class*
          -> class class* id? class*
          -> class* id class*
`
)

var lxr = parlex.MustLexer(simplelexer.New(lexerRules))
var grmr, grmrRdcr = regexgram.Must(grammarRules)
var prsr = packrat.New(grmr)

var rdcr = tree.Merge(grmrRdcr, tree.Reducer{
	"Selectors": tree.RemoveAll("comma"),
})

var runner = parlex.New(lxr, prsr, rdcr)

type nodeChecker interface {
	check(gothichtml.TagNode) bool
}

type checkTag string

func (m checkTag) check(node gothichtml.TagNode) bool {
	return node.Name() == string(m)
}

type checkClass string

func (m checkClass) check(node gothichtml.TagNode) bool {
	for _, class := range gothichtml.Classes(node) {
		if class == string(m) {
			return true
		}
	}
	return false
}

type checkID string

func (m checkID) check(node gothichtml.TagNode) bool {
	id, _ := node.Attribute("id")
	return id == string(m)
}

type locationChecker func([]int) bool

type descendant []int

func newDescendant(where []int) locationChecker { return descendant(where).check }

// checks that loc is a descendant of where
func (where descendant) check(loc []int) bool {
	if len(where) > len(loc) {
		return false
	}
	for i, w := range where {
		if w != loc[i] {
			return false
		}
	}
	return true
}

type child []int

func newChild(where []int) locationChecker { return child(where).check }

func (where child) check(loc []int) bool {
	if len(where)+1 != len(loc) {
		return false
	}
	for i, w := range where {
		if w != loc[i] {
			return false
		}
	}
	return true
}

type after []int

func newAfter(where []int) locationChecker { return after(where).check }

func (where after) check(loc []int) bool {
	if len(where) != len(loc) {
		return false
	}
	for i, w := range where[:len(where)-1] {
		if w != loc[i] {
			return false
		}
	}
	return where[len(where)-1] < loc[len(loc)-1]
}

type nextSib []int

func newNextSib(where []int) locationChecker { return nextSib(where).check }

func (where nextSib) check(loc []int) bool {
	if len(where) != len(loc) {
		return false
	}
	for i, w := range where[:len(where)-1] {
		if w != loc[i] {
			return false
		}
	}
	return where[len(where)-1]+1 == loc[len(loc)-1]
}

type selector struct {
	loc      locationChecker
	checkers []nodeChecker
	next     *selector
	nextLoc  func([]int) locationChecker
}

func (s *selector) check(node gothichtml.TagNode, loc *Location) bool {
	if s.loc != nil && !s.loc(loc.Tag) {
		return false
	}
	for _, c := range s.checkers {
		if !c.check(node) {
			return false
		}
	}
	return true
}

type selectors []*selector

func (s selectors) QueryAll(node gothichtml.Node) []gothichtml.TagNode {
	op := &selectorsOp{
		selectors: s,
	}
	Walk(node, op.compare)
	return op.matches
}

type selectorsOp struct {
	selectors
	matches []gothichtml.TagNode
}

// this is just to organize my thoughts
func (op *selectorsOp) compare(node gothichtml.Node, loc *Location) {
	tag, ok := node.(gothichtml.TagNode)
	if !ok {
		return
	}

	var toAppend []*selector

	matched := false
	for _, sel := range op.selectors {
		if matched && sel.next == nil {
			// no need to compare final acceptors
			continue
		}
		if sel.check(tag, loc) {
			if sel.next == nil {
				op.matches = append(op.matches, tag)
				matched = true
			} else {
				next := &selector{
					loc:      sel.nextLoc(loc.Tag),
					checkers: sel.next.checkers,
					next:     sel.next.next,
					nextLoc:  sel.nextLoc,
				}
				toAppend = append(toAppend, next)
			}
		}
	}

	op.selectors = append(op.selectors, toAppend...)
}

type SelectorQuery interface {
	QueryAll(node gothichtml.Node) []gothichtml.TagNode
}

func Selector(str string) (SelectorQuery, error) {
	parsedSelectors, err := runner.Run(str)
	if err != nil {
		return nil, err
	}

	sels := make(selectors, parsedSelectors.Children())
	for i := range sels {
		sels[i] = parseSelector(parsedSelectors.Child(i))
	}
	return sels, nil
}

func parseSelector(parsedSelector parlex.ParseNode) *selector {
	var root, cur *selector
	for i := 0; i < parsedSelector.Children(); i++ {
		child := parsedSelector.Child(i)
		switch child.Kind().String() {
		case "Element":
			if root == nil {
				root = parseElement(child)
				cur = root
			} else {
				cur.next = parseElement(child)
				cur = cur.next
			}
		case "descd":
			cur.nextLoc = newDescendant
		case "child":
			cur.nextLoc = newChild
		case "after":
			cur.nextLoc = newAfter
		case "nextSib":
			cur.nextLoc = newNextSib
		}
	}
	return root
}

func parseElement(element parlex.ParseNode) *selector {
	s := &selector{}
	for i := 0; i < element.Children(); i++ {
		checkType := element.Child(i)
		switch checkType.Kind().String() {
		case "tag":
			s.checkers = append(s.checkers, checkTag(checkType.Value()))
		case "class":
			v := []byte(checkType.Value())[1:] // remove leading .
			s.checkers = append(s.checkers, checkClass(v))
		case "id":
			v := []byte(checkType.Value())[1:] // remove leading #
			s.checkers = append(s.checkers, checkID(v))
		}
	}
	return s
}
