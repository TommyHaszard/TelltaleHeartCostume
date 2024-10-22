package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
)

type HeartbeatJson struct {
	HeartBeat *string `json:"HeartBeat"`
}

type Heartbeat struct {
	currHeartbeat     int
	audioFile         string
	trailingHeartbeat int
	averageOver       int
	dataChannel       chan int
}

func (heartbeat *Heartbeat) init(dataChannel chan int) {
	// get audio file and set that up
	heartbeat.averageOver = 10
	heartbeat.dataChannel = dataChannel
	heartbeat.currHeartbeat = 60
}

func (h *Heartbeat) recieveHeartbeat(value int) {
	if h.averageOver == 0 {
		h.averageOver = 10
	}
	// get the trailingAvg
	h.trailingHeartbeat = value
	// if trailing is significantly different then change currHeartbeat
	if float64(h.trailingHeartbeat)*1.2 > float64(h.currHeartbeat) || float64(h.trailingHeartbeat)*.7 < float64(h.currHeartbeat) {
		fmt.Println("Heartbeat received from Watch: ", h.trailingHeartbeat)
		h.currHeartbeat = h.trailingHeartbeat
		h.dataChannel <- h.currHeartbeat
	}
}

func (heartbeat Heartbeat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var t HeartbeatJson
	err := decoder.Decode(&t)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t.HeartBeat == nil {
		http.Error(w, "missing field 'HeartBeat' from JSON object", http.StatusBadRequest)
		return
	}

	i, err := strconv.Atoi(*t.HeartBeat)
	if err != nil {
		// ... handle error
		http.Error(w, "Could not convert heartbeat str to int", http.StatusBadRequest)
		return
	}

	heartbeat.recieveHeartbeat(i)
}
