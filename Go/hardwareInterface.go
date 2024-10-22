package main

import (
	"fmt"
	"os"
	"time"
	"os/signal"
	"syscall"
	"github.com/stianeikeland/go-rpio/v4"
)

const P_WAVE_PERCENTAGE = .20
const R_WAVE_PERCENTAGE = .10
const REST_PERCENTAGE = .70
const MAX_BRIGHTNESS = 100

type Hardware struct {
	beat        int
	audioFile   string
	pin         rpio.Pin
	dataChannel chan int
}

func (h *Hardware) init(audioFile string, dataChannel chan int) {
	h.pin = rpio.Pin(13)
	h.audioFile = audioFile
	h.beat = 60
	h.dataChannel = dataChannel
	fmt.Println("init hardware")
}

func (h *Hardware) loop_test() {
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
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()
	h.pin.Mode(rpio.Pwm)
	h.pin.Freq(1000)
	h.pin.DutyCycle(0, MAX_BRIGHTNESS)

	for {
		select {
		case beat, ok := <-h.dataChannel: // Wait for the trigger from the sender
			if !ok {
				fmt.Println("Channel closed, stopping fading...")
				return
			}
			h.beat = beat
		case <-sigs:
			rpio.Close()
			fmt.Println("HANDLING SIG")
		default:
			fmt.Println("Starting fade cycle...")
			speed := float64(60000 / h.beat)            // speed in ms where 60000 = 60 seconds so 70bpm will equal 60000 / 70 = 857ms per beat
			pWave := float64(speed * P_WAVE_PERCENTAGE) // P wave is the initial ba sound
			rWave := float64(speed * R_WAVE_PERCENTAGE)
			rest := speed * REST_PERCENTAGE

			pWaveBrightness := uint32(.40 * float64(MAX_BRIGHTNESS))
			pWaveFadeInLength := uint32(pWave * .70)
			pWaveFadeOutLength := uint32(pWave * .30)

			h.fadeInAndOut(pWaveFadeInLength, pWaveFadeOutLength, pWaveBrightness)

			rWaveBrightness := uint32(MAX_BRIGHTNESS)
			rWaveFadeInLength := uint32(rWave * .50)
			rWaveFadeOutLength := uint32(rWave * .50)

			h.fadeInAndOut(rWaveFadeInLength, rWaveFadeOutLength, rWaveBrightness)

			if h.beat == 0 {
				// make sure there's a reading
				time.Sleep(time.Second * 1)
			} else {
				time.Sleep(time.Millisecond * time.Duration(rest))
			}
		}
	}
}

func (h *Hardware) fadeInAndOut(fadeInLength, fadeOutLength, brightness uint32) {
	for i := uint32(0); i <= fadeInLength; i++ { // increasing brightness
		fadeInRate := i * (brightness / fadeInLength)
		h.pin.DutyCycle(fadeInRate, brightness)
		fmt.Println("Fade in rate: ", fadeInRate, brightness)
	}
	for i := uint32(0); i <= fadeOutLength; i++ { // decreasing brightness
		fadeOutRate := i * (brightness / fadeOutLength)
		h.pin.DutyCycle(brightness-fadeOutRate, brightness)
		fmt.Println("Fade out rate: ", fadeOutRate, brightness)
	}
}
