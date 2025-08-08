package main

import (
	"image/color"
	"strconv"
	"time"

	sphero "github.com/hybridgroup/tinygo-sphero"
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	robot *sphero.Robot

	debug = false
)

func main() {
	time.Sleep(5 * time.Second)
	println("enabling...")

	must("enable BLE interface", adapter.Enable())

	println("start scan...")

	must("start scan", adapter.Scan(scanHandler))

	var err error
	select {
	case result := <-ch:
		device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
		must("connect to peripheral device", err)

		println("connected to ", result.Address.String())
	}

	defer device.Disconnect()

	robot = sphero.NewRobot(&device)
	err = robot.Start()
	if err != nil {
		println("start error: ", err.Error())
		return
	}

	robot.EnableNotifications(func(msg *sphero.Payload) {
		switch {
		case msg.DeviceID == sphero.DeviceSensor && msg.Command == sphero.SensorCommandCollisionDetectedAsync:
			data := msg.Payload
			if len(data) == 0 {
				if debug {
					println("invalid collision notification")
				}

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
		default:
			if debug {
				println("data:", msg.String())
			}
		}
	})

	robot.ConfigureCollisionDetection(sphero.CollisionConfig{
		Method: 1,
		Xt:     20,
		Yt:     20,
		Xs:     20,
		Ys:     20,
		Dead:   10})

	robot.GetBatteryVoltage()

	println("set led red")
	err = robot.SetLEDColor(color.RGBA{R: 255, G: 0, B: 0})
	if err != nil {
		println(err)
	}

	time.Sleep(time.Second)

	// roll
	for i := 0; i < 4; i++ {
		robot.Roll(i*90, 150)
		time.Sleep(1 * time.Second)
	}

	robot.Stop()
	robot.SetLEDColor(color.RGBA{R: 0, G: 0, B: 255})
	time.Sleep(3 * time.Second)

	robot.Sleep()
}

func scanHandler(a *bluetooth.Adapter, d bluetooth.ScanResult) {
	println("device:", d.Address.String(), d.RSSI, d.LocalName())
	if d.Address.String() == connectAddress() {
		a.StopScan()
		ch <- d
	}
}

func must(action string, err error) {
	if err != nil {
		for {
			println("failed to " + action + ": " + err.Error())
			time.Sleep(time.Second)
		}
	}
}
