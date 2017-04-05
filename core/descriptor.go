package core

type FieldRef struct {
	fieldRef [3]uint8
}

func isPrimitiveType(value uint16) bool {
	return (value & 0x8000) == 0x8000
}

func isReferenceType(value uint16) bool {
	return (value & 0x8000) == 0x0
}

type FieldDescriptorInfo struct {
	token      uint8
	pAF        uint8
	pFieldRef  *FieldRef
	pFieldtype uint16
}

type MethodDescriptorInfo struct {
	token                 uint8
	pAF                   uint8 //access flag
	methodOffset          uint16
	typeOffset            uint16
	bytecodeCount         uint16
	exceptionHandlerCount uint16
	exceptionHandlerIndex uint16
}

type ClassDescriptorInfo struct {
	token          uint8
	accessFlags    uint8
	thisClassRef   uint16
	interfaceCount uint8
	fieldCount     uint16
	methodCount    uint16

	interfaces []uint16
	fields     []*FieldDescriptorInfo
	methods    []*MethodDescriptorInfo
}

type TypeDescriptorInfo struct {
	constPoolCount uint16
	//typeDescCount int //Not a standard member
	pConstantPoolTypes []uint16
	pTypeDesc          []*TypeDescriptor
}
type DescriptorComponent struct {
	classCount uint8
	classes    []*ClassDescriptorInfo
	types      *TypeDescriptorInfo
}
