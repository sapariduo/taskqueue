package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/sapariduo/taskqueue"
)

const (
	defaultNumWorkers     int           = 10
	defaultNumJobs        int           = 10000
	defaultJobDuration    time.Duration = time.Microsecond
	defaultJobExpiration  time.Duration = time.Minute * 5
	defaultReportInterval time.Duration = time.Millisecond * 200
)

var (
	numWorkers     = flag.Int("workers", defaultNumWorkers, "number of workers to spawn")
	numJobs        = flag.Int("jobs", defaultNumJobs, "number of jobs to create")
	jobDuration    = flag.Int64("jobduration", int64(defaultJobDuration), "How long jobs last (time.Sleep")
	jobExpiry      = flag.Int64("expiration", int64(defaultJobExpiration), "How long until a finished job is purged")
	report         = flag.Bool("report", false, "Report on random jobs while jobs are running")
	reportInterval = flag.Int64("reportinterval", int64(defaultReportInterval), "Interval on which to report a random job's status (if report enabled)")
	HTTPAddr       = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

// var (
// 	NWorkers = flag.Int("n", 4, "The number of workers to start")
// 	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
// )

type Queue struct {
	Dispatcher *taskqueue.WorkerDispatcher
}

func NewQueue() *Queue {
	q := new(Queue)
	q.Dispatcher = taskqueue.NewWorkerDispatcher(
		taskqueue.Workers(*numWorkers),
		taskqueue.JobExpiry(time.Duration(*jobExpiry)),
	)
	return q
}

func StartDispatcher() *taskqueue.WorkerDispatcher {
	dispatcher := taskqueue.NewWorkerDispatcher(
		taskqueue.Workers(*numWorkers),
		taskqueue.JobExpiry(time.Duration(*jobExpiry)),
	)
	return dispatcher
}

func (queue *Queue) showJobs(w http.ResponseWriter, r *http.Request) {
	//now you can use users.db
}

func main() {
	// Parse the command-line flags.
	flag.Parse()

	// Start the dispatcher.
	fmt.Println("Starting the dispatcher")
	q := NewQueue()
	// StartDispatcher(*NWorkers)

	// Register our collector as an HTTP handler function.
	fmt.Println("Registering the collector")
	http.HandleFunc("/work", q.AddJob)
	http.HandleFunc("/job", q.GetStatus)

	// Start the HTTP server!
	fmt.Println("HTTP server listening on", *HTTPAddr)
	if err := http.ListenAndServe(*HTTPAddr, nil); err != nil {
		fmt.Println(err.Error())
	}
}
