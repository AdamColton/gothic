package builder

import (
	"github.com/adamcolton/gothic/gothichtml"
)

type Builder struct {
	cur   gothichtml.ContainerNode
	stack []gothichtml.ContainerNode
}

func New() *Builder {
	return &Builder{
		cur: gothichtml.NewFragment(),
	}
}

func (b *Builder) Text(text string) *Builder {
	b.cur.AddChildren(gothichtml.NewText(text))
	return b
}

func (b *Builder) Tag(tag string, attrs ...string) *Builder {
	t := gothichtml.NewTag(tag, attrs...)
	b.cur.AddChildren(t)
	b.push(t)
	return b
}

func (b *Builder) Close() *Builder {
	if l := len(b.stack); l > 0 {
		b.cur = b.stack[l-1]
		b.stack = b.stack[:l-1]
	}
	return b
}

func (b *Builder) push(node gothichtml.ContainerNode) {
	b.stack = append(b.stack, b.cur)
	b.cur = node
}

func (b *Builder) Cur() gothichtml.Node {
	return b.cur
}
