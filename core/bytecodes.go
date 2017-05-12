package core

import (
	"JCVM/jcre/api"
	"JCVM/jcre/nativeMethods"
	"errors"
	"fmt"
	"reflect"
)

var (
	ownerPinMap  = map[uint8]uint8{0: 0, 1: 4, 2: 3, 3: 1, 4: 5, 5: 6, 6: 8, 7: 2, 8: 7}
	frameworkAID = []byte{0xa0, 0, 0, 0, 0x62, 1, 1}
)

const (
	ownerpinClass         = 9
	check                 = 1
	isValidated           = 4
	resetAndUnblock       = 6
	update                = 8
	applet                = 3
	register              = 1
	selectingApplet       = 3
	util                  = 16
	arraycopynonAtomic    = 2
	apdu                  = 10
	getbuffer             = 1
	receiveBytes          = 3
	sendBytes             = 4
	sendBytesLong         = 5
	setIncomingAndReceive = 6
	setOutgoing           = 7
	setOutgoingAndSend    = 8
	setOutgoingLength     = 9
	isoException          = 7
	throwit               = 1
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
func bspush(currF *Frame, value uint8) {
	currF.push(int16(value))
}
func sspush(currF *Frame, value uint16) {
	currF.push(int16(value))
}
func aload(currF *Frame, index uint8) {
	val := currF.Localvariables[index]
	currF.push(val.(Reference))
}
func iload(currF *Frame, index uint8) {
	val := currF.Localvariables[index]
	currF.push(val.(int32))
}
func sload(currF *Frame, index uint8) {
	val := currF.Localvariables[index]
	currF.push(val.(int16))
}
func sconst(currF *Frame, value int8) {
	currF.push(int16(value))

}
func aaload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) { // the reference point to an array
	case *ArrayValue:
		if value.componentType == TypeReference {
			c := value.array.([]Reference)[index.(int16)]
			currF.push(c)
		}
	}
}
func baload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeByte || value.componentType == TypeBoolean {
			c := value.array.([]byte)[index.(int16)]
			currF.push(int16(c))
		}
	}

}
func saload(currF *Frame) {
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeShort {
			c := value.array.([]int16)[index.(int16)]
			currF.push(c)
		}
	}
}
func astore(currF *Frame, index uint8) {
	ref := currF.pop()
	switch ref.(type) {
	case Reference:
		currF.Localvariables[index] = Reference(ref.(Reference))
	case ReturnAddress:
		currF.Localvariables[index] = ref.(ReturnAddress)
	}
}
func istore(currF *Frame, index uint8) {
	val := currF.pop()
	currF.Localvariables[index] = val.(int32)

}
func sstore(currF *Frame, index uint8) {
	val := currF.pop()
	currF.Localvariables[index] = val.(int16)
}
func aastore(currF *Frame) {
	refval := currF.pop()
	index := currF.pop()
	arrayref := currF.pop()
	switch value := heap[arrayref.(Reference)].(type) {
	case *ArrayValue: // the reference point to an array
		if value.componentType == TypeReference { //an array of reference
			value.array.([]Reference)[index.(int16)] = refval.(Reference)
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
			value.array.([]byte)[index.(int16)] = uint8(refval.(int16))
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
			value.array.([]int16)[index.(int16)] = refval.(int16)
		}

	}

}

func sinc(currF *Frame, index uint8, constant int8) {
	interm := currF.Localvariables[index].(int16)
	currF.Localvariables[index] = interm + int16(constant)
	fmt.Println("Result ", currF.Localvariables[index].(int16))
}
func popBytecode(currF *Frame) {
	interm := currF.operandStack[currF.opStackTop]
	switch interm.(type) {
	case int16:
		currF.pop()
	case uint16:
		currF.pop()
	}

}
func dup(currF *Frame) {
	interm := currF.operandStack[currF.opStackTop]
	switch interm.(type) {
	case int16:
		currF.push(interm.(int16))
	case uint16:
		currF.push(interm.(uint16))
	case Reference:
		currF.push(interm.(Reference))
	}
}
func dup2(currF *Frame) {
	interm1 := currF.operandStack[currF.opStackTop]
	interm2 := currF.operandStack[currF.opStackTop-1]
	currF.push(interm2.(int16))
	currF.push(interm1.(int16))

}
func dupX(currF *Frame, mn uint8) {
	m := readHighShift(mn)
	n := readLow(mn)
	if n == 0 {
		interm := make([]interface{}, m)
		for i := 0; i < int(m); i++ {
			switch a := currF.operandStack[currF.opStackTop-i].(type) {
			case int32:
				interm[i] = a
			case int16:
				interm[i] = a
			}
		}
		for i := m; i > 0; i++ {
			currF.push(interm[i])
		}
	} else {
		interm := make([]interface{}, n)
		for i := 0; i < int(n); i++ {
			switch a := currF.pop().(type) {
			case int32:
				interm[i] = a
			case int16:
				interm[i] = a
			}
		}
		for i := m; i > 0; i++ {
			currF.push(interm[i])
		}
		for i := n; i > 0; i++ {
			currF.push(interm[i])
		}
	}
}
func iadd(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int32) + value2.(int32)
	currF.push(result)
}
func sadd(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int16) + value2.(int16)
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
func smul(currF *Frame) {
	value1 := currF.pop()
	value2 := currF.pop()
	result := value1.(int16) * value2.(int16)
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
func iint(currF *Frame, index uint8, constant int8) {
	switch currF.Localvariables[index].(type) {
	case int32:
		currF.Localvariables[index] = int32(constant)
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
func sreturn(currF *Frame, invokerF *Frame) {
	objref := currF.pop()
	invokerF.push(objref.(int16))
}

func invokevirtual(currF *Frame, index uint16, pCA *AbstractApplet, vm *VM) {

	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	if byte1&0x80 == 0x80 {
		packageIndex := byte1 & 0x7F
		pPI := pCA.PImport.packages[packageIndex]
		pCL := findLibrary(pPI)
		var classtoken, token uint8
		if pCL != nil {
			//External library which is not framework
			classtoken = pCI.info[1]
			token = pCI.info[2]
			sOffset := pCL.AbsA.PExport.pClassExport[classtoken].classOffset
			pcLInf := pCL.AbsA.PClass.pClasses[sOffset]
			index2 := token - pcLInf.publicMethodTableBase
			newFrame := &Frame{}
			vm.PushFrame(newFrame)
			pCL.AbsA.PMethod.executeByteCode(pcLInf.publicVirtualMethodTable[index2], pCL.AbsA, vm, true, false)
		} else {
			if reflect.DeepEqual(pPI.AID, frameworkAID) { //classes of package framework
				classtoken = pCI.info[1]
				token = pCI.info[2]
				//call appropriate virtual method
				callFrameworkMethods(vm, classtoken, token)
			} else {
				fmt.Println("Error: cannot invoke virtual package not found")
			}

		}

	} else { //Interal class library
		offset := makeU2(pCI.info[0], pCI.info[1])
		token := pCI.info[2]
		pcLInf := pCA.PClass.pClasses[offset-2]
		newFrame := &Frame{}
		vm.PushFrame(newFrame)
		if token&0x80 == 0x80 { //call private method
			index2 := token - pcLInf.packageMethodTableBase
			pCA.PMethod.executeByteCode(pcLInf.packageVirtualMethodTable[index2], pCA, vm, true, false)
		} else { //call public method
			index2 := token - pcLInf.publicMethodTableBase
			pCA.PMethod.executeByteCode(pcLInf.publicVirtualMethodTable[index2], pCA, vm, true, false)
		}
	}

}
func invokespecial(currF *Frame, index uint16, pCA *AbstractApplet, vm *VM) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	if pCI.tag == 0x06 { //static method
		if byte1 == 0x00 { //package internal method
			sOffset := makeU2(pCI.info[1], pCI.info[2])
			newFrame := &Frame{}
			vm.PushFrame(newFrame)
			pCA.PMethod.executeByteCode(sOffset, pCA, vm, true, false)
		} else {
			packageIndex := byte1 & 0x7F
			pPI := pCA.PImport.packages[packageIndex]
			pCL := findLibrary(pPI)
			if pCL != nil {
				//External library, call external static method
				classtoken := pCI.info[1]
				sOffset := pCL.AbsA.PExport.pClassExport[classtoken].pStaticMethodOffsets[pCI.info[2]]
				newFrame := &Frame{}
				vm.PushFrame(newFrame)
				pCL.AbsA.PMethod.executeByteCode(sOffset, pCL.AbsA, vm, true, false)
			} else {
				if reflect.DeepEqual(pPI.AID, frameworkAID) { //classes of package framework
					classtoken := pCI.info[1]
					token := pCI.info[2]
					//call appropriate virtual method
					callFrameworkMethods(vm, classtoken, token)
				} else {
					fmt.Println("Error: cannot invoke special package not found")
				}

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
	newFrame.Localvariables = make([]interface{}, 100)
	newFrame.operandStack = make([]interface{}, 100)
	for i := nargs; i > 0; i-- {
		newFrame.Localvariables[i-1] = currF.pop()
	}
	vm.PushFrame(newFrame)
	objref := newFrame.Localvariables[0].(Reference) //the object reference
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

func vmNew(currF *Frame, index uint16, pCA *AbstractApplet) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	var class *JavaClass
	if byte1&0x80 == 0x80 {
		packageIndex := byte1 & 0x7F
		pPI := pCA.PImport.packages[packageIndex]
		pCL := findLibrary(pPI)
		if pCL != nil {
			//External library
			classtoken := pCI.info[1]
			sOffset := pCL.AbsA.PExport.pClassExport[classtoken].classOffset
			class = createClassInstance(classtoken, sOffset, pCL.AbsA)
		} else {
			if reflect.DeepEqual(pPI.AID, frameworkAID) { //classes of package framework
				classtoken := pCI.info[1]
				createFrameworkClass(currF, classtoken)
			}
			fmt.Println("Error: cannot create class package not found")
		}
	} else {

		offset := makeU2(pCI.info[0], pCI.info[1])
		token := pCI.info[2]
		//	pcLInf := pCA.PClass.pClasses[offset]
		class = createClassInstance(token, offset, pCA)
	}
	jcCount++
	//add in the map between aid and instance references
	pckg := pCA.PHeader.PThisPackage
	aid := api.InitAID(pckg.AID, 0, int16(pckg.AIDLength))
	InstanceRefHeap[aid] = Reference(jcCount)
	fmt.Println(InstanceRefHeap[aid])

	//add instance in memory heap
	heap[Reference(jcCount)] = class
	currF.push(Reference(jcCount)) //arbritrary number
}
func createClassInstance(classtoken uint8, soffset uint16, pCL *AbstractApplet) *JavaClass {
	//javaClass := &JavaClass{}
	superclassref := pCL.PClass.pClasses[soffset-2].superClassRef
	declaredinstancesize := pCL.PClass.pClasses[soffset-2].declaredInstanceSize
	var classInf *ClassDescriptorInfo
	var classref uint16
	for i := 0; i < int(pCL.PDescriptor.classCount); i++ {
		if pCL.PDescriptor.classes[i].token == classtoken {
			classref = pCL.PDescriptor.classes[i].thisClassRef
			classInf = pCL.PDescriptor.classes[i]
			break
		}
	}
	javaClass := &JavaClass{classref, superclassref, declaredinstancesize, nil, nil}
	setInstanceFieldsDefaultValue(classInf, javaClass)
	return javaClass
}
func setInstanceFieldsDefaultValue(classInf *ClassDescriptorInfo, javaClass *JavaClass) {
	var token uint8
	var value interface{}
	javaClass.fields = make([]*instanceField, classInf.fieldCount)
	for i := 0; i < int(classInf.fieldCount); i++ {
		if classInf.fields[i].pAF&AccStatic != AccStatic {
			a := classInf.fields[i].pFieldtype
			switch a {
			case 0x8002:
				//bool
				token = classInf.fields[i].token
				value = int16(0)
			case 0x8003:
				//byte
				token = classInf.fields[i].token
				value = int16(0)

			case 0x8004:
				//short
				token = classInf.fields[i].token
				value = int16(0)

			case 0x8005:
				//int
				token = classInf.fields[i].token
				value = int32(0)
			default:
				//reference type
				token = classInf.fields[i].token
				value = Reference(0)
			}
			javaClass.fields[i] = &instanceField{token, value}
		}
	}
}
func newArray(currF *Frame, atype uint8) {
	count := currF.pop().(int16)
	array := &ArrayValue{}
	switch atype {
	case 10:
		//boolean
		array.componentType = TypeBoolean
		array.length = uint16(count)
		array.array = make([]uint8, count)
	case 11:
		//byte
		array.componentType = TypeByte
		array.length = uint16(count)
		array.array = make([]byte, count)
	case 12:
		//short
		array.componentType = TypeShort
		array.length = uint16(count)
		array.array = make([]int16, count)
	case 13:
		//int
		array.componentType = TypeInt
		array.length = uint16(count)
		array.array = make([]int32, count)
	}
	heap[Reference(arrcount+1)] = array
	currF.push(Reference(arrcount + 1))
}
func anewArray(currF *Frame, index uint16, pCA *AbstractApplet) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	//var class *JavaClass
	if byte1&0x80 == 0x80 {
		//trouver si c'est interface  ou class ensuite creer sur heap
	} else {
		//trouver si c'est interface  ou class ensuite creer sur heap
	}

}
func getFieldThis(currF *Frame, index uint8, pCA *AbstractApplet) {
	pCI := pCA.PConstPool.pConstantPool[index]
	instanceref := currF.Localvariables[0].(Reference)
	class := heap[Reference(instanceref)]
	switch class.(type) {
	case *JavaClass:
		for i := 0; i < len(class.(*JavaClass).fields); i++ {
			if class.(*JavaClass).fields[i].token == pCI.info[2] {
				value := class.(*JavaClass).fields[i].value
				switch value.(type) {
				case uint8:
					currF.push(int16(value.(uint8)))
				case int16:
					currF.push(value.(int16))
				default:
					currF.push(value)
				}
				break
			}
		}
	}
}
func getField(currF *Frame, index uint8, pCA *AbstractApplet) {
	pCI := pCA.PConstPool.pConstantPool[index]
	instanceref := currF.pop().(Reference)
	class := heap[Reference(instanceref)]
	switch class.(type) {
	case *JavaClass:
		for i := 0; i < len(class.(*JavaClass).fields); i++ {
			if class.(*JavaClass).fields[i].token == pCI.info[2] {
				value := class.(*JavaClass).fields[i].value
				switch value.(type) {
				case uint8:
					currF.push(int16(value.(uint8)))
				case int16:
					currF.push(value.(int16))
				default:
					currF.push(value)
				}
				break
			}
		}
	}
}

func ifScmpne(currF *Frame, branch int8, pPC *int) {
	value2 := currF.pop()
	value1 := currF.pop()
	if value1.(int16) != value2.(int16) {
		(*pPC) += int(branch)
		(*pPC) -= 2
	}
}
func putfield(currF *Frame, index uint8, pCA *AbstractApplet) {
	svalue := currF.pop()
	objref := currF.pop()
	pCI := pCA.PConstPool.pConstantPool[index]
	class := heap[objref.(Reference)]
	switch class.(type) {
	case *JavaClass:
		for i := 0; i < len(class.(*JavaClass).fields); i++ {
			token := class.(*JavaClass).fields[i].token
			if token == pCI.info[2] {
				class.(*JavaClass).fields[i].value = svalue
				break
			}
		}
	}
}
func athrow(currF *Frame) { //TODO
	status = uint16(currF.Localvariables[1].(int16))
	return
}
func getstatic(currF *Frame, index uint16, pCA *AbstractApplet, ins int) {
	offset, err := getstaticfieldAddress(index, pCA)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch ins {
	case 0x7B:
		value := readU2(pCA.PStaticField.pStaticFieldImage, &offset)
		currF.push(Reference(value))
	case 0x7C:
		value := pCA.PStaticField.pStaticFieldImage[offset]
		currF.push(int16(value))
	case 0x7D:
		value := readU2(pCA.PStaticField.pStaticFieldImage, &offset)
		currF.push(int16(value))
	case 0x7E:
		value := readU4(pCA.PStaticField.pStaticFieldImage, &offset)
		currF.push(int32(value))
	}
}

func getstaticfieldAddress(index uint16, pCA *AbstractApplet) (int, error) {
	pCI := pCA.PConstPool.pConstantPool[index]
	byte1 := pCI.info[0]
	offset := makeU2(pCI.info[1], pCI.info[2])
	if byte1 == 0 {
		//It is an internal static field.
		return int(offset), nil

	}
	return 0, errors.New("Error getstatic on external address")
}
func putstatic(currF *Frame, index uint16, pCA *AbstractApplet, ins int) {
	offset, err := getstaticfieldAddress(index, pCA)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch ins {
	case 0x7F:
		value := currF.pop().(Reference)
		pCA.PStaticField.pStaticFieldImage[offset] = uint8(value / 0x100)
		pCA.PStaticField.pStaticFieldImage[offset+1] = uint8(value % 0x100)
	case 0x80:
		value := currF.pop().(uint8)
		pCA.PStaticField.pStaticFieldImage[offset] = value
	case 0x81:
		value := currF.pop().(int16)
		pCA.PStaticField.pStaticFieldImage[offset] = uint8(value / 0x100)
		pCA.PStaticField.pStaticFieldImage[offset+1] = uint8(value % 0x100)
	case 0x82:
		//TODO
		/*value := currF.pop().(int16)
		pCA.PStaticField.pStaticFieldImage[offset] = uint8(value / 0x100)
		pCA.PStaticField.pStaticFieldImage[offset+1] = uint8(value % 0x100)
		pCA.PStaticField.pStaticFieldImage[offset+2] = uint8(value / 0x100)
		pCA.PStaticField.pStaticFieldImage[offset+3] = uint8(value % 0x100)*/
	}
}
func callselectingApplet(currF *Frame) {
	_ = currF.pop()
	cond := api.GetSelectingAppletFlag()
	if cond {
		currF.push(int16(1)) //push true
	} else {
		currF.push(int16(0)) //push false
	}
}
func callArraycopynonAtomic(currF *Frame) {
	length := currF.pop().(int16)
	destOffset := currF.pop().(int16)
	destArr := heap[currF.pop().(Reference)].(*ArrayValue)
	srcOffset := currF.pop().(int16)
	srcArr := heap[currF.pop().(Reference)].(*ArrayValue)
	result := api.ArrayCopyNonAtomic(srcArr.array, srcOffset, destArr.array, destOffset, length)
	currF.push(result)
}
func callOwnerPINInit(currF *Frame) {
	maxPINSize := currF.pop().(int16)
	tryLimit := currF.pop().(int16)
	ownerPinRef := currF.pop().(Reference)
	ownerPIN := heap[ownerPinRef].(*api.OwnerPIN)
	ownerPIN.InitOwnerPIN(byte(tryLimit), byte(maxPINSize))
}
func callOwnerPINCheck(currF *Frame) {
	length := currF.pop().(int16)
	offset := currF.pop().(int16)
	arr := heap[currF.pop().(Reference)].(*ArrayValue)
	ownerPinRef := currF.pop().(Reference)
	ownerPIN := heap[ownerPinRef].(*api.OwnerPIN)
	cond := ownerPIN.Check(arr.array, offset, byte(length))
	if cond {
		currF.push(int16(1)) //push true
	} else {
		currF.push(int16(0)) //push false
	}
}
func callIsValidated(currF *Frame) {
	ownerPinRef := currF.pop().(Reference)
	ownerPIN := heap[ownerPinRef].(*api.OwnerPIN)
	cond := ownerPIN.IsValidated()
	if cond {
		currF.push(int16(1)) //push true
	} else {
		currF.push(int16(0)) //push false
	}
}
func callResetAndUnblock(currF *Frame) {
	ownerPinRef := currF.pop().(Reference)
	ownerPIN := heap[ownerPinRef].(*api.OwnerPIN)
	ownerPIN.ResetAndUnblock()
}
func callUpdateOwnerPIN(currF *Frame) {
	length := currF.pop().(int16)
	offset := currF.pop().(int16)
	arr := heap[currF.pop().(Reference)].(*ArrayValue)
	ownerPinRef := currF.pop().(Reference)
	ownerPIN := heap[ownerPinRef].(*api.OwnerPIN)
	n, err := ownerPIN.Update(arr.array, offset, byte(length))
	if n == 1 && err != nil {
		SetStatus(n)
		fmt.Println(err)
		leaveVM = true
	}
}
func callRegister(currF *Frame) {
	ref := currF.pop().(Reference)
	for i, j := range InstanceRefHeap {
		if reflect.DeepEqual(j, ref) {
			AppletTable[i] = ConstantApplet
			return
		}
	}
}
func callFrameworkMethods(vm *VM, classtoken uint8, methodToken uint8) {

	currF := vm.StackFrame[vm.FrameTop]
	switch classtoken {
	case applet:
		if methodToken == register {
			callRegister(currF)
		} else if methodToken == selectingApplet {
			callselectingApplet(currF)

		}
	case util:
		if methodToken == arraycopynonAtomic {
			callArraycopynonAtomic(currF)
		}
	case apdu:
		switch methodToken {
		case getbuffer:
			callgetBuffer(currF)
		case receiveBytes:
			callreceiveBytes(currF)
		case sendBytes:
			length := currF.pop().(int16)
			offset := currF.pop().(int16)
			_ = currF.pop()
			sndBytes(offset, length)
		case sendBytesLong:
			callsendBytesLong(currF)
		case setIncomingAndReceive:
			callsetIncomingandreceive(currF)
		case setOutgoing:
			callsetOutgoing(currF)
		case setOutgoingAndSend:
			callsetOutgoingAndSend(currF)
		case setOutgoingLength:
			callsetOutgoingLength(currF)
		default:
			//nothing
			return
		}

	case ownerpinClass:
		switch methodToken {
		case 0: //init
			callOwnerPINInit(currF)
		case check:
			callOwnerPINCheck(currF)
		case isValidated:
			callIsValidated(currF)
		case resetAndUnblock:
			callResetAndUnblock(currF)
		case update:
			callUpdateOwnerPIN(currF)
		}
	case isoException:
		if methodToken == throwit {
			callthrowit(currF)

		}
	}

}
func callgetBuffer(currF *Frame) {
	arrayref := currF.pop().(Reference) //6000
	arr := heap[arrayref].(*ArrayValue)
	apduarray := arr.array.([]byte)
	copy(apduarray[0:], nativeMethods.BufferRcv[:len(apduarray)])
	currF.push(Reference(6000))
}
func callthrowit(currF *Frame) {
	status := currF.pop().(int16)
	SetStatus(uint16(status))
	leaveVM = true
}
func callreceiveBytes(currF *Frame) {
	//Only APDUs case 3 and 4 are expected to call this method.
	offset := currF.pop().(int16)
	ref := currF.pop().(Reference)
	arr := heap[ref].(*ArrayValue)
	length := nativeMethods.T0RcvData(arr.array.([]byte), offset)
	currF.push(length)
}
func callsetIncomingandreceive(currF *Frame) {
	ref := currF.pop().(Reference)
	arr := heap[ref].(*ArrayValue)
	length := nativeMethods.T0RcvData(arr.array.([]byte), int16(api.OffsetCData))
	currF.push(length)
}
func callsetOutgoing(currF *Frame) {
	_ = currF.pop().(Reference)
	currF.push(nativeMethods.LE)
}
func callsetOutgoingLength(currF *Frame) {
	Len := currF.pop().(int16)
	_ = currF.pop()
	nativeMethods.LR = byte(Len)
}
func sndBytes(offset, length int16) {
	for length > 0 {
		temp := length
		// Need to force GET RESPONSE for Case 4 & for partial blocks
		if length != int16(nativeMethods.LR) || nativeMethods.LR != nativeMethods.LE || nativeMethods.SendInProgressFlag {
			temp = nativeMethods.Send61xx(length) // resets
			nativeMethods.LR -= byte(temp)
			nativeMethods.LE = nativeMethods.LR
		}
		arr := heap[Reference(6000)].(*ArrayValue)
		nativeMethods.T0SendData(arr.array.([]byte), offset, temp)
		nativeMethods.SendInProgressFlag = true
		offset += temp
		length -= temp
	}

}
func callsendBytesLong(currF *Frame) {
	Len := currF.pop().(int16)
	bOff := currF.pop().(int16)
	arrayref := currF.pop().(Reference)
	apduref := currF.pop().(Reference)
	outData := heap[arrayref].(*ArrayValue).array.([]byte)
	apduBuff := heap[apduref].(*ArrayValue).array.([]byte)
	api.CheckArrayArgs(outData, bOff, Len)
	sendLength := int16(len(apduBuff))
	for Len > 0 {
		if Len < sendLength {
			sendLength = Len
		}
		api.ArrayCopy(outData, bOff, apduBuff, 0, sendLength)
		sndBytes(0, sendLength)
		Len -= sendLength
		bOff += sendLength
	}
}
func callsetOutgoingAndSend(currF *Frame) {
	Len := currF.pop().(int16)
	bOff := currF.pop().(int16)
	_ = currF.pop().(Reference)
	nativeMethods.LR = byte(Len)
	sndBytes(bOff, Len)
}
func createFrameworkClass(currF *Frame, classtoken uint8) {
	jcCount++
	switch classtoken {
	case ownerpinClass: //OwnerPIN
		ownerPIN := &api.OwnerPIN{}
		//add instance in memory heap
		heap[Reference(jcCount)] = ownerPIN
		currF.push(Reference(jcCount)) //arbritrary number
	}
}
