package cast

import (
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

func ToUint64(bytesArg []byte) uint {
	return uint(binary.BigEndian.Uint64(bytesArg))
}

func ToInt(bytesArg []byte) int {
	return int(int32(binary.BigEndian.Uint32(bytesArg)))
}

func ToSmallInt(bytesArg []byte) int {
	return int(int16(binary.BigEndian.Uint16(bytesArg)))
}

func ToBool(bytesArg []byte) bool {
	return bytesArg[0] > 0
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

func Uint64ToBytes(arg uint64) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, arg)
	return bytes
}

func IntToBytes(arg int) []byte {
	bytes := make([]byte, 4)
	binary.PutVarint(bytes, int64(arg))
	return bytes
}

func SmallIntToBytes(arg int) []byte {
	bytes := make([]byte, 10)
	binary.PutUvarint(bytes, uint64(arg))
	return bytes
}

func BoolToBytes(arg bool) []byte {
	if arg {
		return []byte{1}
	}
	return []byte{0}
}
