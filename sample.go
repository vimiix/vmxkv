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
    db.Put(1, 1)
    db.Put(2, 2)
    db.Put(3, 3)
    v, err := db.Get(2)
    if err != nil {
    }
    fmt.Println("get:", v)
}
