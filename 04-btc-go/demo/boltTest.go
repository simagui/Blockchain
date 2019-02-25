package main

import (
	"lib/bolt"
	"fmt"
)

const (
	daName      = "test.db"
	bucketName1 = "bucketName1"
)

func main() {
	db, err := bolt.Open(daName, 0600, nil)
	if err != nil {
		panic(err)

	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b:=tx.Bucket([]byte(bucketName1))
	    if b==nil{
	    	b,err=b.CreateBucket([]byte(bucketName1))
		    if err!=nil{
		    	panic(err)
			}
	    	}
		b.Put([]byte("111"),[]byte("hello"))
		v1:=b.Get([]byte("111"))
		fmt.Print(v1)
    return nil
	})
}
