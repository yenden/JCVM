package core

import (
	"fmt"
)

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

func (mComp *MethodComponent) executeByteCode(offset uint16, pCA *AbstractApplet, vm *VM, invokercond bool, processCond bool) {
	iPosm2 := int(offset - 1)
	flags := readU1(mComp.pMethodInfo, &iPosm2)
	currFrame := vm.StackFrame[vm.FrameTop]
	var maxStack, nargs, maxLocals uint8
	if isExtended(flags) {
		maxStack = readU1(mComp.pMethodInfo, &iPosm2)
		nargs = readU1(mComp.pMethodInfo, &iPosm2)
		maxLocals = readU1(mComp.pMethodInfo, &iPosm2)
	} else {
		//if abstract
		maxStack = readLow(flags)
		bitField := readU1(mComp.pMethodInfo, &iPosm2)
		nargs = readHighShift(bitField)
		maxLocals = readLow(bitField)
	}
	fmt.Println("max stack", maxStack, "maxlocal", maxLocals)
	if !processCond {
		currFrame.opStackTop = -1
		currFrame.Localvariables = make([]interface{}, 256)
	}
	currFrame.operandStack = make([]interface{}, 256)
	if invokercond == true {
		invokerframe := vm.StackFrame[vm.FrameTop-1]
		for i := nargs; i > 0; i-- {
			currFrame.Localvariables[i-1] = invokerframe.pop()
		}
	}
	vm.runStatic(mComp.pMethodInfo, &iPosm2, pCA, nargs)

}
