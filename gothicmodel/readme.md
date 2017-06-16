## Model

[![GoDoc](https://godoc.org/github.com/AdamColton/gothic/ent/entbp?status.svg)](https://godoc.org/github.com/AdamColton/gothic/ent/entbp)

A model is a collection of fields and types that will be shared across several
levels of the stack. For instance, a model may need to be stored in a SQL
database, used in Go on a server and represented in JavaScript on the client. A
Model type is just a string, any package that interfaces with Entity needs to
provide a mapping of string to type.