package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := bolt.Open("bolt.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := func(w http.ResponseWriter, r *http.Request) {
		var id uint64
		err := db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("MyBucket"))
			if bucket == nil {
				bucket, err = tx.CreateBucket([]byte("MyBucket"))
			}

			value := bucket.Get([]byte("a"))
			if value == nil {
				value = toBytes(1)
				bucket.Put([]byte("a"), value)
			}

			id = toInt(value)
			err := bucket.Put([]byte("a"), toBytes(id+1))

			return err
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "Hello, you've requested: %s\nBolt: %d", r.URL.Path, id)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func toInt(value []byte) uint64 {
	return binary.BigEndian.Uint64(value)
}

func toBytes(num uint64) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, num)
	return buf.Bytes()
}
