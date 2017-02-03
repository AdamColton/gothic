// Package entbp provides the Blueprint functions for defining structs that
// implement ent. Marshal and Unmarshal functions are automatically generated.
package entbp

/*

person := gothicEnt.New("Person").
  AddField("Name", gothicgo.String()).
  AddField("Age", gothicgo.Int()).
  AddField("Role", gothicgo.String())

personStruct := person.Go(package) // automatically exports the struct and Marshal/Unmarshal
gothicgoMethods.Constructor(personStruct) // creates NewPerson and takes all fields.
gothicgoMethods.Constructor(personStruct, "NewPersonByName", "Name").
  Set("Role", `"user"`)

person.JS() // automatically exports a POJO and the Marshal/Unmarshal methods
person.LapizJS() // exports Lapiz object M/UM

that's it, that will create
*/
