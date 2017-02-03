// Used to generate Go code in Gothic.
package gothicgo

/*

I'm hitting a conceptual issue with types and packages. A struct has a package.
The built in types don't have a package, which is treated as "". A slice of a
type can be treated as belonging to the underlying type's package.

But given map[foo.Foo]bar.Bar, what is the package?

A type doesn't have a package, either it doesn't specify or it can have many
packages. The package name is important for determining imports, which leans
towards the many.

*/
