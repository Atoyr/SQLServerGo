package main

import "github.com/boltdb/bolt"

func openDB (path string) (*bolt.DB, error){
    return bolt.Open(path, 0600, nil)
}

func update (db *bolt.DB, bucket string, key, value []byte) {
  db.Update(func(tx *bolt.Tx) error {
        bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
        if err != nil {
          return err
        }
        err = bucket.Put(key,value)
        if err != nil {
          return err
        }
        return nil
    })
}


