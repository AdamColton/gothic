/*
Ent, short for entity, represents a struct that has a uint64 ID and can be
Marshaled and Unmarshaled to []byte.

For persistance, ent provides tools to save to a bolt database.
*/
package ent

/** Indexing **
All instances of an index must be the same length, so for variable length
fields, we need to define a length. To use a default length, -1.

entBP.Index(personStruct).
  Field("Name",5).
  Field("Age",-1)

type PersonIdxNameAge struct{
  id []byte
  personId []byte
}
func (p *Person) idxNameAge(){

}
*/
