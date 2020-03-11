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
			if (bucket == nil) {
				bucket, err = tx.CreateBucket([]byte("MyBucket"))
			}

			value := bucket.Get([]byte("a"))
			var num uint64 = 1
			if (value == nil) {
				buf := new(bytes.Buffer)
				binary.Write(buf, binary.BigEndian, num)
				value = buf.Bytes()
				bucket.Put([]byte("a"), value)
			}

			id = binary.BigEndian.Uint64(value)
			num = id + 1
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.BigEndian, num)

			err := bucket.Put([]byte("a"), buf.Bytes())

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
