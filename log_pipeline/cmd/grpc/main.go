package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ochadipa/log_pipeline/internal"
	"github.com/ochadipa/log_pipeline/processor"
	"github.com/ochadipa/log_pipeline/storage"
)

func main() {

	var repo internal.LogRepository

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := storage.NewPostgressConnection(&user, &password, &host, &dbName)

	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}

	workerPool := processor.NewWorkerPool(db, 3, 1000)

	repo = internal.NewRepository(workerPool)
	defer db.Close()

	log.Println("Listening on port 50051...")

	service := internal.NewLogService(repo)

	log.Fatal(internal.ListenGRPC(service, 50051))
}
