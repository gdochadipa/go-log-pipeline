package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ochadipa/log_pipeline/producer"
	"go.uber.org/zap"
)


func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// This channel will receive the OS signal. setup shutdown channel
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		logger.Fatal("failed start server", zap.Error(err))
	}

	defer listener.Close()

	logger.Info("Log aggregator listening on port 8080")

	concurrencyLimit := 10
	slot := make(chan struct{}, concurrencyLimit)

	var wg sync.WaitGroup

	go func() {
		for{
			conn, err := listener.Accept()
			if err != nil {
				// When we call listener.Close() during shutdown, Accept() will return
				// an error. We check for this specific error to exit the loop cleanly.
				if errors.Is(err, net.ErrClosed) {
					logger.Info("Listener closed, stopping accept loop.")
					break // Exit the loop
				}

				logger.Error("could not accep connection", zap.Error(err))
				continue
			}
			slot <- struct{}{} // add slot
			wg.Add(1)
			go func(c net.Conn){
				defer func() {
					<-slot
					wg.Done()
			 	}() // remove slot
				producer.HandleProducer(conn, logger)
			}(conn)
		}
	}()
	//This is a blocking call. The main function will pause here until a signal is received.
	// If there's no value in the channel, this call will block until another goroutine sends something into the channel.
	// so we need to close this, or there was data inside
	<-shutdownChan
	logger.Info("Shutdown signal received, starting graceful shutdown...")

	//Stop accepting new connections.
	if err := listener.Close(); err != nil {
		logger.Error("Failed to close listener", zap.Error(err))
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	doneChan := make(chan struct{})
	go func(){
		wg.Wait()
		close(doneChan)
	}()

	select {
		case <-doneChan:
			logger.Info("All connections handled gracefully.")
		case <-shutdownCtx.Done():
			log.Println("Shutdown timed out, forcing exit.")
		}

		logger.Info("Server shut down gracefully.")

}
