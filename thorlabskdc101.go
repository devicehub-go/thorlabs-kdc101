package main

import (
	"github.com/devicehub-go/thorlabs-kdc101/protocol"
	"github.com/devicehub-go/unicomm"
)

type KDC101 = protocol.KDC101
type MotorType string
type StageType string

const (
	Brushed   MotorType = "Brushed"
	Brushless MotorType = "Brushless"

	MTS25Z8 StageType = "MTS25-Z8"
	MTS50Z8 StageType = "MTS50-Z8"
	Z8xx    StageType = "Z8xx"
	Z6xx    StageType = "Z6xx"
	PRM1Z8  StageType = "PRM1-Z8"
	PRMTZ8  StageType = "PRMTZ8"
	CR1Z7   StageType = "CR1-Z7"
	KVS30   StageType = "KVS30"
)

/*
Creates a new instance of Brushed Motor Controller KDC101 that
allows to communicate and control the connected motors
*/
func New(stage StageType, motor MotorType, options unicomm.UnicommOptions) *KDC101 {
	oem750 := &KDC101{
		Communication: unicomm.New(options),
		StageType: string(stage),
		MotorType: string(motor),
	}
	return oem750
}
