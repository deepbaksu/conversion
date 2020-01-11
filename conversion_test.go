package conversion

import (
	"bytes"
	"errors"
	"fmt"

	"testing"

	"log"
)

var optBig = &Option{Endian: BigEndian}
var optLittle = &Option{Endian: LittleEndian}

func ExampleIntToFloat32() {
	var x int = 3
	fmt.Println(IntToFloat32(x))
	// Output: 3 <nil>
}

func TestFlattenInternal(t *testing.T) {
	testCases := []struct {
		name     string
		input    [][]byte
		inputDim int
		expected []byte
		err      error
	}{
		{
			name:     "it should flatten [][4]byte to []byte",
			input:    [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}},
			inputDim: 4,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			err:      nil,
		},
		{
			name:     "if the input is not float32, it should return an error",
			input:    [][]byte{{1, 2, 3, 4, 5}, {5, 6, 7, 8, 9}},
			inputDim: 4,
			expected: nil,
			err:      errors.New("expected [][4]byte, but received [][5]byte"),
		},
		{
			name:     "it should flatten [][8]byte to []byte",
			input:    [][]byte{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15, 16}},
			inputDim: 8,
			expected: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			err:      nil,
		},
	}

	for _, testCase := range testCases {

		// capture the value
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ys, err := flattenInternal(testCase.input, testCase.inputDim)

			errMsg := `
Expected: %v 
Received: %v`

			if !areSameErrors(err, testCase.err) {
				t.Errorf(errMsg, testCase.err, err)
			}
			if !bytes.Equal(ys, testCase.expected) {
				t.Errorf(errMsg, testCase.expected, ys)
			}
		})
	}
}

func TestUnflattenInternal(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		inputDim int
		expected [][]byte
		err      error
	}{
		{
			name:     "it should convert []byte to [][8]byte array",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			inputDim: 8,
			expected: [][]byte{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15, 16}},
		},
		{
			name:     "it should return an error if the input is not a float64 slice",
			input:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			inputDim: 8,
			err:      errors.New("expected the length of xs to be a multiple of 8 but its size is 15"),
		},
	}

	for _, testCase := range testCases {

		// capture the value
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			ys, err := unflattenInternal(testCase.input, testCase.inputDim)

			errMsg := `
Expected: %v 
Received: %v`

			if !areSameErrors(err, testCase.err) {
				t.Errorf(errMsg, testCase.err, err)
			}
			if !areSame2dBytes(ys, testCase.expected) {
				t.Errorf(errMsg, testCase.expected, ys)
			}
		})
	}
}

func areSame2dBytes(ys [][]byte, expected [][]byte) bool {
	if len(ys) != len(expected) {
		return false
	}
	for i := range ys {
		if !bytes.Equal(ys[i], expected[i]) {
			return false
		}
	}

	return true
}

// Returns true if both err and err2 are <nil> or have the same error message.
func areSameErrors(err error, err2 error) bool {
	if err != nil && err2 != nil {
		return err.Error() == err2.Error()
	}
	return err == err2
}

func Example_flattenBytes32() {
	xs := [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}}
	fmt.Print(flattenBytes32(xs))
	// Output: [1 2 3 4 5 6 7 8] <nil>
}

func Example_flattenBytes64() {
	xs := [][]byte{{1, 2, 3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14, 15, 16}}
	fmt.Print(flattenBytes64(xs))
	// Output: [1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16] <nil>
}

func Example_unflattenBytes32() {
	xs := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Print(unflattenBytes32(xs))
	// Output: [[1 2 3 4] [5 6 7 8]] <nil>
}

func Example_unflattenBytes64() {
	xs := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	fmt.Print(unflattenBytes64(xs))
	// Output: [[1 2 3 4 5 6 7 8]] <nil>
}

func ExampleFloat32ToBytes() {
	x := float32(-561.2863) // -561.2863, 0xc40c5253
	bsBig := Float32ToBytes(x, optBig)
	fmt.Printf("%#02v\n", bsBig)
	bsLittle := Float32ToBytes(x, optLittle)
	fmt.Printf("%#02v\n", bsLittle)
	// Output: []byte{0xc4, 0x0c, 0x52, 0x53}
	// []byte{0x53, 0x52, 0x0c, 0xc4}
}

func ExampleBytesToFloat32() {
	xs := []byte{0xc4, 0x0c, 0x52, 0x53} // -561.2863, 0xc40c5253
	fx, err := BytesToFloat32(xs, optBig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", fx)
	// Output: -561.2863
}

func ExampleFloat64ToBytes() {
	x := float64(-561.2863) // -561.2863, 0xc0818a4a57a786c2
	bs := Float64ToBytes(x, optBig)
	fmt.Printf("%#02v\n", bs)
	// Output: []byte{0xc0, 0x81, 0x8a, 0x4a, 0x57, 0xa7, 0x86, 0xc2}
}

func ExampleBytesToFloat64() {
	xs := []byte{0xc0, 0x81, 0x8a, 0x4a, 0x57, 0xa7, 0x86, 0xc2} // -561.2863, 0xc0818a4a57a786c2
	fx, err := BytesToFloat64(xs, optBig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", fx)
	// Output: -561.2863

}
