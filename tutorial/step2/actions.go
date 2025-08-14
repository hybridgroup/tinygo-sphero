package main

import (
	"image/color"
	"time"
)

func actions() {
	println("set led red")
	err := robot.SetLEDColor(color.RGBA{R: 255, G: 0, B: 0})
	if err != nil {
		println(err)
	}

	println("rolling forward")
	robot.Roll(0, 150)

	time.Sleep(2 * time.Second)

	println("rolling back")
	robot.Roll(180, 150)

	time.Sleep(2 * time.Second)

	println("going to sleep")
	robot.Sleep()
}
