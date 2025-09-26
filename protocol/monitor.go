package protocol

import (
	"github.com/devicehub-go/thorlabs-kdc101/internal/utils"
)

type DCStatusUpdate struct {
	Channel    uint16
	Position   int32
	Velocity   uint16
	Current    int16
	StatusBits uint32
}

type DCStatusBits struct {
	CWHardLimit      bool
	CCWHardLimit     bool
	CWSoftLimit      bool
	CCWSoftLimit     bool
	InMotionCW       bool
	InMotionCCW      bool
	JoggingCW        bool
	JoggingCCW       bool
	IsConnected      bool
	IsHoming 	     bool
	IsHomed          bool
	IsInitializing   bool
	IsTracking       bool
	IsSettled        bool
	PositionError    bool
	InstructionError bool
	Interlock        bool
	OverTemperature  bool
	BusVoltageFault  bool
	CommutationError bool
	Overload         bool
	EncoderFault     bool
	OverCurrent      bool
	BusCurrentFault  bool
	PowerOk          bool
	IsActive         bool
	Error            bool
	IsEnabled        bool
}

type DCStatusUpdateSI struct {
	Channel    uint16
	Position   float64
	Velocity   float64
	Current    float64
	StatusBits DCStatusBits
}

/*
Request a status update for the specified DC motor channel
*/
func (k *KDC101) GetDCStatusUpdate(channel uint8) (DCStatusUpdate, error) {
	if channel != 1 {
		return DCStatusUpdate{}, ErrChannelNotSupported
	}
	msg := HeaderMessage{
		ID:          0x0490,
		Parameter1:  byte(1 << (channel - 1)),
		Parameter2:  0x00,
		Destination: GenericUnit,
		Source:      Host,
	}
	response, err := k.RequestData(msg)
	if err != nil {
		return DCStatusUpdate{}, err
	}
	data := response.Data
	return DCStatusUpdate{
		Channel:    utils.BytesToWord(data[0:2]),
		Position:   utils.BytesToLong(data[2:6]),
		Velocity:   utils.BytesToWord(data[6:8]),
		Current:    utils.BytesToShort(data[8:10]),
		StatusBits: utils.BytesToDword(data[10:14]),
	}, nil
}

/*
Re-scales the DC status update data according to the motor
and stage type
*/
func (k *KDC101) DCStatusUpdateToSI(update DCStatusUpdate) DCStatusUpdateSI {
	return DCStatusUpdateSI{
		Channel:  update.Channel,
		Position: k.CountsToPosition(update.Position),
		Velocity: k.CountsToVelocity(uint32(update.Velocity)),
		Current:  float64(update.Current),
		StatusBits: k.ParseDCStatusBits(update.StatusBits),
	}
}

/*
Parses the status bits from the DC status update
*/
func (k *KDC101) ParseDCStatusBits(statusBits uint32) DCStatusBits {
	return DCStatusBits{
		CWHardLimit:      (statusBits & 0x00000001) != 0,
		CCWHardLimit:     (statusBits & 0x00000002) != 0,
		CWSoftLimit:      (statusBits & 0x00000004) != 0,
		CCWSoftLimit:     (statusBits & 0x00000008) != 0,
		InMotionCW:       (statusBits & 0x00000010) != 0,
		InMotionCCW:      (statusBits & 0x00000020) != 0,
		JoggingCW:        (statusBits & 0x00000040) != 0,
		JoggingCCW:       (statusBits & 0x00000080) != 0,
		IsConnected:      (statusBits & 0x00000100) != 0,
		IsHoming:         (statusBits & 0x00000200) != 0,
		IsHomed:          (statusBits & 0x00000400) != 0,
		IsInitializing:   (statusBits & 0x00000800) != 0,
		IsTracking:       (statusBits & 0x00001000) != 0,
		IsSettled:        (statusBits & 0x00002000) != 0,
		PositionError:    (statusBits & 0x00004000) != 0,
		InstructionError: (statusBits & 0x00008000) != 0,
		Interlock:        (statusBits & 0x00010000) != 0,
		OverTemperature:  (statusBits & 0x00020000) != 0,
		BusVoltageFault:  (statusBits & 0x00040000) != 0,
		CommutationError: (statusBits & 0x00080000) != 0,
		Overload:         (statusBits & 0x01000000) != 0,
		EncoderFault:     (statusBits & 0x02000000) != 0,
		OverCurrent:      (statusBits & 0x04000000) != 0,
		BusCurrentFault:  (statusBits & 0x08000000) != 0,
		PowerOk:          (statusBits & 0x10000000) != 0,
		IsActive:         (statusBits & 0x20000000) != 0,
		Error:            (statusBits & 0x40000000) != 0,
		IsEnabled:        (statusBits & 0x80000000) != 0,
	}
}
