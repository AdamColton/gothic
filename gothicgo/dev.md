## Dev Notes

### Todo

* set order on methods

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

#### Package rebuild
Right now, we do relative package comparison by package name. But that won't
handle collision. It also won't handle files where a package has been aliased.

Packages should be referenced by their import. And I'll probably need to do a
subtle but significant rebuild of Types so that a Type can indicate what package
it belongs to then the file can insert the package name.