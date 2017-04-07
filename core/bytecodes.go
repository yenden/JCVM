package core

import (
	"fmt"
)

func aconstNull(currF *Frame) {
	currF.push(NullType(0))
}
func iconst(currF *Frame, value int) {
	currF.push(int32(value))
}
func bipush(currF *Frame, byte1 uint8) {
	currF.push(int32(byte1))
}
func sipush(currF *Frame, sValue int16) {
	currF.push(int32(sValue))
}
func aload(currF *Frame, index uint8) {
	val := currF.localvariables[index]
	currF.push(val)
}
func iload(currF *Frame, index uint8) {
	val := currF.localvariables[index]
	currF.push(val.(int32))
}
func aaload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) { // the reference point to an array
	case *ArrayValue:
		if value.componentType == TypeReference {
			c := value.array[index.(int16)]
			currF.push(c.(Reference))
		}
	}
}
func baload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeByte || value.componentType == TypeBoolean {
			c := value.array[index.(int16)]
			currF.push(int16(c.(int8)))
		}
	}

}
func saload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeShort {
			c := value.array[index.(int16)]
			currF.push(c.(int16))
		}
	}
}
func astore(currF *Frame, index uint8) {
	ref := currF.pop()
	switch ref.(type) {
	case Reference:
		currF.localvariables[index] = ref.(Reference)
	case ReturnAddress:
		currF.localvariables[index] = ref.(ReturnAddress)
	}
}
func istore(currF *Frame, index uint8) {
	val := currF.pop()
	currF.localvariables[index] = val.(int32)

}
func aastore(currF *Frame) {
	refval := currF.pop()
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeReference { //an array of reference
			value.array[index.(int16)] = refval.(Reference)
		}

	}
}
func bastore(currF *Frame) {
	refval := currF.pop()
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeByte || value.componentType == TypeBoolean { //an array of byte or boolean
			value.array[index.(int16)] = int8(refval.(int16))
		}
	}
}

func sastore(currF *Frame) {
	refval := currF.pop()
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeShort { //an array of byte or boolean
			value.array[index.(int16)] = refval.(int16)
		}

	}

}
func popBytecode(currF *Frame) interface{} {
	interm := currF.operandStack[currF.opStackTop]
	switch interm.(type) {
	case int16:
		return currF.pop()
	}
	return NullType(0)
}
func dup(currF *Frame) {
	interm := currF.operandStack[currF.opStackTop]
	currF.push(interm.(int16))

}
func dup2(currF *Frame) {
	interm1 := currF.operandStack[currF.opStackTop]
	interm2 := currF.operandStack[currF.opStackTop-1]
	currF.push(interm2.(int16))
	currF.push(interm1.(int16))

}
func dupX(currF *Frame) {
	//todo
}
func iadd(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) + value2.(int32)
	currF.push(result)
}
func isub(currF *Frame) {
	value2 := currF.pop()
	value1 := currF.pop()
	result := value1.(int32) - value2.(int32)
	currF.push(result)
}
func imul(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) * value2.(int32)
	currF.push(result)
}
func idiv(currF *Frame) {
	value2 := currF.pop()
	value1 := currF.pop()
	result := value1.(int32) / value2.(int32)
	currF.push(result)
}
func irem(currF *Frame) {
	value2 := currF.pop()
	value1 := currF.pop()
	result := value1.(int32) - (value1.(int32)/value2.(int32))*value2.(int32)
	currF.push(result)
}
func ishl(currF *Frame) {
	value2 := currF.pop()
	value1 := currF.pop()
	s := uint(value2.(int32) & 0x0000001F)
	result := value1.(int32) << s
	currF.push(result)
}
func iushr(currF *Frame) {
	value2 := currF.pop()
	value1 := currF.pop()
	s := uint(value2.(int32) & 0x0000001F)
	result := value1.(int32) >> s
	currF.push(result)
}
func iand(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) & value2.(int32)
	currF.push(result)
}
func ior(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) | value2.(int32)
	currF.push(result)
}
func ixor(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) ^ value2.(int32)
	currF.push(result)
}
func i2b(currF *Frame) {
	value := currF.pop()
	interm := int16(value.(int32) & 0x0000FFFF)
	result := int16(int8(interm))
	currF.push(result)
}
func i2s(currF *Frame) {
	value := currF.pop()
	interm := int16(value.(int32) & 0x0000FFFF)
	currF.push(interm)
}
func ifeq(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) == 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func ifne(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) != 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func iflt(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) < 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func ifge(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) >= 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func ifgt(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) > 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func ifle(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(int16) <= 0 {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}

func ifnull(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(NullType) == NullType(0) {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}

func ifnnnull(currF *Frame, branch int8, pPC *int) {
	value := currF.pop()
	if value.(Reference) != Reference(NullType(0)) {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func goTo(currF *Frame, branch int8, pPC *int) {
	(*pPC) += int(branch)
	(*pPC) -= 2
}
func areturn(currF *Frame, invokerF *Frame) {
	objref := currF.pop()
	invokerF.push(objref.(Reference))
}
func ireturn(currF *Frame, invokerF *Frame) {
	objref := currF.pop()
	invokerF.push(objref.(int32))
}
func invokevirtual(currF *Frame, index uint16, pCA *AbstractApplet, vm *VM) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	if byte1&0x80 == 0x80 {
		packageIndex := byte1 & 0x7F
		pPI := pCA.PImport.packages[packageIndex]
		pCL := findLibrary(pPI)
		if pCL != nil {
			//External library
			classtoken := pCI.info[1]
			sOffset := pCL.AbsA.PExport.pClassExport[classtoken].classOffset
			pcLInf := pCL.AbsA.PClass.pClasses[sOffset]
			token := pCI.info[2]
			index2 := token - pcLInf.publicMethodTableBase
			newFrame := &Frame{}
			vm.pushFrame(newFrame)
			pCL.AbsA.PMethod.executeByteCode(pcLInf.publicVirtualMethodTable[index2], pCL.AbsA, vm, true)
		} else {
			fmt.Println("Error: cannot invoke virtual package not found")
		}

	} else { //Interal class library
		offset := makeU2(pCI.info[0], pCI.info[1])
		token := pCI.info[2]
		pcLInf := pCA.PClass.pClasses[offset]
		newFrame := &Frame{}
		vm.pushFrame(newFrame)
		if token&0x80 == 0x80 {
			index2 := token - pcLInf.packageMethodTableBase
			pCA.PMethod.executeByteCode(pcLInf.packageVirtualMethodTable[index2], pCA, vm, true)
		} else {
			index2 := token - pcLInf.publicMethodTableBase
			pCA.PMethod.executeByteCode(pcLInf.publicVirtualMethodTable[index2], pCA, vm, true)
		}
	}

}
func invokespecial(currF *Frame, index uint16, pCA *AbstractApplet, vm *VM) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	if pCI.tag == 0x06 { //static method
		if byte1 == 0x00 {
			sOffset := makeU2(pCI.info[1], pCI.info[2])
			pCA.PMethod.executeByteCode(sOffset, pCA, vm, true)
		} else {
			packageIndex := byte1 & 0x7F
			pPI := pCA.PImport.packages[packageIndex]
			pCL := findLibrary(pPI)
			if pCL != nil {
				//External library
				classtoken := pCI.info[1]
				sOffset := pCL.AbsA.PExport.pClassExport[classtoken].classOffset
				pcLInf := pCL.AbsA.PClass.pClasses[sOffset]
				token := pCI.info[2]
				index2 := token - pcLInf.publicMethodTableBase
				newFrame := &Frame{}
				vm.pushFrame(newFrame)
				pCL.AbsA.PMethod.executeByteCode(pcLInf.publicVirtualMethodTable[index2], pCL.AbsA, vm, true)
			} else {
				fmt.Println("Error: cannot invoke special package not found")
			}
		}
	}
	if pCI.tag == 0x04 { //super method virtual
		invokevirtual(currF, index, pCA, vm)
	}
}
func invokestatic(currF *Frame, index uint16, pCA *AbstractApplet, vm *VM) {
	invokespecial(currF, index, pCA, vm)
}
func invokeinterface(currF *Frame, pCA *AbstractApplet, vm *VM, nargs uint8, index uint16, methodToken uint8) {
	pCI := pCA.PConstPool.pConstantPool[index] //pCI is a reference to an interface
	interfaceRef := makeU2(pCI.info[1], pCI.info[2])
	newFrame := &Frame{}
	newFrame.opStackTop = -1
	newFrame.localvariables = make([]interface{}, 100)
	newFrame.operandStack = make([]interface{}, 100)
	for i := nargs; i > 0; i-- {
		newFrame.localvariables[i-1] = currF.pop()
	}
	vm.pushFrame(newFrame)
	objref := newFrame.localvariables[0].(Reference) //the object reference
	byte1 := uint8((objref & 0xFF00) >> 4)
	byte2 := uint8(objref & 0x00FF)
	//todo array type
	var token uint8
	if byte1&0x80 == 0x80 {
		packageIndex := byte1 & 0x7F
		pPI := pCA.PImport.packages[packageIndex]
		pCL := findLibrary(pPI)
		if pCL != nil {
			//External library
			classtoken := byte2
			sOffset := pCL.AbsA.PExport.pClassExport[classtoken].classOffset
			pcLInf := pCL.AbsA.PClass.pClasses[sOffset]
			for j := 0; j < len(pcLInf.interfaces); j++ {
				if pcLInf.interfaces[j].interfaces == interfaceRef {
					token = pcLInf.interfaces[j].index[int(methodToken)]
					break
				}
			}
			index2 := token - pcLInf.publicMethodTableBase
			ipos := int(index2)
			flag := readU1(pCL.AbsA.PMethod.pMethodInfo, &ipos)
			if isExtended(flag) {
				ipos += 3
			} else {
				ipos++
			}
			ipos2 := int(pcLInf.publicVirtualMethodTable[ipos])
			vm.runStatic(pCL.AbsA.PMethod.pMethodInfo, &ipos2, pCL.AbsA, nargs)
		} else {
			fmt.Println("Error: cannot invoke interface package not found")
		}
	} else {
		offset := makeU2(byte1, byte2)
		pcLInf := pCA.PClass.pClasses[offset]
		for j := 0; j < len(pcLInf.interfaces); j++ {
			if pcLInf.interfaces[j].interfaces == interfaceRef {
				token = pcLInf.interfaces[j].index[int(methodToken)]
				break
			}
		}
		if token&0x80 == 0x80 {
			index2 := token - pcLInf.packageMethodTableBase
			ipos := int(index2)
			flag := readU1(pCA.PMethod.pMethodInfo, &ipos)
			if isExtended(flag) {
				ipos += 3
			} else {
				ipos++
			}
			ipos2 := int(pcLInf.packageVirtualMethodTable[ipos])
			vm.runStatic(pCA.PMethod.pMethodInfo, &ipos2, pCA, nargs)

		} else {
			index2 := token - pcLInf.publicMethodTableBase
			ipos := int(index2)
			flag := readU1(pCA.PMethod.pMethodInfo, &ipos)
			if isExtended(flag) {
				ipos += 3
			} else {
				ipos++
			}
			ipos2 := int(pcLInf.publicVirtualMethodTable[ipos])
			vm.runStatic(pCA.PMethod.pMethodInfo, &ipos2, pCA, nargs)
		}
	}
}

/*
*****Byte code verification for later*****
func verifyInvokeVirtual(index uint16, pCA *AbstractApplet) bool {
	var protected bool = false
	var class *ClassDescriptorInfo
	pCI := pCA.PConstPool.pConstantPool[index]
	//verify that it is a virtual method
	if pCI.tag == 3 {
		byte1 := pCI.info[0]
		// verify that this method is in the current package or external package
		superClassRef := pCA.PClass.pClasses[classref].superClassRef
		//PCLInf := pCA.PClass.pClasses[superClassRef]
		if byte1&0x80 == 0x80 { //external package
			if class.methods[j].pAF&AccProtected == AccProtected {
				protected = true
			}
			break

		} else {
			//internal package
			classref := makeU2(pCI.info[0], pCI.info[1])
			token := pCI.info[2]
			for i := 0; i < int(pCA.PDescriptor.classCount); i++ {
				if pCA.PDescriptor.classes[i].thisClassRef == classref {
					class = pCA.PDescriptor.classes[i]
					break
				}
			}
			for j := 0; j < int(class.methodCount); j++ {
				if class.methods[j].token == token {
					if class.methods[j].pAF&AccInit == AccInit {
						return false
					}
				}
			}

		}
	}
}
*/
