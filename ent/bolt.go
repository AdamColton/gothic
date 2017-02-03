package ent

import (
	"fmt"
	"github.com/boltdb/bolt"
	"runtime"
	"unsafe"
)

var BucketID = []byte("entities")
var FileName = "entities.db"
var weakmap = make(map[uint64]uintptr)
var db *bolt.DB

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

func Register(ent Entity) {
	crcID := CRC64(ent.ID())
	if _, ok := weakmap[crcID]; !ok {
		weakmap[crcID] = uintptr(unsafe.Pointer(&ent))
		runtime.SetFinalizer(ent, finalizerRemove)
	}
}

func finalizerRemove(ent Entity) {
	delete(weakmap, CRC64(ent.ID()))
	go Store(ent)
}

func SaveAll() {
	var ent Entity
	for _, ptr := range weakmap {
		ent = *((*Entity)(unsafe.Pointer(ptr)))
		Store(ent)
	}
}

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
}
