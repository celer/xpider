package main

import (
	"../robot"
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// Connect to the xpider
	x := &robot.Controller{}
	err := x.Connect("192.168.100.1:80")
	if err != nil {
		panic(err)
	}

	//Set a new random seed
	rand.Seed(time.Now().UTC().UnixNano())

	//Calculate a random direction to walk
	r := rand.Float32()*(3.14*2.0) - 3.14

	// Randomly walk around until we run into something
	for true {

		// Set our front LEDs to be green and red
		x.FrontLED(0, 0xFF, 0, 0xFF, 0, 0)

		// We need to make sure we have a mutex around
		// Getting the status
		state:=x.GetState()

		fmt.Printf("Updated %v\n",state.Updated)
		fmt.Printf("Observed Distance %d\n", state.ObsticalDistance)

		x.AutoMove(100, r, 100, 10)

		time.Sleep(time.Second * 10)

		if !state.Updated.IsZero() && (state.ObsticalDistance < 100) {
			break
		}
		r = rand.Float32()*(3.14*2.0) - 3.14
	}

}
