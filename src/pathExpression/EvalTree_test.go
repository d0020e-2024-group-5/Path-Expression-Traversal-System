package pathExpression

import (
	"log"
	"testing"
)

var str string

func TestIsValid(t *testing.T) {
	str = "s/bomba*}f{*y"
	got := IsValid(str)
	want := true

	if got != want {
		log.Print(got)
		t.Errorf("got %t, wanted %t", got, want)
	}

}
// Checks for invalid operator combinations and returns nil if no invalid combination is found.
