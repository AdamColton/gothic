/*
Package gothicserial provides tools for defining serialization in Blueprints.
*/
package gothicserial

/*

1) Defining serialization functions: This is what SerialDef does, it defines
   how a type is serialized and deserialized.
2) Boilerplate serialization: Given a serialDef, it is repetitive to define
   serializing slices or maps.

Questions

If I need []int, where does that go? We need a default output package for
generated serializers.

Given map[foo.Foo]bar.Bar, does that go in foo or bar?

*/
