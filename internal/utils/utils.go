/*
Author: Leonardo Rossi Leao
Created at: September 26th, 2025
Last update: September 26th, 2025
*/

package utils

import "encoding/binary"

/*
Converts a little-endian four byte array to a
32-bit signed integer (long)
*/
func BytesToLong(data []byte) int32 {
	return int32(binary.LittleEndian.Uint32(data))
}

/*
Converts a long integer to a little-endian four byte array
*/
func LongToBytes(value int32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(value))
	return data
}

/*
Converts a little-endian four byte array to a
32-bit unsigned integer (dword)
*/
func BytesToDword(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

/*
Converts a dword integer to a little-endian four byte array
*/
func DwordToBytes(value uint32) []byte {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, value)
	return data
}

/*
Converts a little-endian two byte array to a 16-bit
unsigned integer (word)
*/
func BytesToWord(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

/*
Converts a word integer to a little-endian two byte array
*/
func WordToBytes(value uint16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, value)
	return data
}

/*
Converts a little-endian two byte array to a 16-bit
signed integer (short)
*/
func BytesToShort(data []byte) int16 {
	return int16(binary.LittleEndian.Uint16(data))
}

/*
Converts a short integer to a little-endian two byte array
*/
func ShortToBytes(value int16) []byte {
	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(value))
	return data
}