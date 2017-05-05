package test

import (
	"JCVM/jcre/api/javacard/framework"
	"testing"
)

func TestArrCopy(t *testing.T) {
	src := make([]byte, 4)
	dest := make([]byte, 6)
	src[0] = 0
	src[1] = 1
	src[2] = 2
	src[3] = 3
	framework.ArrayCopy(src, 1, dest, 2, 2)
	if dest[2] != 1 && dest[3] != 3 {
		t.Error("Problem Array copy")
	}
}
