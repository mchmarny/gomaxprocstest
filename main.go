package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

var (
	logger     = log.New(os.Stdout, "[app] ", 0)
	numOfCores = runtime.NumCPU()
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", homeHandler)
	r.GET("/cores/:core/concurrency/:count/calcs/:calc", workHandler)

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	logger.Printf("Server starting: %s \n", addr)
	if err := r.Run(addr); err != nil {
		logger.Fatal(err)
	}

}

func homeHandler(c *gin.Context) {
	u := fmt.Sprintf("/cores/%d/concurrency/%d/calcs/1000000000",
		numOfCores, numOfCores)
	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"cpu_num": numOfCores,
		"example": u,
	})
}

func workHandler(c *gin.Context) {

	// URL params
	cores := paramAsInt(c, "core")
	counts := paramAsInt(c, "count")
	calcs := paramAsInt64(c, "calc")
	r := runWork(cores, counts, calcs)
	c.IndentedJSON(http.StatusOK, r)
	return
}

func paramAsStr(c *gin.Context, k string) string {
	logger.Printf("Parsing '%s'", k)
	a := c.Param(k)
	if a == "" {
		logger.Fatalf("Parameter not defined '%s'", k)
	}
	return a
}

func paramAsInt64(c *gin.Context, k string) int64 {
	a := paramAsStr(c, k)
	i, err := strconv.ParseInt(a, 10, 64)
	if err != nil {
		logger.Fatalf("Error while parsing %s parameter: %v", k, err)
	}
	return i
}

func paramAsInt(c *gin.Context, k string) int {
	a := paramAsStr(c, k)
	i, err := strconv.Atoi(a)
	if err != nil {
		logger.Fatalf("Error while parsing %s parameter: %v", k, err)
	}
	return i
}

func runWork(cores, counts int, calcs int64) *runResponse {

	logger.Printf("Running cores:%d counts:%d calcs:%d", cores, counts, calcs)

	runtime.GOMAXPROCS(int(cores))
	logger.Printf("Setting max CPU: %d/%d", cores, numOfCores)
	start := time.Now()

	r := &runResponse{
		AvailableCores: numOfCores,
		MaxCores:       cores,
		Concurrency:    counts,
		Calculations:   calcs,
		Details:        []calcDetail{},
	}

	done := make(chan int)

	for i := 1; i <= int(counts); i++ {
		logger.Printf("Core %d start", i)
		go func(worker int) {
			doWork(calcs)
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
			r.add(k, e.String())
			doneWorkerCount++
			if doneWorkerCount == int(counts) {
				break W
			}
		}
	}

	r.TotalDuration = time.Since(start).String()
	logger.Printf("Total duration: %s", r.TotalDuration)

	return r

}

func doWork(c int64) {
	p := int64(0)
	for i := int64(1); i <= c; i++ {
		p = p + i
		p = p - i
	}
}

type runResponse struct {
	AvailableCores int          `json:"available_cores"`
	MaxCores       int          `json:"max_cores"`
	Concurrency    int          `json:"concurrency"`
	Calculations   int64        `json:"calculations"`
	TotalDuration  string       `json:"duration"`
	Details        []calcDetail `json:"details"`
}

func (r *runResponse) add(c int, m string) {
	r.Details = append(r.Details, calcDetail{
		Routine:  c,
		Duration: m,
	})
}

type calcDetail struct {
	Routine  int    `json:"goroutine"`
	Duration string `json:"duration"`
}
