package entity

import (
	"fmt"
	"github.com/boltdb/bolt"
	"runtime"
	"unsafe"
)

type Entity interface {
	ID() uint64
	Marshal() []byte
}

var bucketID = []byte("test")
var weakmap = make(map[uint64]uintptr)
var db *bolt.DB

func Get(id uint64, ent *Entity, unmarshal func(*[]byte) Entity) {
	if entPtr, ok := weakmap[id]; ok {
		*ent = *((*Entity)(unsafe.Pointer(entPtr)))
	} else {
		db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketID)
			if bucket == nil {
				return fmt.Errorf("Bucket %q not found!", bucketID)
			}

			val := bucket.Get(idToBytes(id))
			*ent = unmarshal(&val)

			return nil
		})
	}
}

func Register(ent Entity) {
	weakmap[ent.ID()] = uintptr(unsafe.Pointer(&ent))
	runtime.SetFinalizer(ent, finalizerRemove)
}

func finalizerRemove(ent Entity) {
	delete(weakmap, ent.ID())
	go Store(ent)
}

func Store(ent Entity) {
	if db == nil {
		db, _ = bolt.Open("test.db", 0644, nil)
	}
	id := ent.ID()
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucketID)
		if err != nil {
			return err
		}

		err = bucket.Put(idToBytes(id), ent.Marshal())
		if err != nil {
			return err
		}
		return nil
	})
}
