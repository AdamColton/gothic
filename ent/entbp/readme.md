## Ent BP

[![GoDoc](https://godoc.org/github.com/AdamColton/gothic/ent/entbp?status.svg)](https://godoc.org/github.com/AdamColton/gothic/ent/entbp)

Ent BP is used to describe Entities in a Blueprint project. An Entity is
anything that has an ID and can be Marshaled and Unmarshaled. It is not a Go
specific concept (though Go types are used).

The workflow for Blueprint project will be to define an Entity, then export it
to the targets (such as Go and JavaScript) and possibly add additional target
specific features.