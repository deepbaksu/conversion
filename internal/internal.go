package internal

import "fmt"

// Flattens [n][4]byte to [n * 4]byte
func FlattenBytes32(xs [][]byte) ([]byte, error) {
	return FlattenInternal(xs, 4)
}

// Flattens [n][8]byte to [n * 8]byte
func FlattenBytes64(xs [][]byte) ([]byte, error) {
	return FlattenInternal(xs, 8)
}

func FlattenInternal(xs [][]byte, size int) ([]byte, error) {
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

func UnflattenBytes32(xs []byte) ([][]byte, error) {
	return UnflattenInternal(xs, 4)
}
func UnflattenBytes64(xs []byte) ([][]byte, error) {
	return UnflattenInternal(xs, 8)
}

// unflattenInternal transforms []byte to [][size]byte
// If the original data was []float32 then it can be transformed into [][4]byte.
// Then, [][4]byte can get flatten to []byte.
// This function "undo" the flatten by making []byte back to [][4]byte.
func UnflattenInternal(xs []byte, size int) ([][]byte, error) {
	if len(xs)%size != 0 {
		return nil, fmt.Errorf("expected the length of xs to be a multiple of %v but its size is %v", size, len(xs))
	}

	var ret [][]byte = nil

	for i := 0; i < len(xs); i += size {
		ret = append(ret, xs[i:i+size])
	}

	return ret, nil
}

