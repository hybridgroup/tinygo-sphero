package sphero

import (
	"encoding/hex"
	"fmt"
)

type Payload struct {
	Flags    uint8
	DeviceID uint8
	Command  uint8
	Sequence uint8
	Error    uint8
	Payload  []byte
}

func (p *Payload) Encode() []byte {
	sendBytes := []byte{
		DataPacketStart, // first byte is always 0x8d
		p.Flags,         // set the flags
		p.DeviceID,      // send is for the given device id
		p.Command,       // with the command
		p.Sequence,      // set the sequence id to ensure that packets are orderable
	}

	sendBytes = append(sendBytes, p.Payload...)
	cs := calculateChecksum(sendBytes[1:])
	sendBytes = append(sendBytes, cs, DataPacketEnd)

	return sendBytes
}

// Decode expects the entire packet, include start, header, payload, checksum, and end bytes.
func (p *Payload) Decode(d []byte) error {
	p.Flags = d[1]
	p.DeviceID = d[2]
	p.Command = d[3]
	p.Sequence = d[4]
	checksum := d[len(d)-2]
	cc := calculateChecksum(d[1 : len(d)-2])

	if checksum != cc {
		return fmt.Errorf("decode checksum for %s invalid: expected %x, received: %x", hex.EncodeToString(d), cc, checksum)
	}

	if len(d) > 7 {
		p.Payload = d[5 : len(d)-2]
	}

	return nil
}

func (p *Payload) String() string {
	return fmt.Sprintf("Flags: %d, DeviceID: %d, Command: %d, Sequence: %d, Payload: %s",
		p.Flags,
		p.DeviceID,
		p.Command,
		p.Sequence,
		hex.EncodeToString(p.Payload))
}

func calculateChecksum(b []byte) uint8 {
	var calculatedChecksum uint16
	for i := range b {
		calculatedChecksum += uint16(b[i])
	}
	return uint8(^(calculatedChecksum % 256))
}
