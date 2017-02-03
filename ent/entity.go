package ent

// Entity is a type that returns an ID. The ID should be globally unique. An
// entity can be marshalled into a byte slice. It should also be paired with
// an Unmarshaler.
type Entity interface {
	ID() []byte
	Marshal() []byte
}

// Unmarshaler is a functino that takes a pointer to a byte slice and returns
// an entity. Using a pointer makes it easier to recursivly unmarshal.
type Unmarshaler func(*[]byte) Entity
