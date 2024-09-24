package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Heartbeat struct {
	HeartBeat *string `json:"HeartBeat"`
}

type HeartbeatHandler struct {
}

func (wsh HeartbeatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var t Heartbeat
	err := decoder.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.HeartBeat == nil {
		http.Error(w, "missing field 'test' from JSON object", http.StatusBadRequest)
		return
	}

	log.Println(*t.HeartBeat)

}

func main() {
	heartbeatHandler := HeartbeatHandler{}
	http.Handle("/heartbeat", heartbeatHandler)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("192.168.1.104:8080", nil))
}
