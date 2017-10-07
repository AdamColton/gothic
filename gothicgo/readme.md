## Gothic Go
[![GoDoc](https://godoc.org/github.com/AdamColton/gothic/gothicgo?status.svg)](https://godoc.org/github.com/AdamColton/gothic/gothicgo)

Gothic Go is a code generation tool. It is also no where near done, but it is
approaching functional.

Being part of the Gothic project, it's based around 2-pass code generation.
After the project is initilized there is a Prepare pass then a Generate pass.
For Go, the Prepare pass allows a lot of the references to be set so they will
resolve correctly in the Generate pass.

To put Gothic Go in context; Reflect performs operations on code that is
running. AST performs operations on code that is already written. Gothic Go
describes code to be generated.

### Development and things that are lacking

I've taken the approach of writing things as I need them. As a result, several
important things are missing. A notable exampe is Channels. I plan to support
them, but I haven't needed them so I haven't written them.

Early in the project I tried to build things based on how I thought they would
work. Those sections tended to need major reworking when I actually tried to use
them. But the cases where I've started with a specific need and worked backwards
have produced the best code.

So if you'd like to help, try to make something and show me the walls you run
into, or (even better) fix the walls and open a pull request.

### Package References

The most challenging aspect of generating Go code is keeping package references
correct. This is particularly difficult to get right without over-burdening the
user.

Gothic Go relies on two concepts to handle this. PackageRefs and Prefixers. A
PackageRef is really just a wrapper around a string. That string is what you
would put in the imports section to import a package, just the part that appears
in quotes.

A Prefixer is actually going to be a Go file (most of the time). It represents
the imports declarations, including the aliases. A Prefixer takes a PackageRef
and returns the prefix that represents that PackageRef in the context of the
prefixer.

Confusing, but an example will make it clear. Lets say we have 3 files

```Go
package foo

import(
  "github.com/adamcolton/bar"
)
```

```Go
package foo

import(
  acBar "github.com/adamcolton/bar"
  "go/bar"
)
```

```Go
// this file is in github.com/adamcolton/bar
package bar

import(
  "go/bar"
)
```

If I want to reference a function "Hello" in the package
github.com/adamcolton/bar the prefix will be "bar." in the first file, "acBar."
in the second and "" in the third. The PackageRef is "github.com/adamcolton/bar"
and it is the same in all cases, but the Prefixer depends on the imports and
aliases.

If you create a Package, it fulfills the PackageRef interface and a File
fulfills Prefixer. Most of this is useful to know but won't effect your code
unless you're writing a code Generator as opposed to using a code Generator.

### Types

The next tricky bit is types. A lot of that has to do with the package problem.
Each Type knows it's PackageRef. And the Type interface includes a RelStr method
that returns the type with the correct prefix.

All instances of a type (such as a Struct or a Func) either fulfill the Type
interface or have a method to return their Type.

### Everything Else

Despite being a fairly large package, outside of the complications about
resolving packages, types and imports, most of the tools are pretty straight
forward.

Finally, if you're looking to generate code, you probably don't want to use
gothic directly. Gothic Model is used to generate models and handle lots of
boilerplate. [Buttress](https://github.com/adamcolton/buttress) provides a
number of supporting tools for gothic and it may provide the types of tools you
want.

If you're looking to build a code generator that will abstract away some
specific piece of boilerplate, then gothic is probably the right tool.
