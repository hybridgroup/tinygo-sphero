package main

import (
	"image/color"
	"strconv"
	"time"

	sphero "github.com/hybridgroup/tinygo-sphero"
)

var (
	red  = color.RGBA{R: 255, G: 0, B: 0}
	blue = color.RGBA{R: 0, G: 0, B: 255}

	currentColor = red
)

func actions() {
	robot.ConfigureCollisionDetection(sphero.CollisionConfig{
		Method: 1,
		Xt:     20,
		Yt:     20,
		Xs:     20,
		Ys:     20,
		Dead:   10})

	robot.EnableNotifications(func(msg *sphero.Payload) {
		switch {
		case msg.DeviceID == sphero.DeviceSensor && msg.Command == sphero.SensorCommandCollisionDetectedAsync:
			data := msg.Payload
			if len(data) == 0 {
				println("invalid collision notification")
				return
			}
			axis, yMag, xMag := data[6], data[8], data[10]
			direction := "forward/backward"
			if axis == 1 {
				direction = "left/right"
			}

			println("collision", direction, "yMag", yMag, "xMag", xMag)
		case msg.DeviceID == sphero.DevicePowerInfo && msg.Command == sphero.PowerCommandsBatteryVoltage:
			data := msg.Payload
			voltage := float32(int32(data[2]) + int32(data[1])*256 + int32(data[0])*65536)
			println("battery:", strconv.FormatFloat(float64(voltage)/100, 'f', 2, 64), "V")
		}
	})

	for range 30 {
		if currentColor == red {
			currentColor = blue
		} else {
			currentColor = red
		}

		robot.SetLEDColor(currentColor)
		robot.GetBatteryVoltage()
		time.Sleep(time.Second)
	}

	println("going to sleep")
	robot.Sleep()
}
