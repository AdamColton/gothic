package gothichtml

// Clone takes a Node and clones it and all it's descendants. It performs a deep
// copy.
func Clone(node Node) Node {
	switch n := node.(type) {
	case *Doctype:
		d := &Doctype{
			doctype: n.doctype,
		}
		d.parent = newParent(d)
		return d
	case *Fragment:
		f := &Fragment{}
		f.fragment = cloneFragment(n.fragment, f)
		return f
	case *Tag:
		t := &Tag{
			tag:        n.tag,
			attributes: cloneAttributes(n.attributes),
		}
		t.fragment = cloneFragment(n.fragment, t)
		return t
	case *VoidTag:
		t := &VoidTag{
			tag:        n.tag,
			attributes: cloneAttributes(n.attributes),
		}
		t.parent = newParent(t)
		return t
	case *Text:
		t := &Text{
			Text: n.Text,
			Wrap: n.Wrap,
		}
		t.parent = newParent(t)
		return t
	default:
		panic("Unknown type, this should be unreachable")
	}
}

func cloneAttributes(in attributes) attributes {
	out := make(attributes, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func cloneFragment(n *fragment, p Node) *fragment {
	f := &fragment{
		children: make([]Node, n.Children()),
		parent: &parent{
			parent: p,
		},
	}
	f.self = f
	for i, c := range n.children {
		cl := Clone(c)
		cl.setParent(f)
		f.children[i] = cl
	}
	return f
}
