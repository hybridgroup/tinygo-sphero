package main

import (
	"time"

	sphero "github.com/hybridgroup/tinygo-sphero"
	"tinygo.org/x/bluetooth"
)

var (
	adapter = bluetooth.DefaultAdapter
	device  bluetooth.Device
	ch      = make(chan bluetooth.ScanResult, 1)

	robot *sphero.Robot
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

		println("connected to", result.Address.String())
	}

	defer device.Disconnect()

	robot = sphero.NewRobot(&device)
	err = robot.Start()
	if err != nil {
		failMessage("start error: " + err.Error())
		return
	}

	actions()

	done()
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
		failMessage("failed to " + action + ": " + err.Error())
	}
}
