package core

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
	javaClassArray [256]*JavaClass
	jcCount        = 0
	heap           = make(map[Reference]interface{})
	arrcount       = 511
	interfcount    = 255
	//heap           = make(map[int16]*ArrayValue)
	//heapClass      = make(map[int16]*JavaClass)
)

type ArrayValue struct {
	componentType typeValue
	length        uint16
	array         []interface{}
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

/*
func allocate(object interface{}, ref int16) {
	heap[ref] = object
	switch value := object.(type) {
	case *ArrayValue:
		value.allocateArray()
	case *JavaClass:
		//createClass()
	default:
		//nothing
	}
}
func (arr *ArrayValue) allocateArray() {
	if arr.componentType == TypeBoolean {
		c := make([]bool, arr.length)
		for i := range c {
			arr.array[i] = c[i]
		}
	}
	if arr.componentType == TypeByte {
		c := make([]uint8, arr.length)
		for i := range c {
			arr.array[i] = c[i]
		}
	}
	if arr.componentType == TypeShort {
		c := make([]uint16, arr.length)
		for i := range c {
			arr.array[i] = c[i]
		}
	}
	if arr.componentType == TypeInt {
		c := make([]uint32, arr.length)
		for i := range c {
			arr.array[i] = c[i]
		}
	}
	if arr.componentType == TypeReference {
		c := make(map[uint16]*JavaClass, arr.length)
		for i := range c {
			arr.array[i] = c[i]
		}
	}
}

type ArrayValue struct {
	componentType typeValue
	length        uint8
	array         []interface{}
}
type JavaClass struct {
	superclassref        uint16
	declaredinstancesize uint8
	fields               []uint16
	fieldInit            []bool
}

func createClass(superclassref uint16, declaredinstancesize uint8) uint16 {
	fieldtab := make([]uint16, int(declaredinstancesize))
	fieldinittab := make([]bool, int(declaredinstancesize))
	jcCount++
	javaClassArray[jcCount] = &JavaClass{superclassref, declaredinstancesize, fieldtab, fieldinittab}
	return (uint16(jcCount)) /*+0x100
}*/

//type Object interface{}
