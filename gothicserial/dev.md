## Dev Notes

Once again, I need to rebuild a bunch of this.

Serialization needs a context. So you might have binary serialization, json
serialization, gob serialization. You might even have different forms of them -
jsonUI and jsonStorage. Which means that the outer most layer is a serialization
context.