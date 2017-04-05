package core

type ExceptionHandlerInfo struct {
	startOffset    uint16
	activeLength   uint16
	handlerOffset  uint16
	catchTypeIndex uint16
}

/*
type MethodHeaderInfo struct {
	flags     uint8
	maxStack  uint8
	nargs     uint8
	maxLocals uint8
}
type MethodInfo struct {
	pMethodHeaderInfo *MethodHeaderInfo
	bytecodes         []uint8
}*/
type MethodComponent struct {
	handlerCount       uint8
	pExceptionHandlers []*ExceptionHandlerInfo
	//pMethodInfo        []*MethodInfo
	pMethodInfo []uint8
}

func isExtended(flag uint8) bool {
	return (flag & 0x80) == 0x80
}
func isAbstract(flag uint8) bool {
	return (flag & 0x40) == 0x40
}

func (mComp *MethodComponent) executeByteCode(offset uint16, pCA *AbstractApplet, vm *VM) {

	iPosm2 := int(offset)
	flags := readU1(mComp.pMethodInfo, &iPosm2)
	if isExtended(flags) {
		maxStack := readU1(mComp.pMethodInfo, &iPosm2)
		nargs := readU1(mComp.pMethodInfo, &iPosm2)
		maxLocals := readU1(mComp.pMethodInfo, &iPosm2)
		vm.runStatic(mComp.pMethodInfo, &iPosm2, pCA, maxStack, nargs, maxLocals)
	} else {
		maxStack := readLow(flags)
		bitField := readU1(mComp.pMethodInfo, &iPosm2)
		nargs := readHighShift(bitField)
		maxLocals := readLow(bitField)
		vm.runStatic(mComp.pMethodInfo, &iPosm2, pCA, maxStack, nargs, maxLocals)
	}
}
