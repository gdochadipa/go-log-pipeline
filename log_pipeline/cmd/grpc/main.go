package main

import (
	"fmt"
	"log"

	"github.com/ochadipa/log_pipeline/internal"
	"github.com/ochadipa/log_pipeline/storage"
)

func main() {

	var repo internal.LogRepository

	user := "user"
	password := "password"
	host := "host"
	dbName := "dbName"

	db, err := storage.NewPostgressConnection(&user, &password, &host, &dbName)

	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}


	repo = internal.NewRepository(db)
	defer db.Close()

	log.Println("Listening on port 8080...")

	service := internal.NewLogService(repo)

	log.Fatal(internal.ListenGRPC(service,8080))


}
