package main

import (
	"fmt"
	"log"

	"github.com/vimiix/vmxkv/server"
)

func main() {
	db, err := server.OpenDB("vmxkv.db", 0666)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed opening db: %s", err))
	}
	defer db.Close()
	//db.Put(1, 1)
	//db.Put(2, 2)
	//db.Put(3, 3)
	v, err := db.Get(2)
	if err != nil {
	}
	fmt.Println("get:", v)
	fmt.Println("range find:")
	db.List(1, 3, func(k, v uint64) bool {
		fmt.Println(k, v)
		return true
	})
	db.Del(3)
	fmt.Println(db.Get(3))
	// Output
	// get: 2
	// range find:
	// 1 1
	// 2 2
	// 3 3
	// 0 not found
}
