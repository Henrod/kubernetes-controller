package http

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// StatusUpdater updates the status of a pod
type StatusUpdater interface {
	Update(name, status string)
}

type updateHandler struct {
	statusUpdater StatusUpdater
}

type body struct {
	Status string
	Name   string
}

func (u *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("called http api")

	defer r.Body.Close()
	bts, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("err %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := new(body)
	err = json.Unmarshal(bts, body)
	if err != nil {
		log.Printf("err %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("pod %s has status %s\n", body.Name, body.Status)

	u.statusUpdater.Update(body.Name, body.Status)
	w.WriteHeader(http.StatusOK)
}

// Start starts a server
func Start(statusUpdater StatusUpdater) {
	log.Println("starting server")

	http.Handle("/statuss", &updateHandler{statusUpdater})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
