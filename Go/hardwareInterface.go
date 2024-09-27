package main

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

type Hardware struct {
	beat        int
	audioFile   string
	pin         rpio.Pin
	dataChannel chan int
}

func (h *Hardware) init(audioFile string, dataChannel chan int) {
	h.pin = rpio.Pin(10)
	h.audioFile = audioFile
	h.beat = 60
	h.dataChannel = dataChannel
}

func (h *Hardware) function_loop() {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	h.pin.Output()

	// Toggle pin 20 times
	for {
		h.pin.Toggle()
		time.Sleep(time.Second / 5)
	}
}

func (h *Hardware) ReceiveAndFade() {
	ticker := time.NewTicker(100 * time.Millisecond) // Print every 100ms
	defer ticker.Stop()

	for {
		select {
		case beat, ok := <-h.dataChannel: // Wait for the trigger from the sender
			if !ok {
				fmt.Println("Channel closed, stopping fading...")
				return
			}
			h.beat = beat

		case <-ticker.C:
			fmt.Println("Starting fade cycle...")
			// print(bpm)
			speed := 60 / h.beat
			brighter_time := speed / 2 // Spend half a beat getting brighter
			dimmer_time := speed / 2   // Spend half a beat getting dimmer
			// Smoothly fade in and out
			for i := 0; i <= brighter_time; i++ { // increasing brightness
				r.pin.DutyCycle(i, 32)
				time.Sleep(time.Second / brighter_time)
			}
			for i := dimmer_time; i > 0; i-- { // decreasing brightness
				r.pin.DutyCycle(i, dimmer_time)
				time.Sleep(time.Second / dimmer_time)
			}
		}
	}
}
