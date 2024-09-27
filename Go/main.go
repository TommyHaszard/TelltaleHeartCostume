package main

import (
	"log"
	"net/http"
)

func main() {
	channel := make(chan int)
	heartbeat := Heartbeat{}
	hardwareInterface := Hardware{}
	hardwareInterface.init("string", channel)
	heartbeat.init(channel)
	http.Handle("/heartbeat", heartbeat)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("192.168.1.104:8080", nil))
}
