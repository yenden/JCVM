package core

import "JCVM/jcre/api"

type typeValue uint16

/*Represent an array's contents type*/
const (
	typeBoolean typeValue = 0x0001 + iota
	typeByte
	typeShort
	typeInt
	typeReference
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
	javaClassArray [256]*JavaClass
	//jcCount is Java Class count in the main heap
	jcCount = -1
	//heap represent the main memory heap
	heap = make(map[Reference]interface{})
	//InstanceRefHeap is a map between AID and memry reference in the main heap
	InstanceRefHeap = make(map[*api.AID]Reference)
	//arrcount indexes arrays in the main heap
	arrcount = 511
	//interfcount indexes interfaces in the main heap
	interfcount = 255
	apduBuffLen = 0
)

/*ArrayValue represents an array in the heap*/
type ArrayValue struct {
	componentType typeValue
	length        uint16
	array         interface{} // []interface{}
}

/*JavaClass represents an instance class in the heap*/
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

//InitApduArr inits apdu array for incoming and outgoing apdus
//It will be called once
func InitApduArr() {
	array := &ArrayValue{}
	array.componentType = typeByte
	array.length = uint16(37)
	array.array = make([]byte, 37)
	heap[Reference(6000)] = array
}

/*FillApduArr fills the apdu array in the heap with the receives apdu buffer*/
func FillApduArr(apdu []byte, ref Reference) {
	//Reference 6000 has been fixed for apdu array in the heap
	arr := (heap[Reference(6000)]).(*ArrayValue)
	for i := 0; i < len(apdu); i++ {
		arr.array.([]byte)[i] = apdu[i]
	}
	apduBuffLen = len(apdu)
}
