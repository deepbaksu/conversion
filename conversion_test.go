package conversion_test

import (
	"fmt"
	"github.com/dl4ab/conversion"
)

func ExampleIntToFloat32() {
	var x int = 3
	fmt.Println(conversion.IntToFloat32(x))
	// Output: 3 <nil>
}
