package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/ochadipa/goroutine-practice/models"
	"github.com/ochadipa/goroutine-practice/processor"
	"github.com/ochadipa/goroutine-practice/producer"
)


func main() {

runtime.GOMAXPROCS(4)

// defer when need cancel context in background
ctx, cancel := context.WithCancel(context.Background())

// when all done, lets cancel the run
defer cancel();
// raw and error channel for streamline to get logs from each server
rawLogChan := make(chan models.RawLog)
errChan := make(chan models.ErrorLog)

// parsed channel is the main channel, after exctract from raw log and err channels
resParsedChan := make(chan models.ParsedLog)

// Handle OS signals for graceful shutdown
// like press Ctrl+C to exit.
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt)

var wg sync.WaitGroup

// there will be two gorutine, so preapre 2 waitgroup
// It's crucial to call Add for *each* goroutine before starting it.
wg.Add(2)

// example to get log from ServerA
go func(){
	defer wg.Done()// when done, set the goroutine are done.
	producer.GenerateRawLogs(ctx, "ServerA", rawLogChan, errChan)
}()

// example to get log from ServerB
// Server A -> sendlog -> logPipeline Services
go func(){
	defer wg.Done()// when done, set the goroutine are done.
	producer.GenerateRawLogs(ctx, "ServerB", rawLogChan, errChan)
}()

// 2. Start Processor Worker Pool (3 workers)
for i := range 3{
	wg.Add(1)
	go func(id int){
		defer wg.Done()// when done, set the goroutine are done.
		processor.ParsedLogs(ctx, id, rawLogChan,resParsedChan, errChan)
	}(i)
}


// 3. Start Storage Consumer
wg.Add(1) // tell to WG, there will a new gouritne to add
go func() {
	defer wg.Done()// when done, set the goroutine are done.
 // TODO: implement storage of parsed logs
        // storeLogs(ctx, parsedLogCh, errorCh)
}()

// monitor / shutdown log
go func() {
	select {
		case sig := <-sigCh:
			fmt.Printf("Received signal: %v, shutting down...\n",sig)
			cancel()
		case err := <- errChan:
			fmt.Printf("Error occurred: %v, shutting down...\n", err)
            cancel()
        case <-time.After(10 * time.Second):
            fmt.Println("Time limit reached, shutting down...")
            cancel()
	}
}()


// waiting another wait group / gouritne until done `wg.Done()`
wg.Wait()

//close all channel
close(rawLogChan)
close(errChan)
close(resParsedChan)

fmt.Println("Pipeline stopped gracefully.")

}
