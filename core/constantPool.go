package core

var cpCount int

/*CpInfo is description of each Constant pool entry */
type CpInfo struct {
	tag  uint8
	info []uint8
}

/*ConstantPoolComponent of the CAP file*/
type ConstantPoolComponent struct {
	count         uint16
	pConstantPool []*CpInfo
}
