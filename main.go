package main

import (
	"log"
)

func main() {
	store, err := NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer("127.0.0.1:8080", store)
	server.Run()
}
