package core

var cpCount int = 0

type CpInfo struct {
	tag  uint8
	info []uint8
}

type ConstantPoolComponent struct {
	count         uint16
	pConstantPool []*CpInfo
}
