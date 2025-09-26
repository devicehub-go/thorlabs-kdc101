/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package protocol

import (
	"fmt"
	"time"

	"github.com/devicehub-go/unicomm"
)

type Endpoint byte

type HeaderMessage struct {
	ID          uint16
	Parameter1  byte
	Parameter2  byte
	Destination Endpoint
	Source      Endpoint
}

type DataMessage struct {
	ID          uint16
	DataLength  uint16
	Destination Endpoint
	Source      Endpoint
	Data        []byte
}

type KDC101 struct {
	Communication unicomm.Unicomm
	StageType string // e.g., "MTS25-Z8", "MTS50-Z8", etc.
	MotorType string // e.g., "Brushed", "Brushless"
}

const (
	Host        Endpoint = 0x01
	Rack        Endpoint = 0x02
	GenericUnit Endpoint = 0x50
)

var ErrChannelNotSupported = fmt.Errorf("KDC101 just supports channel 1")
var ErrInvalidResponseLength = fmt.Errorf("invalid response length")
var InvalidHeader HeaderMessage = HeaderMessage{}
var InvalidData   DataMessage = DataMessage{}

/*
Establishes a connection with the device
*/
func (k *KDC101) Connect() error {
	if err := k.Communication.Connect(); err != nil {
		return err
	}
	return nil
}

/*
Closes the connection with the device
*/
func (k *KDC101) Disconnect() error {
	return k.Communication.Disconnect()
}

/*
Returns true if device is connected
*/
func (k *KDC101) IsConnected() bool {
	return k.Communication.IsConnected()
}

/*
Writes a header only message
*/
func (k *KDC101) WriteHeaderOnly(msg HeaderMessage) error {
	bytes := []byte{
		byte(msg.ID & 0x00FF),
		byte(msg.ID >> 8),
		msg.Parameter1,
		msg.Parameter2,
		byte(msg.Destination),
		byte(msg.Source),
	}
	return k.Communication.Write(bytes)
}

/*
Writes a data message
*/
func (k *KDC101) WriteData(msg DataMessage) error {
	bytes := []byte{
		byte(msg.ID & 0x00FF),
		byte(msg.ID >> 8),
		byte(msg.DataLength & 0x00FF),
		byte(msg.DataLength >> 8),
		byte(msg.Destination) | 0x80,
		byte(msg.Source),
	}
	bytes = append(bytes, msg.Data...)
	return k.Communication.Write(bytes)
}

/*
Reads a header only response
*/
func (k *KDC101) ReadHeaderOnly() (HeaderMessage, error) {
	response, err := k.Communication.Read(6)
	if err != nil {
		return InvalidHeader, err
	}
	msg := HeaderMessage{
		ID:          uint16(response[1])<<8 | uint16(response[0]),
		Parameter1:  response[2],
		Parameter2:  response[3],
		Destination: Endpoint(response[4]),
		Source:      Endpoint(response[5]),
	}
	return msg, nil
}

/*
Reads a message which contains header and data
*/
func (k *KDC101) ReadData() (DataMessage, error) {
	response, err := k.Communication.Read(6)
	if err != nil {
		return InvalidData, err
	}
	msg := DataMessage{
		ID:          uint16(response[1])<<8 | uint16(response[0]),
		DataLength:  uint16(response[3])<<8 | uint16(response[2]),
		Destination: Endpoint(response[4]),
		Source:      Endpoint(response[5]),
	}
	if msg.DataLength < 1 {
		return InvalidData, fmt.Errorf("invalid data length: %d", msg.DataLength)
	}
	data, err := k.Communication.Read(uint(msg.DataLength))
	if err != nil {
		return InvalidData, err
	}
	msg.Data = data
	return msg, nil
}

/*
Sends a header only message to device and waits for a 
header only response.
*/
func (k *KDC101) RequestHeaderOnly(msg HeaderMessage) (HeaderMessage, error) {
	err := k.WriteHeaderOnly(msg)
	if err != nil {
		return InvalidHeader, err
	}
	time.Sleep(15 * time.Millisecond)
	return k.ReadHeaderOnly()
}

/*
Sends a header only message to device and waits for a
data message response.
*/
func (k *KDC101) RequestData(msg HeaderMessage) (DataMessage, error) {
	err := k.WriteHeaderOnly(msg)
	if err != nil {
		return InvalidData, err
	}
	time.Sleep(50 * time.Millisecond)
	return k.ReadData()
}
