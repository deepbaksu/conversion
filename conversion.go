// Package conversion TODO(kkweon): Write description
package conversion

import (
	"github.com/deepbaksu/conversion/internal"

	"encoding/binary"
	"errors"
	"math"
)

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
)

type Option struct {
	Endian Endian
}

// IntToFloat32 converts int to float32
func IntToFloat32(x int) (float32, error) {
	return float32(x), nil
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
		return fx, errors.New("length of []byte should be 8 or bigger")
	}

	var ux uint64
	if o == nil || o.Endian == BigEndian {
		ux = binary.BigEndian.Uint64(x)
	} else {
		ux = binary.LittleEndian.Uint64(x)
	}
	return math.Float64frombits(ux), nil
}

// Float32sToBytes converts []float32 to []byte.
// Step 1. Make 2D byte slice by []float32 -> [][4]byte
// Step 2. Flatten [][4]byte to []byte
func Float32sToBytes(xs []float32, o *Option) ([]byte, error) {
	bs := make([][]byte, len(xs))
	for i, x := range xs {
		bs[i] = Float32ToBytes(x, o)
	}
	return internal.FlattenBytes32(bs)
}

// BytesToFloat32s converts []byte to []float32
// Step 1. Unflatten []byte to [][4]byte
// Step 2. [][4]byte to []float32
func BytesToFloat32s(xs []byte, o *Option) ([]float32, error) {
	var fs []float32
	xxs, err := internal.UnflattenBytes32(xs)
	if err != nil {
		return nil, err
	}

	for _, xs := range xxs {
		xs, err := BytesToFloat32(xs, o)
		if err != nil {
			return nil, err
		}
		fs = append(fs, xs)
	}

	return fs, nil
}

// Float64sToBytes converts []float64 to []byte.
// Step 1. Make 2D byte slice by []float64 -> [][8]byte
// Step 2. Flatten [][8]byte to []byte
func Float64sToBytes(xs []float64, o *Option) ([]byte, error) {
	bs := make([][]byte, len(xs))
	for i, x := range xs {
		bs[i] = Float64ToBytes(x, o)
	}
	return internal.FlattenBytes64(bs)
}

// BytesToFloat64s converts []byte to []float64
// Step 1. Unflatten []byte to [][8]byte
// Step 2. [][8]byte to []float64
func BytesToFloat64s(xs []byte, o *Option) ([]float64, error) {
	var fs []float64
	xxs, err := internal.UnflattenBytes64(xs)
	if err != nil {
		return nil, err
	}

	for _, xs := range xxs {
		xs, err := BytesToFloat64(xs, o)
		if err != nil {
			return nil, err
		}
		fs = append(fs, xs)
	}

	return fs, nil
}
