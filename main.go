package main

func main() {
	InitLoger()
	store, err := NewPostgresStore()
	if err != nil {
		FatalLogger.Println("Error creating the store: ", err)
	}

	if err := store.Init(); err != nil { // this will create the table if it does not exist
		FatalLogger.Println("Error initializing the store: ", err)
	}

	apiServer := NewApiServer(":8080", store)
	apiServer.Run()
}
