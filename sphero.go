package sphero

import (
	"errors"
	"image/color"
	"time"

	"tinygo.org/x/bluetooth"
)

type Robot struct {
	device            *bluetooth.Device
	apiService        *bluetooth.DeviceService
	apiCharacteristic *bluetooth.DeviceCharacteristic

	spheroService         *bluetooth.DeviceService
	antiDOSCharacteristic *bluetooth.DeviceCharacteristic

	buf                     []byte
	sequenceNo              int
	expectedCommandSequence int
}

var (
	APIServiceUUID, _ = bluetooth.ParseUUID("00010001-574f-4f20-5370-6865726f2121")
	ApiV2UUID, _      = bluetooth.ParseUUID("00010002-574f-4f20-5370-6865726f2121")

	SpheroServiceUUID, _ = bluetooth.ParseUUID("00020001-574f-4f20-5370-6865726f2121")
	DFUUUID, _           = bluetooth.ParseUUID("00020002-574f-4f20-5370-6865726f2121")
	DFUv2UUID, _         = bluetooth.ParseUUID("00020004-574f-4f20-5370-6865726f2121")
	AntiDOSUUID, _       = bluetooth.ParseUUID("00020005-574f-4f20-5370-6865726f2121")
)

// NewRobot creates a new Sphero Mini robot.
func NewRobot(dev *bluetooth.Device) *Robot {
	r := &Robot{
		device: dev,
		buf:    make([]byte, 0, 255),
	}

	return r
}

func (r *Robot) Start() (err error) {
	srvcs, err := r.device.DiscoverServices([]bluetooth.UUID{
		APIServiceUUID,
		SpheroServiceUUID,
	})
	if err != nil || len(srvcs) == 0 {
		return errors.New("could not find services")
	}

	// TODO: look into why this is reversed
	r.spheroService = &srvcs[0]
	r.apiService = &srvcs[1]

	if debug {
		println("found API service", r.apiService.UUID().String())
		println("found Sphero service", r.spheroService.UUID().String())
	}

	chars, err := r.apiService.DiscoverCharacteristics([]bluetooth.UUID{
		ApiV2UUID,
	})
	if err != nil || len(chars) == 0 {
		return errors.New("could not find API characteristic")
	}

	r.apiCharacteristic = &chars[0]

	chars, err = r.spheroService.DiscoverCharacteristics([]bluetooth.UUID{
		AntiDOSUUID,
	})
	if err != nil || len(chars) == 0 {
		return errors.New("could not find Anti DOS characteristic")
	}

	r.antiDOSCharacteristic = &chars[0]

	r.AntiDOS()
	r.Wake()

	return
}

func (r *Robot) Halt() (err error) {
	return
}

// AntiDOS sends message that prevents Sphero from disconnecting.
func (r *Robot) AntiDOS() error {
	_, err := r.antiDOSCharacteristic.WriteWithoutResponse([]byte("usetheforce...band"))
	return err
}

// Wake brings the device out of sleep mode
func (r *Robot) Wake() error {
	_, err := r.send(r.apiCharacteristic, DevicePowerInfo, PowerCommandsWake, true, []byte{})

	return err
}

// Sleep puts the device into sleep mode to save battery
func (r *Robot) Sleep() error {
	_, err := r.send(r.apiCharacteristic, DevicePowerInfo, PowerCommandsSleep, true, []byte{})

	return err
}

// SetLEDColor sets the Sphero LED to the given color.
func (r *Robot) SetLEDColor(c color.RGBA) error {
	payload := []byte{0x00, 0x0e, c.R, c.G, c.B}

	_, err := r.send(r.apiCharacteristic, DeviceUserIO, UserIOCommandsAllLEDs, true, payload)
	return err
}

// Roll towards heading given in degrees 0-360 at speed as an integer 0-255
func (r *Robot) Roll(heading, speed int) error {
	speedH := uint8((speed & 0xFF00) >> 8)
	speedL := uint8(speed & 0xFF)
	headingH := uint8((heading & 0xFF00) >> 8)
	headingL := uint8(heading & 0xFF)

	payload := []byte{speedL, headingH, headingL, speedH}

	_, err := r.send(r.apiCharacteristic, DeviceDriving, DrivingCommandsWithHeading, true, payload)
	return err
}

// Stop moving
func (r *Robot) Stop() (err error) {
	return r.Roll(0, 0)
}

// Enable API notifications
func (r *Robot) EnableNotifications(f func(message *Payload)) error {
	return r.apiCharacteristic.EnableNotifications(func(data []byte) {
		r.buf = append(r.buf, data...)

		for {
			// Find start marker
			startIdx := -1
			for i, b := range r.buf {
				if b == DataPacketStart {
					startIdx = i
					break
				}
			}
			if startIdx == -1 {
				// No start marker, clear buffer
				r.buf = nil
				break
			}

			// Find end marker after start
			endIdx := -1
			for i := startIdx; i < len(r.buf); i++ {
				if r.buf[i] == DataPacketEnd {
					endIdx = i
					break
				}
			}
			if endIdx == -1 {
				// No end marker yet, wait for more data
				if startIdx > 0 {
					// Remove data before start marker
					r.buf = r.buf[startIdx:]
				}
				break
			}

			// Extract packet
			packet := r.buf[startIdx : endIdx+1]
			r.buf = r.buf[endIdx+1:]

			// Decode and handle
			p := &Payload{}
			err := p.Decode(packet)
			if err != nil {
				if debug {
					println("notification error:", err.Error())
				}
				continue
			}

			f(p)
		}
	})
}

// ConfigureCollisionDetection configures collision detection events.
func (r *Robot) ConfigureCollisionDetection(cc CollisionConfig) error {
	payload := []byte{cc.Method, cc.Xt, cc.Xs, cc.Yt, cc.Ys, cc.Dead}

	_, err := r.send(r.apiCharacteristic, DeviceSensor, SensorCommandConfigureCollision, true, payload)
	return err
}

// GetBatteryVoltage requests a battery voltage notification.
func (r *Robot) GetBatteryVoltage() error {
	_, err := r.send(r.apiCharacteristic, DevicePowerInfo, PowerCommandsBatteryVoltage, true, []byte{})

	return err
}

// https://github.com/MProx/Sphero_mini/blob/1dea6ff7f59260ea5ecee9cb9a7c9f46f1f8a6d9/sphero_mini.py#L243
func (r *Robot) send(dc *bluetooth.DeviceCharacteristic, deviceID, commandID byte, expectResponse bool, message []byte) (*Payload, error) {
	// sequence ensures we can associate a request with a response
	r.sequenceNo += 1
	if r.sequenceNo > 255 {
		r.sequenceNo = 0
	}

	// are we expecting a response
	if expectResponse {
		r.expectedCommandSequence = r.sequenceNo
	}

	// define the header for the send request
	p := Payload{
		Flags:    FlagResetsInactivityTimeout + FlagRequestsResponse, // set the flags
		DeviceID: deviceID,                                           // send is for the given device id
		Command:  commandID,                                          // with the command
		Sequence: byte(r.sequenceNo),                                 // set the sequence id to ensure that packets are orderable
		Payload:  message,
	}

	data := p.Encode()

	if debug {
		println("sending data", "bytes", data)
	}

	_, err := dc.WriteWithoutResponse(data)
	if err != nil {
		return nil, err
	}

	if !expectResponse {
		return nil, nil
	}

	// TODO: wait for real response?
	time.Sleep(100 * time.Millisecond)

	return nil, nil
}
