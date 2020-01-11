// Package conversion TODO(kkweon): Write description
package conversion

import (
	"fmt"

	"encoding/binary"
	"errors"
	"math"
)

type Endian int

const BigEndian Endian = 0
const LittleEndian Endian = 1

type Option struct {
	Endian Endian
}

// IntToFloat32 converts int to float32
func IntToFloat32(x int) (float32, error) {
	return float32(x), nil
}

// Flattens [n][4]byte to [n * 4]byte
func flattenBytes32(xs [][]byte) ([]byte, error) {
	return flattenInternal(xs, 4)
}

// Flattens [n][8]byte to [n * 8]byte
func flattenBytes64(xs [][]byte) ([]byte, error) {
	return flattenInternal(xs, 8)
}

func flattenInternal(xs [][]byte, size int) ([]byte, error) {
	ret := make([]byte, len(xs)*size)
	for i, bytes := range xs {
		if len(bytes) != size {
			return nil, fmt.Errorf("expected [][%v]byte, but received [][%v]byte", size, len(bytes))
		}
		for j, b := range bytes {
			ret[i*size+j] = b
		}
	}
	return ret, nil
}

func unflattenBytes32(xs []byte) ([][]byte, error) {
	return unflattenInternal(xs, 4)
}
func unflattenBytes64(xs []byte) ([][]byte, error) {
	return unflattenInternal(xs, 8)
}

// unflattenInternal transforms []byte to [][size]byte
// If the original data was []float32 then it can be transformed into [][4]byte.
// Then, [][4]byte can get flatten to []byte.
// This function "undo" the flatten by making []byte back to [][4]byte.
func unflattenInternal(xs []byte, size int) ([][]byte, error) {
	if len(xs)%size != 0 {
		return nil, fmt.Errorf("expected the length of xs to be a multiple of %v but its size is %v", size, len(xs))
	}

	var ret [][]byte = nil

	for i := 0; i < len(xs); i += size {
		ret = append(ret, xs[i:i+size])
	}

	return ret, nil
}

// Float32ToBytes converts float32 to []byte with length 4
func Float32ToBytes(x float32, o *Option) []byte {
	b := make([]byte, 4)
	ux := math.Float32bits(x)
	if o == nil || o.Endian == BigEndian {
		binary.BigEndian.PutUint32(b, ux)
	} else {
		binary.LittleEndian.PutUint32(b, ux)
	}
	return b
}

// BytesToFloat32 converts []byte to float32
// length of []byte should be 4 or bigger
func BytesToFloat32(x []byte, o *Option) (float32, error) {
	var fx float32

	if len(x) < 4 {
		return fx, errors.New("length of []byte should be 4 or bigger")
	}

	var ux uint32
	if o == nil || o.Endian == BigEndian {
		ux = binary.BigEndian.Uint32(x)
	} else {
		ux = binary.LittleEndian.Uint32(x)
	}
	return math.Float32frombits(ux), nil
}

// Float64ToBytes converts float64 to []byte with length 8
func Float64ToBytes(x float64, o *Option) []byte {
	b := make([]byte, 8)
	ux := math.Float64bits(x)
	if o == nil || o.Endian == BigEndian {
		binary.BigEndian.PutUint64(b, ux)
	} else {
		binary.LittleEndian.PutUint64(b, ux)
	}
	return b
}

// BytesToFloat64 converts []byte to float64
// length of []byte should be 8 or bigger
func BytesToFloat64(x []byte, o *Option) (float64, error) {
	var fx float64

	if len(x) < 8 {
		return fx, errors.New("length of []byte should be 4 or bigger")
	}

	var ux uint64
	if o == nil || o.Endian == BigEndian {
		ux = binary.BigEndian.Uint64(x)
	} else {
		ux = binary.LittleEndian.Uint64(x)
	}
	return math.Float64frombits(ux), nil
}
