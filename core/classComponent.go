package core

/*ClassComponent represents the CAP class components*/
type ClassComponent struct {
	interfaces []*InterfaceInfo
	pClasses   []*ClassInfo
}

/*TypeDescriptor describes the type of a field*/
type TypeDescriptor struct {
	nibbleCount uint8
	typeArray   []uint8
}

func checkBitField(bitfield uint8) (bool, int) {
	//isShareable := ((bitfield & 0x40) == 0x40)
	isInterf := ((bitfield & 0x80) == 0x80)
	iinterfaceCount := (int)(bitfield & 0x0F)
	return isInterf, iinterfaceCount
}

/*InterfaceInfo ... See CAP components*/
type InterfaceInfo struct {
	bitfield        uint8
	superinterfaces []uint16 //classref struct represent by uint16
}

/*ImplementedInterfaceInfo ... See CAP components */
type ImplementedInterfaceInfo struct {
	interfaces uint16
	count      uint8
	index      []uint8
}

/*ClassInfo is the class description in a package*/
type ClassInfo struct {
	bitfield                uint8
	superClassRef           uint16
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
