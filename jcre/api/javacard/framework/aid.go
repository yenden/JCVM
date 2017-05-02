package framework

import (
	"JCVM/jcre/api/com/sun/javacard/impl"
	"log"
	"reflect"
)

type AID struct {
	theAID []byte
	bArray []byte
}

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

func (aid *AID) equals(object interface{}) bool {
	if object == nil {
		return false
	}
	if aid.bArray == nil {
		aid.bArray = make([]byte, 16) //TODO JCSYstem.makeTransientByteArray((short) 16, JCSystem.CLEAR_ON_RESET);
	}
	switch object.(type) {
	case *AID:
		object.(*AID).getBytes(aid.bArray, 0)
		return reflect.DeepEqual(aid.bArray, aid.theAID)
	default:
		return false
	}
}
func equals(aid *AID, bArray []byte, offset int16, length int16) bool {
	err := impl.CheckArrayArgs(bArray, offset, int16(len(aid.theAID)))
	if err == nil {
		log.Fatal(err)
	}
	return int16(len(aid.theAID)) == length && reflect.DeepEqual(aid.bArray, aid.theAID)
}
func (aid *AID) partialEquals(bArray []byte, offset int16, length int16) bool {
	err := impl.CheckArrayArgs(bArray, offset, int16(len(aid.theAID)))
	if err == nil {
		log.Fatal(err)
	}
	return reflect.DeepEqual(aid.bArray, aid.theAID)
}
func (aid *AID) ridEquals(otherAID *AID) bool {
	if otherAID == nil {
		return false
	}
	if aid.bArray == nil {
		aid.bArray = make([]byte, 16) //TODO JCSYstem.makeTransientByteArray((short) 16, JCSystem.CLEAR_ON_RESET);
	}
	otherAID.getBytes(aid.bArray, 0)
	return reflect.DeepEqual(aid.bArray[:5], aid.theAID[:5])

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
