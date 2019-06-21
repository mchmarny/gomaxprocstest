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
	defaultCoreNum   = "1"
)

var (
	logger     = log.New(os.Stdout, "[app] ", 0)
	numOfCores = runtime.NumCPU()
)

func main() {

	logger.Printf("CPU cores: %d", numOfCores)
	gin.SetMode(gin.ReleaseMode)

	// router
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/", homeHandler)
	r.GET("/cores/:cores", coreHandler)
	r.GET("/cores", coreHandler)

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
	c.JSON(http.StatusOK, map[string]interface{}{
		"numCPU": numOfCores,
	})
}

func coreHandler(c *gin.Context) {
	k := c.Param("cores")
	if k == "" {
		k = defaultCoreNum
	}

	cc, err := strconv.Atoi(k)
	if err != nil {
		logger.Printf("Error while parsing core parameter: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Argument",
			"status":  "BadRequest",
		})
		return
	}

	r := runCores(cc)
	c.JSON(http.StatusOK, r)
	return
}

const numOfLoops = 10

func runCores(n int) *runResponse {

	runtime.GOMAXPROCS(n)
	logger.Printf("Setting max CPU: %d/%d", n, numOfCores)

	r := &runResponse{
		StartTime:   time.Now(),
		Messages:    []runMsg{},
		TottalCores: numOfCores,
		UsedCores:   n,
	}

	done := make(chan int)
	x := int64(0)

	for i := 0; i < numOfLoops; i++ {
		logger.Printf("Core %d start", i)
		r.add(i, fmt.Sprintf("Core %d start", i))
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
			e := time.Since(r.StartTime)
			logger.Printf("Core %d done in %s", k, e)
			r.add(k, fmt.Sprintf("Core %d done in %s", k, e))
			doneWorkerCount++
			if doneWorkerCount == n {
				break W
			}
		}
	}

	r.EndTime = time.Now()
	r.TotalDuration = time.Since(r.StartTime).String()
	logger.Printf("Total duration: %s", r.TotalDuration)

	return r

}

const workLoops = 100000000

func doWork(p *int64) {
	for i := int64(1); i <= workLoops; i++ {
		*p = i
	}
}

type runResponse struct {
	TottalCores   int       `json:"tottalCores"`
	UsedCores     int       `json:"usedCores"`
	RunCores      int       `json:"runCores"`
	NumOfLoops    int       `json:"numOfLoops"`
	Messages      []runMsg  `json:"messages"`
	StartTime     time.Time `json:"startTime"`
	EndTime       time.Time `json:"endTime"`
	TotalDuration string    `json:"totalDuration"`
}

func (r *runResponse) add(c int, m string) {
	r.Messages = append(r.Messages, runMsg{
		CoreIndex: c,
		Message:   m,
	})
}

type runMsg struct {
	CoreIndex int    `json:"coreIndex"`
	Message   string `json:"message"`
}
