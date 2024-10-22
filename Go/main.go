package main

import (
	"log"
	"net/http"
)

func main() {
	channel := make(chan int, 1)
	heartbeat := Heartbeat{}
	hardwareInterface := Hardware{}
	hardwareInterface.init("string", channel)
	heartbeat.init(channel)
	go hardwareInterface.ReceiveAndFade()
	// go hardwareInterface.loop_test()
	http.Handle("/heartbeat", heartbeat)
	log.Print("Starting server...")
	log.Fatal(http.ListenAndServe("192.168.1.119:8080", nil))
}
