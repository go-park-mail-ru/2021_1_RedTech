package cast

import (
	"encoding/binary"
	"math"
)

func ToString(bytes []byte) string {
	return string(bytes)
}

func ToFloat(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func ToUint(bytes []byte) uint {
	return uint(binary.BigEndian.Uint32(bytes))
}

func ToInt(bytes []byte) int {
	return int(binary.BigEndian.Uint32(bytes))
}

func ToSmallInt(bytes []byte) int {
	return int(binary.BigEndian.Uint16(bytes))
}
