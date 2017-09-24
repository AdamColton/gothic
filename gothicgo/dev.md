## Dev Notes

### Todo

* set order on methods
* add chan type
* add array type

Change NameType to an interface:
type NameType interace{
  Type
  Name() string
}

#### Func generation
I need to clean up how funcs are generated, I think there's still some issues
with spaces around return values.

#### Interface from Struct
Often I have an interface that only represents a single struct, it would be nice
to automate that, though it should probably happen in buttress.

#### Package
In Package, add exists checks for struct and interface. Add a generators list
and expose an add method. Move the logic for adding references into File.