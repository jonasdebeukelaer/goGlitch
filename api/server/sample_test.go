package server

import "testing"

func Test_testableFunction(t *testing.T) {
	sq := testableFunction(3)

	if sq != 9 {
		t.Errorf("Squaring didn't work! Expected %d got %v", 9, sq)
	}
}
