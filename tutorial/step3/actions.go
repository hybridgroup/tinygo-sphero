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

	println("rolling in a square")

	for i := 0; i < 4; i++ {
		robot.Roll(i*90, 150)
		time.Sleep(1 * time.Second)
	}

	println("set led blue")
	robot.SetLEDColor(color.RGBA{R: 0, G: 0, B: 255})
	time.Sleep(2 * time.Second)

	println("going to sleep")
	robot.Sleep()
}
