package core

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
	currF.push(val)
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
		currF.localvariables[index] = ref
	case ReturnAddress:
		currF.localvariables[index] = ref
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
	result := value1.(int16) + value2.(int16)
	currF.push(result)
}
