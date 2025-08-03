package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ochadipa/log_pipeline/internal"
	"github.com/ochadipa/log_pipeline/internal/db"
	"github.com/ochadipa/log_pipeline/processor"
)

func main() {

	var repo internal.LogRepository

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := db.NewPostgressConnection(&user, &password, &host, &dbName)
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	workerPool := processor.NewWorkerPool(db, 3, 100)

	defer workerPool.Shutdown(ctx)
	defer cancel()

	repo = internal.NewRepository(ctx, workerPool)

	log.Println("Listening on port 50051...")

	service := internal.NewLogService(repo)

	// running grpc service
	log.Fatal(internal.ListenGRPC(service, 50051))
}

// another option of main, incase if you want run with go grpc many routine
func exampleMain() {
	var repo internal.LogRepository

	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	db, err := db.NewPostgressConnection(&user, &password, &host, &dbName)
	sqlDB, err := db.DB()
	defer sqlDB.Close()

	if err != nil {
		fmt.Sprintln("Database load failed %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	workerPool := processor.NewWorkerPool(db, 3, 100)
	defer workerPool.Shutdown(ctx)

	repo = internal.NewRepository(ctx, workerPool)

	log.Println("Listening on port 50051...")

	service := internal.NewLogService(repo)

	signalChan := make(chan os.Signal, 1)
	// notify by syscall.SIGINT and syscall.SIGTERM, trigger by os
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal ...")
		cancel()
	}()

	if err := internal.ListenGRPC2(ctx, service, 50051); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}

	<-ctx.Done()
	log.Println("Application shutting down...")
}
