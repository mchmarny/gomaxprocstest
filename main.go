package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	logger     = log.New(os.Stdout, "[app] ", 0)
	numOfCores = runtime.NumCPU()
)

func main() {

	logger.Printf("CPU cores: %d", numOfCores)
	runtime.GOMAXPROCS(numOfCores)

	mux := http.NewServeMux()

	// Handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})
	mux.HandleFunc("/run", runHandler)

	// Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	logger.Printf("Server starting on port %s \n", port)
	logger.Fatal(server.ListenAndServe())

}

func runHandler(w http.ResponseWriter, r *http.Request) {

	logger.Printf("CPU cores: %d", numOfCores)
	fmt.Fprintf(w, "CPU cores: %d\n", numOfCores)

	start := time.Now()
	done := make(chan int)
	x := int64(0)

	for i := 0; i < numOfCores; i++ {
		logger.Printf("Core %d start", i)
		fmt.Fprintf(w, "Core %d start\n", i)
		go func(worker int) {
			doWork(&x)
			done <- worker
		}(i)
	}

	doneWorkerCount := 0
W:
	for {
		select {
		case k := <-done:
			e := time.Since(start)
			logger.Printf("Core %d done in %s", k, e)
			fmt.Fprintf(w, "Core %d done in %s\n", k, e)
			doneWorkerCount++
			if doneWorkerCount == numOfCores {
				break W
			}
		}
	}

	elapsed := time.Since(start)
	logger.Printf("Total duration: %s", elapsed)
	fmt.Fprintf(w, "Total duration: %s\n", elapsed)

}

func doWork(p *int64) {
	for i := int64(1); i <= 9000000000; i++ {
		*p = i
	}
}
