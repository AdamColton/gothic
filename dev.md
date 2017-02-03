## Dev Notes

Todo
- Struct Method Lib
  - Constructor
  - getter
  - setter
  - appender
  - mux guard?
- Add interfaces (see below)
  - gothicgo.Package
  - gothicgo.File
- Add indexing to Ent
- Add migration to ent
- Export ent to Lapiz
- HTTP Api
- Logging (see below)

Someday
- Validation
- Serialize Struct Stub
- Serialize Struct Full (iterate over fields)

### Replace existing func generators
A little different than I originally thought. The Ent wrapper will ONLY marshal
and unmarshal the fields in the Ent. This lets the user add fields to the
struct, which might be useful...

Need to be careful on the ordering, I think I'll want to generate all the
methods during GS.Prepare(), and Generate will not really do anything.

Probably want to break Ref out into it's own full blown struct in the file.

### Validation
This is tricky.

I know I need a validation collector interface.

I think the way this will work is that a validator has a sort of interface, you
give it the names of the variables it's validating and the collector, and it
produces a chunk of code.

So a min length validator takes the var name it's validating, a length or var
containing the length and the var name of the validation collector and produces
a chunk of code on it:

`
if len(%s) < %s {
  %s.Add("Min Length: %s must be at least %s")
}
`

fmt.Sprintf(minLenF, varName, minLength, collector, varname, minLength)
given:
  varName = "name"
  minLength = "3"
  collector = "valErrs"
we get
if len(name) < minLength {
  valErrs.Add("Min Length: name must be at least 3")
}

### Lazy package resolver
I originally came across this for File, but it could be useful elsewhere (struct
func). A lazy package resolver can handle paths, packages and aliases. It holds
packages names until it's told to resolve them.

### Logging
Right now, if there's a formatting error, it goes to fmt.Println, which means I
can't check it in tests. Better would be for that to go to some sort of log, or
at least an out stream so tests can check it and users can redirect it.

### Ent Version Numbers
Part of the Ent spec is going to be a Version number. This is the first value in
a serialized struct. By default, it is one byte. A version of 0 represents a nil
instance. If the version number reaches 255, then 2 more bytes are added
supporting version numbers up to 65788. More revisions than that, and you've got
bigger problems.

### Ent ID
Right now, an Ent ID is a byte slice that is supposed to be 8 bytes long. That
feels a little sloppy. I should either make it an array, or have it cope with
variation.

But part of the reason I did that was to make it easy to use as a bolt key. So
the length really needs to be fixed.

### Marshaling
To do marshaling efficiently, it may get complex down the road. Right now, I
have a single Marshal method per Ent. But for storage it makes more sense to not
include the ID in the data because I'll always have access to it when reading
from the DB. And, when I add indexing, I won't want to store that either, again,
so a user can do a quick read on just IDs, or just indexed fields.

### Snippets
I've started working on this, but I'll add a few notes just in case.

SnippetContainers hold snippets in buckets. Likely buckets would be
"validators", "mutexLocks", "mutexUnlocks".

Both a Field and a Struct are SnippetContainers. So a setter only needs to
access snippets on the field. Things like validators that span multiple fields
can be called in constructors.

## Cleaner Ent Models
Look at demo.

This turns into a nasty problem quickly. We really want to be able to have Ents
reference eachother rather than translating to Go structs and reference those.
But what about non-Ent structs. So Ent really needs a way to define an object
that isn't an ent. Consider Date - we'd use the builtin for Go and JavaScript.
We need a way to define Date at each level (Ent, Go, JS).

I want to keep using the Go type system. Ent needs some type system and the Go
type system is as good as any, no reason to build a new one. The one place where
this sort of breaks down is Packages, but that's not too bad, we'll just
incorporate the idea of packages in to entities. While it's not a formal part of
the language, it's present in most JavaScript anyway.

## Add interfaces
I really like the pattern of a private struct with a public interface.
Particularly for what I'm doing here, it makes it easy for someone to build a
replacement. I need to convert a few of the structs over to this style.