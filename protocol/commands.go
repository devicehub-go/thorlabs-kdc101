/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package protocol

import (
	"fmt"

	"github.com/devicehub-go/thorlabs-kdc101/internal/utils"
)

type Direction uint8
type StopMode  uint8

type HwInformation struct {
	SerialNumber    int32
	Model           string
	Type            uint16
	FirmwareVersion []byte
	HardwareVersion uint16
	ModState        uint16
	NumberChannels  uint16
}

const (
	Forward Direction = 0x01
	Reverse Direction = 0x02

	Abrupt StopMode = 0x01
	Soft   StopMode = 0x02
)

/*
Instruct hardware unit to identify itself by flashing
its front panel LEDs.
*/
func (k *KDC101) Identify(channel uint8) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0223,
		Parameter1:  byte(1 << (channel - 1)),
		Parameter2:  0x00,
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Request hardware information from the controller
*/
func (k *KDC101) GetInformation() (HwInformation, error) {
	response, err := k.RequestData(HeaderMessage{
		ID:          0x0005,
		Parameter1:  0x00,
		Parameter2:  0x00,
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return HwInformation{}, err
	}
	data := response.Data
	if len(data) < 84 {
		return HwInformation{}, fmt.Errorf("invalid response length")
	}
	return HwInformation{
		SerialNumber:    utils.BytesToLong(data[0:4]),
		Model:           string(data[4:12]),
		Type:            utils.BytesToWord(data[12:14]),
		FirmwareVersion: data[14:18],
		HardwareVersion: utils.BytesToWord(data[78:80]),
		ModState:        utils.BytesToWord(data[80:82]),
		NumberChannels:  utils.BytesToWord(data[82:84]),
	}, nil
}

/*
Start a home move sequence on the specified channel
in accordance with the home paramters set
*/
func (k *KDC101) StartHomeMove(channel uint8) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0443,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Starts a relative move on the specified channel
in accordance with the distance parameters set
*/
func (k *KDC101) StartRelativeMove(channel uint8) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0448,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Starts a relative move on the specified channel
with the target distance
*/
func (k *KDC101) MoveRelativeDistance(channel uint8, distance float64) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	var counts int32 = k.PositionToCounts(distance)
	data := []byte{
		byte(1 << (channel - 1)),
		0x00,
	}
	data = append(data, utils.LongToBytes(counts)...)
	return k.WriteData(DataMessage{
		ID:          0x0448,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Starts an absolute mvoe on the specified channel
in accordance with the absolute move parameters set
*/
func (k *KDC101) StartAbsoluteMove(channel uint8) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0453,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Starts an absolute move on the specified channel
with the target position
*/
func (k *KDC101) MoveAbsolutePosition(channel uint8, position float64) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	var counts int32 = k.PositionToCounts(position)
	data := []byte{
		byte(1 << (channel - 1)),
		0x00,
	}
	data = append(data, utils.LongToBytes(counts)...)
	return k.WriteData(DataMessage{
		ID:          0x0453,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Start a jog move on the specified motor channel
*/
func (k *KDC101) StartJogMove(channel uint8, direction Direction) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x046A,
		Parameter1:  byte(1 << (channel - 1)),
		Parameter2:  byte(direction),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Moves the motor continuously in the specified direction
using the velocity parameters set until a stop command
or limit is reached.
*/
func (k *KDC101) MoveContinuous(channel uint8, direction Direction) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0457,
		Parameter1:  byte(1 << (channel - 1)),
		Parameter2:  byte(direction),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Stops the motor on the specified channel
*/
func (k *KDC101) Stop(channel uint8, mode StopMode) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	return k.WriteHeaderOnly(HeaderMessage{
		ID:          0x0465,
		Parameter1:  byte(1 << (channel - 1)),
		Parameter2:  byte(mode),
		Destination: GenericUnit,
		Source:      Host,
	})
}
