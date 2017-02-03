/*
Gothic is a set of code generation tools.

Terminology

A Blueprint is a Go program that is run like a script (it's generally not built,
just invoked with go run ...) and exports code.

A Project is the collection of files produced by a Blueprint. Generally these
will include Go files where the packages are a mix of generated and written
code.

Directory Layout

Given gothic/foo, foo will be a package intended for use in a project.

Given gothic/foo/fooBP, fooBP will be a representation of foo for use in a
Blueprint.

Given gothic/gothicFoo, gothicFoo will define a set of interfaces to be used in
Blueprints. The gothicFoo package may also provide structs that implement the
interface.
*/
package gothic
