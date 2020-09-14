package conversion

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"

	"testing"

	"log"
)

const (
	binaryWidthTwo = "%#02v\n"
)

var (
	optBig    = &Option{Endian: BigEndian}
	optLittle = &Option{Endian: LittleEndian}
)

var (
	testInputF64      = float64(-561.2863)
	testInputF64Bytes = []byte{0xc0, 0x81, 0x8a, 0x4a, 0x57, 0xa7, 0x86, 0xc2}
)

var (
	testInputF32      = float32(-561.2863)
	testInputF32Bytes = []byte{0xc4, 0x0c, 0x52, 0x53}
)

func ExampleIntToFloat32() {
	var x int = 3
	fmt.Println(IntToFloat32(x))
	// Output: 3 <nil>
}

// areSameErrors returns true if both err and err2 are <nil> or have the same error message.
func areSameErrors(err error, err2 error) bool {
	if err != nil && err2 != nil {
		return err.Error() == err2.Error()
	}
	return err == err2
}

func ExampleFloat32ToBytes() {
	x := float32(-561.2863) // -561.2863, 0xc40c5253
	bsBig := Float32ToBytes(x, optBig)
	fmt.Printf(binaryWidthTwo, bsBig)
	bsLittle := Float32ToBytes(x, optLittle)
	fmt.Printf(binaryWidthTwo, bsLittle)
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

func TestBytesToFloat32_failBadInput(t *testing.T) {
	// bad input as 32bits should have a length of 4.
	xs := []byte{0xc4, 0x0c, 0x52}
	_, err := BytesToFloat32(xs, optBig)
	assert.EqualError(t, err, "length of []byte should be 4 or bigger")
}

func TestBytesToFloat32_littleEndian(t *testing.T) {
	xs := []byte{0x53, 0x52, 0x0c, 0xc4} // -561.2863
	fx, err := BytesToFloat32(xs, optLittle)
	assert.NoError(t, err)
	assert.InEpsilon(t, -561.2863, fx, 1e-4)
}

func ExampleFloat64ToBytes() {
	x := -561.2863 // -561.2863, 0xc0818a4a57a786c2
	bs := Float64ToBytes(x, optBig)
	fmt.Printf(binaryWidthTwo, bs)
	// Output: []byte{0xc0, 0x81, 0x8a, 0x4a, 0x57, 0xa7, 0x86, 0xc2}
}

func TestFloat64ToBytes_littleEndian(t *testing.T) {
	x := -561.2863
	bs := Float64ToBytes(x, optLittle)
	assert.Equal(t, []byte{0xc2, 0x86, 0xa7, 0x57, 0x4a, 0x8a, 0x81, 0xc0}, bs)
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

func TestBytesToFloat64_invalidInput(t *testing.T) {
	invalidInput := []byte{0x1, 0x2, 0x3}
	_, err := BytesToFloat64(invalidInput, nil)
	assert.EqualError(t, err, "length of []byte should be 8 or bigger")
}

func TestBytesToFloat64_littleEndian(t *testing.T) {
	xs := []byte{0xc2, 0x86, 0xa7, 0x57, 0x4a, 0x8a, 0x81, 0xc0} // -561.2863, 0xc0818a4a57a786c2
	fx, err := BytesToFloat64(xs, optLittle)
	assert.NoError(t, err)
	assert.InEpsilon(t, -561.2863, fx, 1e-4)
}

func TestFloat32sToBytes(t *testing.T) {
	t.Parallel()

	type args struct {
		xs []float32
		o  *Option
	}

	testCases := map[string]struct {
		args    args
		want    []byte
		wantErr error
	}{
		"it should convert []float32 to []byte": {
			args:    args{[]float32{float32(-561.2863)}, nil},
			want:    []byte{0xc4, 0x0c, 0x52, 0x53},
			wantErr: nil,
		},
		"it should return empty []byte if []float32 is empty ": {
			args:    args{[]float32{}, nil},
			want:    []byte{},
			wantErr: nil,
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got, got1 := Float32sToBytes(tc.args.xs, tc.args.o)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("Float32sToBytes() got = %v, want %v", got, tc.want)
			}
			if !areSameErrors(got1, tc.wantErr) {
				t.Errorf("Float32sToBytes() got1 = %v, want %v", got1, tc.wantErr)
			}
		})
	}
}

func TestBytesToFloat32s(t *testing.T) {
	t.Parallel()

	type args struct {
		xs []byte
		o  *Option
	}

	testCases := map[string]struct {
		args    args
		want    []float32
		wantErr bool
	}{
		"it should convert []byte to []float32": {
			args: args{
				xs: testInputF32Bytes,
				o:  nil,
			},
			want:    []float32{testInputF32},
			wantErr: false,
		},
		"it should return an error when the input is invalid": {
			args: args{
				xs: []byte{0x0c, 0x52, 0x53},
				o:  nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got, err := BytesToFloat32s(tc.args.xs, tc.args.o)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFloat64sToBytes(t *testing.T) {
	t.Parallel()

	type args struct {
		xs []float64
		o  *Option
	}

	testCases := map[string]struct {
		name    string
		args    args
		want    []byte
		wantErr error
	}{
		"it should convert []float64 to []byte": {
			args:    args{[]float64{testInputF64}, nil},
			want:    testInputF64Bytes,
			wantErr: nil,
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got, got1 := Float64sToBytes(tc.args.xs, tc.args.o)
			if !bytes.Equal(got, tc.want) {
				t.Errorf("Float64sToBytes() got = %v, want %v", got, tc.want)
			}
			if !areSameErrors(got1, tc.wantErr) {
				t.Errorf("Float64sToBytes() got1 = %v, want %v", got1, tc.wantErr)
			}
		})
	}
}

func TestBytesToFloat64s(t *testing.T) {
	t.Parallel()

	type args struct {
		xs []byte
		o  *Option
	}

	testCases := map[string]struct {
		args    args
		want    []float64
		wantErr bool
	}{
		"it should convert []byte to []float64": {
			args: args{
				xs: testInputF64Bytes,
				o:  nil,
			},
			want:    []float64{testInputF64},
			wantErr: false,
		},
		"it should fail if the input is not []float64": {
			args: args{
				xs: []byte{0x1, 0x2, 0x3},
				o:  nil,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			got, err := BytesToFloat64s(tc.args.xs, tc.args.o)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("BytesToFloat64s() got = %v, want %v", got, tc.want)
			}
		})
	}
}
