/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package protocol

var MotorTFactor = map[string]float64{
	"Brushed":   2048.0 / (6.0 * 1e6),
	"Brushless": 2048.0 / (6.0 * 1e6),
}

var StageScalingFactor = map[string]float64{
	"MTS25-Z8": 34554.96,
	"MTS50-Z8": 34554.96,
	"Z8xx":     34554.96,
	"Z6xx":     24600.0,
	"PRM1-Z8":  1919.6418578623391,
	"PRMTZ8":   1919.6418578623391,
	"CR1-Z7":   12288.0,
	"KVS30":    20000.0,
}

/*
Converts position in millimeters to encoder counts
*/
func (k *KDC101) PositionToCounts(position float64) int32 {
	encCount := StageScalingFactor[k.StageType]
	return int32(position * encCount)
}

/*
Converts encoder counts to position in millimeters
*/
func (k *KDC101) CountsToPosition(counts int32) float64 {
	encCount := StageScalingFactor[k.StageType]
	return float64(counts) / encCount
}

/*
Converts velocity in millimeters per second to encoder
counts per second
*/
func (k *KDC101) VelocityToCounts(velocity float64) uint32 {
	encCount := StageScalingFactor[k.StageType]
	T := MotorTFactor[k.MotorType]
	return uint32(velocity * T * 65536 * encCount)
}

/*
Converts encoder counts per second to velocity in millimeters
*/
func (k *KDC101) CountsToVelocity(counts uint32) float64 {
	encCount := StageScalingFactor[k.StageType]
	T := MotorTFactor[k.MotorType]
	return float64(counts) / (T * 65536 * encCount)
}

/*
Converts acceleration in millimeters per second squared to
encoder counts per second squared
*/
func (k *KDC101) AccelerationToCounts(acceleration float64) uint32 {
	encCount := StageScalingFactor[k.StageType]
	T := MotorTFactor[k.MotorType]
	return uint32(acceleration * (T * T) * 65536 * encCount)
}

/*
Converts encoder counts per second squared to acceleration
in millimeters per second squared
*/
func (k *KDC101) CountsToAcceleration(counts int32) float64 {
	encCount := StageScalingFactor[k.StageType]
	T := MotorTFactor[k.MotorType]
	return float64(counts) / (T * T * 65536 * encCount)
}