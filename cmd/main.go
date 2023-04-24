package main

import (
	"log"

	"key-value-ts/domains/storage/files"
	"key-value-ts/servers"
)

func main() {
	storer, err := files.New("store")
	if err != nil {
		log.Fatal(err)
	}
	server := servers.New(storer)
	log.Fatal(server.Start(":8080"))
}
