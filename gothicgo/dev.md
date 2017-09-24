## Dev Notes

### Todo

* set order on methods

there should be a contextual wrapper around Resolver to prevent adding a package
to itself.

Change NameType to an interface:
type NameType interace{
  Type
  Name() string
}

#### Comments
Add a comments utility to take a comment string, wrap it at a given col width (
80 by default) and prepend // to each line. Functions, Structs and Types should
all have a way to attach a comment. A file should also be able to have a package
comment. And it would be nice to be able to generate a doc.go file from a
comment on a package.


#### Func generation
I need to clean up how funcs are generated, I think there's still some issues
with spaces around return values.

#### Interface from Struct
Often I have an interface that only represents a single struct, it would be nice
to automate that. 

test