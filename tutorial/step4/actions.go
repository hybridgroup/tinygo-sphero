package main

import (
	"image/color"
	"time"
)

func actions() {
	println("set led red")
	robot.SetLEDColor(color.RGBA{R: 255, G: 0, B: 0})
	robot.Roll(60, 150)
	time.Sleep(1 * time.Second)

	println("set led green")
	robot.SetLEDColor(color.RGBA{R: 0, G: 255, B: 0})
	robot.Roll(150, 150)
	time.Sleep(1 * time.Second)

	println("set led blue")
	robot.SetLEDColor(color.RGBA{R: 0, G: 0, B: 255})
	robot.Roll(270, 150)
	time.Sleep(1 * time.Second)

	println("set led purple")
	robot.SetLEDColor(color.RGBA{R: 128, G: 0, B: 128})
	robot.Stop()

	time.Sleep(2 * time.Second)

	println("going to sleep")
	robot.Sleep()
}
