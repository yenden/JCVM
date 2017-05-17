package test

/*
import (
	"JCVM/jcre/api/javacard/framework"
	"testing"
)

func TestAIDfunctions(t *testing.T) {
	array := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x00, 0x03}
	array2 := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x01, 0x00}
	aid := framework.InitAID(array, 0, 8)
	aid2 := framework.InitAID(array2, 0, 7)
	if !framework.Equals(aid, array, 0, 8) {
		t.Error("Erreur on aid, arrays not matched")
	}
	cond := framework.Equals(aid, array, 0, 4)
	if cond {
		t.Error("Erreur on aids length")
	}
	cond2 := aid.RidEquals(aid2)
	if !cond2 {
		t.Error("Erreur: RIDs don't match")
	}
}
*/
