package core

import "JCVM/jcre/api"

type typeValue uint16

const (
	TypeBoolean typeValue = 0x0001 + iota
	TypeByte
	TypeShort
	TypeInt
	TypeReference
)

/*TypeVoid             = 0x0001
TypeBoolean          = 0x0002
TypeByte             = 0x0003
TypeShort            = 0x0004
TypeInt              = 0x0005
TypeReference        = 0x0006
TypeArrayOfBoolean   = 0x000A
TypeArrayofByte      = 0x000B
TypeArrayOfShort     = 0x000C
TypeArrayOfInt       = 0x000D
TypeArrayOfReference = 0x000E*/

var (
	javaClassArray  [256]*JavaClass
	jcCount         = -1
	heap            = make(map[Reference]interface{})
	InstanceRefHeap = make(map[*api.AID]Reference)
	arrcount        = 511
	interfcount     = 255
	apduBuffLen     = 0
	//heap           = make(map[int16]*ArrayValue)
	//heapClass      = make(map[int16]*JavaClass)
)

type ArrayValue struct {
	componentType typeValue
	length        uint16
	array         interface{} // []interface{}
}
type JavaClass struct {
	classref             uint16
	superclassref        uint16
	declaredinstancesize uint8
	fields               []*instanceField
	fieldInit            []bool
}
type instanceField struct {
	token uint8
	value interface{}
}

//Apdu array for incoming and outgoing apdus
//It will be called once
func InitApduArr() {
	array := &ArrayValue{}
	array.componentType = TypeByte
	array.length = uint16(37)
	array.array = make([]byte, 37)
	heap[Reference(6000)] = array
}

func FillApduArr(apdu []byte, ref Reference) {
	arr := (heap[Reference(6000)]).(*ArrayValue)
	for i := 0; i < len(apdu); i++ {
		arr.array.([]byte)[i] = apdu[i]
	}
	apduBuffLen = len(apdu)
}
