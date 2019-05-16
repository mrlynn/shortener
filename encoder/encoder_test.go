package encoder

import (
	"fmt"
	"testing"

	fuzz "github.com/google/gofuzz"
)

func TestEncode(t *testing.T) {
	for i := 0; i < 10; i++ {
		f := fuzz.New()

		var testInput uint32

		f.Fuzz(&testInput)

		fmt.Println(Encode(int64(testInput)), testInput)
	}
}
