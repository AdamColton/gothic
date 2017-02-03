package ent

import (
	"fmt"
	"github.com/boltdb/bolt"
	"runtime"
	"unsafe"
)

// BucketID is the ID used for the entities bucket
var BucketID = []byte("entities")

// FileName is the filename used for the entities database
var FileName = "entities.db"

var weakmap = make(map[uint64]uintptr)
var db *bolt.DB

// Get returns an Entity by ID. It will first check if there is a reference to
// the entity is use and if not, it will hydrate the entity from the database
func Get(id []byte, ent *Entity, unmarshal Unmarshaler) {
	if entPtr, ok := weakmap[CRC64(id)]; ok {
		*ent = *((*Entity)(unsafe.Pointer(entPtr)))
	} else {
		if db == nil {
			db, _ = bolt.Open(FileName, 0644, nil)
		}
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(BucketID)
			if bucket == nil {
				return fmt.Errorf("Bucket %q not found!", BucketID)
			}

			val := bucket.Get(id)
			if len(val) == 0 {
				return fmt.Errorf("Record not found")
			}
			*ent = unmarshal(&val)

			return nil
		})
		if err == nil {
			Register(*ent)
		}
	}
}

// Register must be called when an entity is first created. Once Register is
// called, the entity will automatically be saved when it is garbage collected.
func Register(ent Entity) {
	crcID := CRC64(ent.ID())
	if _, ok := weakmap[crcID]; !ok {
		weakmap[crcID] = uintptr(unsafe.Pointer(&ent))
		runtime.SetFinalizer(ent, finalizerRemove)
	}
}

func finalizerRemove(ent Entity) {
	id := CRC64(ent.ID())
	if _, ok := weakmap[id]; ok {
		delete(weakmap, id)
		go Store(ent)
	}
}

// SaveAll will save all the entities that are currently active. This can be
// expensive.
func SaveAll() {
	var ent Entity
	for _, ptr := range weakmap {
		ent = *((*Entity)(unsafe.Pointer(ptr)))
		Store(ent)
	}
}

// Store saves an entity to the database. The reference to the entity is still
// valid.
func Store(ent Entity) {
	Register(ent)
	if db == nil {
		db, _ = bolt.Open(FileName, 0644, nil)
	}
	id := ent.ID()
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BucketID)
		if err != nil {
			return err
		}

		err = bucket.Put(id, ent.Marshal())
		if err != nil {
			return err
		}
		return nil
	})
}

// Delete removes the entity from the database and prevents it from
func Delete(ent Entity) {
	if db == nil {
		db, _ = bolt.Open(FileName, 0644, nil)
	}
	id := ent.ID()
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(BucketID)
		if err != nil {
			return err
		}

		err = bucket.Delete(id)
		if err != nil {
			return err
		}
		return nil
	})
	id64 := CRC64(id)
	if _, ok := weakmap[id64]; ok {
		delete(weakmap, id64)
	}
}
