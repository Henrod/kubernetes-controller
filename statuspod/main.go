package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Pod statuses
const (
	ready      = "ready"
	occupied   = "occupied"
	terminated = "terminated"
)

var (
	status   = ready
	url      string
	wg       sync.WaitGroup
	hostname string
)

func getHandler(w http.ResponseWriter, r *http.Request) {
	bts, err := json.Marshal(map[string]interface{}{
		"name":   hostname,
		"status": status,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(bts)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("changing status, current:", status)
	if status == ready {
		status = occupied
	} else {
		status = ready
	}

	sendStatus()
}

func bodyReader() io.Reader {
	type body struct {
		Status string
	}
	bts, _ := json.Marshal(&body{status})
	return bytes.NewReader(bts)
}

func sendStatus() {
	log.Println("sending request to", url, "status", status)

	fullURL := url + "/statuss"

	client := &http.Client{}
	resp, err := client.Post(fullURL, "application/json", bodyReader())
	if err != nil {
		log.Printf("error: %q\n", err)
		return
	}

	defer resp.Body.Close()
	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error: %q\n", err)
		return
	}

	log.Printf("server response: %s", []byte(bts))
}

func keepSendingStatus() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sendStatus()
		}
	}
}

func captureSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	status = terminated
	sendStatus()
	wg.Done()
}

func main() {
	log.Println("starting server")
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/toggle", toggleHandler)

	flag.StringVar(&url, "url", "http://localhost:8080", "api url to send status")
	flag.Parse()

	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	println("henrod", hostname)

	wg.Add(1)
	go captureSignal()
	go keepSendingStatus()
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	wg.Wait()
}
