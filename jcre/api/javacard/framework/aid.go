package framework

import (
	"JCVM/jcre/api/com/sun/javacard/impl"
	"JCVM/jcre/api/share"
	"log"
	"reflect"
)

type AID struct {
	theAID []byte
}

var bArray []byte

func InitAID(bArray []byte, offset int16, length int16) *AID /*(*AID, error)*/ {
	aid := &AID{}
	err := impl.CheckArrayArgs(bArray, offset, length)
	if err == nil {
		log.Fatal(err)
	}
	if length < 5 || length > 16 {
		log.Fatal("AID Length is not correct")
	}
	aid.theAID = make([]byte, length)
	ArrayCopy(bArray, offset, aid.theAID, 0, length)
	return aid
}
func (aid *AID) getBytes(dest []byte, offset int16) byte {
	err := impl.CheckArrayArgs(dest, offset, int16(len(aid.theAID)))
	if err == nil {
		log.Fatal(err)
	}
	ArrayCopy(aid.theAID, 0, dest, offset, int16(len(aid.theAID)))
	return byte(len(aid.theAID))
}

func (aid *AID) Equals(object interface{}) bool {
	if object == nil {
		return false
	}
	if bArray == nil {
		bArray = share.MakeTransientByteArray(16, share.ClearOnReset)
	}
	switch object.(type) {
	case *AID:
		object.(*AID).getBytes(bArray, 0)
		return reflect.DeepEqual(bArray, aid.theAID)
	default:
		return false
	}
}
func Equals(aid *AID, bArray []byte, offset int16, length int16) bool {
	err := impl.CheckArrayArgs(bArray, offset, int16(len(aid.theAID)))
	if err == nil {
		log.Fatal(err)
	}
	return int16(len(aid.theAID)) == length && reflect.DeepEqual(bArray, aid.theAID)
}
func (aid *AID) partialEquals(bArray []byte, offset int16, length int16) bool {
	err := impl.CheckArrayArgs(bArray, offset, int16(len(aid.theAID)))
	if err == nil {
		log.Fatal(err)
	}
	return reflect.DeepEqual(bArray, aid.theAID)
}
func (aid *AID) RidEquals(otherAID *AID) bool {
	if otherAID == nil {
		return false
	}
	if bArray == nil {
		bArray = share.MakeTransientByteArray(16, share.ClearOnReset)
	}
	otherAID.getBytes(bArray, 0)
	return reflect.DeepEqual(bArray[:5], aid.theAID[:5])

}
func (aid *AID) getPartialBytes(aidOffset int16, dest []byte, oOffset int16, oLength int16) {
	copyLen := oLength
	if oLength == 0 {
		copyLen = int16(len(aid.theAID)) - aidOffset
	}
	err := impl.CheckArrayArgs(dest, oOffset, copyLen)
	if err == nil {
		log.Fatal(err)
	}
	ArrayCopy(aid.theAID, aidOffset, dest, oOffset, copyLen)
}
