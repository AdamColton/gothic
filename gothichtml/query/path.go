package query

import (
	"github.com/adamcolton/gothic/gothichtml"
)

type path []int

func Path(ps ...int) path {
	return path(ps)
}

func (p path) Query(n gothichtml.Node) gothichtml.Node {
	cur := n

	for _, idx := range p {
		parent, ok := cur.(gothichtml.ContainerNode)
		if !ok {
			return nil
		}
		cur = parent.Child(idx)
	}

	return cur
}

func (p path) QueryTag(n gothichtml.Node) *gothichtml.Tag {
	n = p.Query(n)
	if tag, ok := n.(*gothichtml.Tag); ok {
		return tag
	}
	return nil
}

func (p path) QueryVoidTag(n gothichtml.Node) *gothichtml.VoidTag {
	n = p.Query(n)
	if tag, ok := n.(*gothichtml.VoidTag); ok {
		return tag
	}
	return nil
}

func (p path) Clone() path {
	cln := make(path, len(p))
	copy(cln, p)
	return cln
}

type Paths []path

func (ps Paths) Query(n gothichtml.Node) []gothichtml.Node {
	nodes := make([]gothichtml.Node, 0, len(ps))
	for _, p := range ps {
		if n := p.Query(n); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}
