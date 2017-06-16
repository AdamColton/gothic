## Dev Notes

Todo
- Fix Snippets
- Models
- Ents
- Struct Method Lib
  - Constructor
  - getter
  - setter
  - appender
  - mux guard?
- Add interfaces (see below)
  - gothicgo.Package
  - gothicgo.File
- HTTP Api
- Logging (see below)

Someday
- Validation
- Serialize Struct Stub
- Serialize Struct Full (iterate over fields)


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

### Logging
Right now, if there's a formatting error, it goes to fmt.Println, which means I
can't check it in tests. Better would be for that to go to some sort of log, or
at least an out stream so tests can check it and users can redirect it.

### Snippets
I've started working on this, but I'll add a few notes just in case.

SnippetContainers hold snippets in buckets. Likely buckets would be
"validators", "mutexLocks", "mutexUnlocks".

Both a Field and a Struct are SnippetContainers. So a setter only needs to
access snippets on the field. Things like validators that span multiple fields
can be called in constructors.

See notes on Generators

### Add interfaces
I really like the pattern of a private struct with a public interface.
Particularly for what I'm doing here, it makes it easy for someone to build a
replacement. I need to convert a few of the structs over to this style.

## Moon-Jumper Plan
I've got a good chunk of this in place. No where near done, but enough to try it
out. I'm going to start building the Go-Moon-Jumper with it and add things as I
need them.