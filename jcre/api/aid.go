package api

import (
	"log"
	"reflect"
)

type AID struct {
	TheAID []byte
}

var bArray []byte

func InitAID(bArray []byte, offset int16, length int16) *AID /*(*AID, error)*/ {
	aid := &AID{}
	err := CheckArrayArgs(bArray, offset, length)
	if err == nil {
		log.Fatal(err)
	}
	if length < 5 || length > 16 {
		log.Fatal("AID Length is not correct")
	}
	aid.TheAID = make([]byte, length)
	ArrayCopy(bArray, offset, aid.TheAID, 0, length)
	return aid
}
func (aid *AID) getBytes(dest []byte, offset int16) byte {
	err := CheckArrayArgs(dest, offset, int16(len(aid.TheAID)))
	if err == nil {
		log.Fatal(err)
	}
	ArrayCopy(aid.TheAID, 0, dest, offset, int16(len(aid.TheAID)))
	return byte(len(aid.TheAID))
}

func (aid *AID) Equals(object interface{}) bool {
	if object == nil {
		return false
	}
	if bArray == nil {
		bArray = MakeTransientByteArray(16, ClearOnReset)
	}
	switch object.(type) {
	case *AID:
		object.(*AID).getBytes(bArray, 0)
		return reflect.DeepEqual(bArray, aid.TheAID)
	default:
		return false
	}
}
func Equals(aid *AID, bArray []byte, offset int16, length int16) bool {
	err := CheckArrayArgs(bArray, offset, int16(len(aid.TheAID)))
	if err == nil {
		log.Fatal(err)
	}
	return int16(len(aid.TheAID)) == length && reflect.DeepEqual(bArray, aid.TheAID)
}
func (aid *AID) partialEquals(bArray []byte, offset int16, length int16) bool {
	err := CheckArrayArgs(bArray, offset, int16(len(aid.TheAID)))
	if err == nil {
		log.Fatal(err)
	}
	return reflect.DeepEqual(bArray, aid.TheAID)
}
func (aid *AID) RidEquals(otherAID *AID) bool {
	if otherAID == nil {
		return false
	}
	if bArray == nil {
		bArray = MakeTransientByteArray(16, ClearOnReset)
	}
	otherAID.getBytes(bArray, 0)
	return reflect.DeepEqual(bArray[:5], aid.TheAID[:5])

}
func (aid *AID) getPartialBytes(aidOffset int16, dest []byte, oOffset int16, oLength int16) {
	copyLen := oLength
	if oLength == 0 {
		copyLen = int16(len(aid.TheAID)) - aidOffset
	}
	err := CheckArrayArgs(dest, oOffset, copyLen)
	if err == nil {
		log.Fatal(err)
	}
	ArrayCopy(aid.TheAID, aidOffset, dest, oOffset, copyLen)
}
