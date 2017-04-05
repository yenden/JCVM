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

/*
func (cpc *ConstantPoolComponent) createCPC() {
	cpc.pConstantPool = make([]*CpInfo, cpc.count)
}*/ /*
func (cpc *ConstantPoolComponent) addConst(tag uint8, byte1 uint8, byte2 uint8, byte3 uint8) {
	info := make([]uint8, 3)
	info[0] = byte1
	info[1] = byte2
	info[2] = 0
	switch tag {
	case 1:
		cpc.pConstantPool[cpCount] = &CpInfo{tag, info}
	default:
		info[2] = byte3
		cpc.pConstantPool[cpCount] = &CpInfo{tag, info}
	}
	cpCount++
}*/
