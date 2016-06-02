## Gothic Serial

Performs serialization on basic types.

### Core
The file core.go contains the core serialization algorithms.

Most serialization/deserialzation will at some point envoke MarshalUintL/UnmarshalUintL. These are used to respectively marshal and unmarshal unsigned ints with a specified maximum byte length.

MarshalIntL/UnmarshalIntL will marshal signed ints. But, most negative values tend to be closer to 0 than the max int, which means with twos compliment, we end up saving a lot of bytes we don't necessarily need. So these convert the sign to a bit and place it in the least significant bit. The assumption is that the cost of convertion is significantly lower than the cost of either disk or network.

MarshalByteSlice prepends a byte slice with a length as a uint64 encoded with MarshalUintL.

### Wrappers
Provides wrappers for all the built-in types.

### Blueprints
The serialBP package can be used to produce serialization for almost any type.

### To-do
* complex

Maybe add a reflect class, though it's a little antithetical to Gothic.