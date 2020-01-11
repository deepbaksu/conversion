// Package conversion TODO(kkweon): Write description
package conversion

import (
	"fmt"
)

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
