package cast

import (
	"bytes"
	"encoding/binary"
	"math"
)

func ToString(bytesArg []byte) string {
	return string(bytesArg)
}

func ToFloat(bytesArg []byte) float32 {
	bits := binary.BigEndian.Uint32(bytesArg)
	return math.Float32frombits(bits)
}

func ToUint(bytesArg []byte) uint {
	return uint(binary.BigEndian.Uint32(bytesArg))
}

func ToInt(bytesArg []byte) int {
	reader := bytes.NewReader(bytesArg)
	num, _ := binary.ReadVarint(reader)
	return int(num)
}

func ToSmallInt(bytesArg []byte) int {
	reader := bytes.NewReader(bytesArg)
	num, _ := binary.ReadVarint(reader)
	return int(num)
}

func StrToBytes(arg string) []byte {
	return []byte(arg)
}

func FloatToBytes(arg float32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, math.Float32bits(arg))
	return bytes
}

func UintToBytes(arg uint) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(arg))
	return bytes
}

func IntToBytes(arg int) []byte {
	bytes := make([]byte, 4)
	binary.PutVarint(bytes, int64(arg))
	return bytes
}

func SmallIntToBytes(arg int) []byte {
	bytes := make([]byte, 8)
	binary.PutVarint(bytes, int64(arg))
	return bytes
}
