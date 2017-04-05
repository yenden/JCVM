package core

type ClassComponent struct {
	interfaces []*InterfaceInfo
	pClasses   []*ClassInfo
}

//TODO create a  number of interfaces and classes

type TypeDescriptor struct {
	nibbleCount uint8
	typeArray   []uint8
}

/*
func (td *TypeDescriptor) isVoid() bool {
	return (td.typeArray[0]&0xF0 == TypeVoid)
}
func (td *TypeDescriptor) isBoolean() bool {
	return (td.typeArray[0]&0xF0 == TypeBoolean)
}
func (td *TypeDescriptor) isByte() bool {
	return (td.typeArray[0]&0xF0 == TypeByte)
}
func (td *TypeDescriptor) isShort() bool {
	return (td.typeArray[0]&0xF0 == TypeShort)
}
func (td *TypeDescriptor) isInt() bool {
	return (td.typeArray[0]&0xF0 == TypeInt)
}
func (td *TypeDescriptor) isRef() bool {
	return (td.typeArray[0]&0xF0 == TypeReference)
}
func (td *TypeDescriptor) isArrayBool() bool {
	return (td.typeArray[0]&0xF0 == TypeArrayOfBoolean)
}
func (td *TypeDescriptor) isArrayByte() bool {
	return (td.typeArray[0]&0xF0 == TypeArrayofByte)
}
func (td *TypeDescriptor) isArrayShort() bool {
	return (td.typeArray[0]&0xF0 == TypeArrayOfShort)
}
func (td *TypeDescriptor) isArrayInt() bool {
	return (td.typeArray[0]&0xF0 == TypeArrayOfInt)
}
func (td *TypeDescriptor) isArrayRef() bool {
	return (td.typeArray[0]&0xF0 == TypeArrayOfReference)
}

/*type AbstractClassInfo struct {
	isShareable     bool
	isInterf        bool
	iinterfaceCount int
}*/

func checkBitField(bitfield uint8) (bool, int) {
	//isShareable := ((bitfield & 0x40) == 0x40)
	isInterf := ((bitfield & 0x80) == 0x80)
	iinterfaceCount := (int)(bitfield & 0x0F)
	return isInterf, iinterfaceCount
}

type InterfaceInfo struct {

	//bitfield        *AbstractClassInfo
	bitfield        uint8
	superinterfaces []uint16 //classref struct represent by uint16
}
type ImplementedInterfaceInfo struct {
	interfaces uint16
	count      uint8
	index      []uint8
}
type ClassInfo struct {
	//bitfield                *AbstractClassInfo
	bitfield      uint8
	superClassRef uint16
	//	superClassRef           ClassRef
	declaredInstanceSize    uint8
	firstReferenceToken     uint8
	referenceCount          uint8
	publicMethodTableBase   uint8
	publicMethodTableCount  uint8
	packageMethodTableBase  uint8
	packageMethodTableCount uint8

	publicVirtualMethodTable  []uint16
	packageVirtualMethodTable []uint16
	interfaces                []*ImplementedInterfaceInfo
}

/*func isInterface(bitfield uint8) bool {
	if bitfield&0x80 == 0x80 {
		return true
	}
	return false
}*/

/*func (cl *ClassInfo) isInterface() bool {
	return false
}
func (interf *InterfaceInfo) isInterface() bool {
	return true
}*/
