/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package protocol

import (
	"github.com/devicehub-go/thorlabs-kdc101/internal/utils"
)

type VelocityProfile struct {
	MinVelocity float64
	MaxVelocity float64
	Acceleration float64
}

type JogParameters struct {
	Mode         uint16
	StepSize     float64
	MinVelocity  float64
	Acceleration float64
	MaxVelocity  float64
	StopMode     uint16
}

/*
Sent to enable or disable the specified drive channel.
*/
func (k *KDC101) Enable(channel uint8, enable bool) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	msg := HeaderMessage{
		ID:          0x0224,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	}
	if enable {
		msg.Parameter2 = 0x01
	} else {
		msg.Parameter2 = 0x02
	}
	return k.WriteHeaderOnly(msg)
}

/*
Get the enabled state of a channel
*/
func (k *KDC101) IsEnabled(channel uint8) (bool, error) {
	if channel != 1 {
		return false, ErrChannelNotSupported
	}
	response, err := k.RequestHeaderOnly(HeaderMessage{
		ID:          0x0211,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return false, err
	}
	return response.Parameter2 == 0x01, nil
}

/*
Sets trapezoidal velocity parameters for the specified
motor channel.
*/
func (k *KDC101) SetTrapezoidalVelocity(channel uint8, profile VelocityProfile) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	minVel := k.VelocityToCounts(profile.MinVelocity)
	accel := k.AccelerationToCounts(profile.Acceleration)
	maxVel := k.VelocityToCounts(profile.MaxVelocity)

	data := []byte{
		byte(1 << (channel - 1)),
		0x00,
	}
	data = append(data, utils.DwordToBytes(minVel)...)
	data = append(data, utils.DwordToBytes(accel)...)
	data = append(data, utils.DwordToBytes(maxVel)...)

	return k.WriteData(DataMessage{
		ID:          0x0413,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Get trapezoidal velocity parameters for the specified
motor channel
*/
func (k *KDC101) GetTrapezoidalVelocity(channel uint8) (VelocityProfile, error) {
	if channel != 1 {
		return VelocityProfile{}, ErrChannelNotSupported
	}
	response, err := k.RequestData(HeaderMessage{
		ID:          0x0414,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return VelocityProfile{}, err
	}

	data := response.Data
	if len(data) < 14 {
		return VelocityProfile{}, ErrInvalidResponseLength
	}

	minVel := k.CountsToVelocity(utils.BytesToDword(data[2:6]))
	accel := k.CountsToAcceleration(utils.BytesToLong(data[6:10]))
	maxVel := k.CountsToVelocity(utils.BytesToDword(data[10:14]))

	return VelocityProfile{
		MinVelocity: minVel,
		Acceleration: accel,
		MaxVelocity: maxVel,
	}, nil
}

/*
Set the velocity jog paramaters for the specified channel.
*/
func (k *KDC101) SetJogParameters(channel uint8, params JogParameters) error {
	if channel != 1 {
		return ErrChannelNotSupported
	}
	stepSize := k.PositionToCounts(params.StepSize)
	minVel := k.VelocityToCounts(params.MinVelocity)
	accel := k.AccelerationToCounts(params.Acceleration)
	maxVel := k.VelocityToCounts(params.MaxVelocity)

	data := []byte{
		byte(1 << (channel - 1)),
		0x00,
	}
	data = append(data, utils.WordToBytes(params.Mode)...)
	data = append(data, utils.LongToBytes(stepSize)...)
	data = append(data, utils.DwordToBytes(minVel)...)
	data = append(data, utils.DwordToBytes(accel)...)
	data = append(data, utils.DwordToBytes(maxVel)...)
	data = append(data, utils.WordToBytes(params.StopMode)...)

	return k.WriteData(DataMessage{
		ID:          0x0416,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Get the jog parameters for the specified channel.
*/
func (k *KDC101) GetJogParameters(channel uint8) (JogParameters, error) {
	if channel != 1 {
		return JogParameters{}, ErrChannelNotSupported
	}
	response, err := k.RequestData(HeaderMessage{
		ID:          0x0417,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return JogParameters{}, err
	}

	data := response.Data
	if len(data) < 22 {
		return JogParameters{}, ErrInvalidResponseLength
	}

	mode := utils.BytesToWord(data[2:4])
	stepSize := k.CountsToPosition(utils.BytesToLong(data[4:8]))
	minVel := k.CountsToVelocity(utils.BytesToDword(data[8:12]))
	accel := k.CountsToAcceleration(utils.BytesToLong(data[12:16]))
	maxVel := k.CountsToVelocity(utils.BytesToDword(data[16:20]))
	stopMode := utils.BytesToWord(data[20:22])

	return JogParameters{
		Mode:         mode,
		StepSize:     stepSize,
		MinVelocity:  minVel,
		Acceleration: accel,
		MaxVelocity:  maxVel,
		StopMode:     stopMode,
	}, nil
}

/*
Sets the relative move distance that will be used
the next time that a relative move is initiated
*/
func (k *KDC101) SetRelativeMoveDistance(channel uint8, distance float64) error {
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
		ID:          0x0445,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Gets the target distance for the next relative move
*/
func (k *KDC101) GetRelativeMoveDistance(channel uint8) (float64, error) {
	if channel != 1 {
		return 0, ErrChannelNotSupported
	}
	response, err := k.RequestData(HeaderMessage{
		ID:          0x0446,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return 0, err
	}
	data := response.Data
	if len(data) < 6 {
		return 0, ErrInvalidResponseLength
	}
	counts := utils.BytesToLong(data[2:6])
	return k.CountsToPosition(int32(counts)), nil
}

/*
Sets the absolute move parameters for the specified channel.
The target position is used the next time that an absolute
move is initiated.
*/
func (k *KDC101) SetAbsoluteMoveDistance(channel uint8, position float64) error {
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
		ID:          0x0450,
		Data:        data,
		DataLength:  uint16(len(data)),
		Destination: GenericUnit,
		Source:      Host,
	})
}

/*
Gets the target position for the next absolute move
*/
func (k *KDC101) GetAbsoluteMoveDistance(channel uint8) (float64, error) {
	if channel != 1 {
		return 0, ErrChannelNotSupported
	}
	response, err := k.RequestData(HeaderMessage{
		ID:          0x0451,
		Parameter1:  byte(1 << (channel - 1)),
		Destination: GenericUnit,
		Source:      Host,
	})
	if err != nil {
		return 0, err
	}
	data := response.Data
	if len(data) < 6 {
		return 0, ErrInvalidResponseLength
	}
	counts := utils.BytesToLong(data[2:6])
	return k.CountsToPosition(int32(counts)), nil
}