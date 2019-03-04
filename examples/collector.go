package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// // A buffered channel that we can send work requests on.
// var WorkQueue = make(chan WorkRequest, 100)

func (queue *Queue) AddJob(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP POST request.
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Parse the delay.
	delay, err := time.ParseDuration(r.FormValue("delay"))
	if err != nil {
		http.Error(w, "Bad delay value: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Check to make sure the delay is anywhere from 1 to 10 seconds.
	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		http.Error(w, "The delay must be between 1 and 10 seconds, inclusively.", http.StatusBadRequest)
		return
	}

	// Now, we retrieve the person's name from the request.
	name := r.FormValue("name")

	// Just do a quick bit of sanity checking to make sure the client actually provided us with a name.
	if name == "" {
		http.Error(w, "You must specify a name.", http.StatusBadRequest)
		return
	}

	jobFunc := func() error {

		time.Sleep(time.Duration(delay))
		fmt.Println(name)

		return nil
	}

	thisjob := queue.Dispatcher.QueueFunc(jobFunc)
	id := thisjob.Status()
	status, _ := json.Marshal(id)

	// // Now, we take the delay, and the person's name, and make a WorkRequest out of them.
	// work := WorkRequest{Name: name, Delay: delay}

	// // Push the work onto the queue.
	// WorkQueue <- work
	fmt.Println("Work request queued with ID", id)

	// And let the user know their work request was created.
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(status))
	return
}

func (queue *Queue) GetStatus(w http.ResponseWriter, r *http.Request) {
	// Make sure we can only be called with an HTTP GET request.
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//Parse Job ID
	id := r.URL.Query().Get("JobID")
	jobid, _ := strconv.ParseUint(id, 10, 64)

	tracker, _ := queue.Dispatcher.JobStatus(uint(jobid))
	status, _ := json.Marshal(tracker.Status())

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(status))
	return
}
