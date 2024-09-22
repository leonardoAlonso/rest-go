package main

import "log"

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil { // this will create the table if it does not exist
		log.Fatal(err)
	}

	apiServer := NewApiServer(":8080", store)
	apiServer.Run()
}
