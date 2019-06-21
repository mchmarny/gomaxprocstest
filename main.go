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
	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"cpu_num": numOfCores,
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
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Argument",
			"status":  "BadRequest",
		})
		return
	}

	r := runCores(cc)
	c.IndentedJSON(http.StatusOK, r)
	return
}

// e := json.NewEncoder(w)
// e.SetEscapeHTML(true)
// e.SetIndent("", "\t")
// e.Encode(o)

func runCores(n int) *runResponse {

	runtime.GOMAXPROCS(n)
	logger.Printf("Setting max CPU: %d/%d", n, numOfCores)
	start := time.Now()

	r := &runResponse{
		Messages:    []runMsg{},
		TottalCores: numOfCores,
		MaxCores:    n,
	}

	done := make(chan int)
	x := int64(0)

	for i := 1; i <= n; i++ {
		logger.Printf("Core %d start", i)
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
			r.add(k, fmt.Sprintf("Done: %s", e))
			doneWorkerCount++
			if doneWorkerCount == n {
				break W
			}
		}
	}

	r.TotalDuration = time.Since(start).String()
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
	TottalCores   int      `json:"total_cores"`
	MaxCores      int      `json:"max_cores"`
	TotalDuration string   `json:"duration"`
	Messages      []runMsg `json:"messages"`
}

func (r *runResponse) add(c int, m string) {
	r.Messages = append(r.Messages, runMsg{
		CoreIndex: c,
		Message:   m,
	})
}

type runMsg struct {
	CoreIndex int    `json:"core"`
	Message   string `json:"message"`
}
