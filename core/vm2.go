package core

const (
	maxlocals  = 256
	maxOpStack = 256
	maxFrame   = 256
)

type ReturnAddress int
type Reference int16
type NullType int
type Frame struct {
	localvariables [maxlocals]interface{}
	operandStack   [maxOpStack]interface{}
	opStackTop     int
}
type VM struct {
	stackFrame [maxFrame]*Frame
	frameTop   int
}

/****Push functions****/

func (frame *Frame) push(value interface{}) bool {
	if frame.opStackTop+1 == maxOpStack {
		return false
	}
	frame.opStackTop++
	frame.operandStack[frame.opStackTop] = value
	return true

}

func (vm *VM) pushFrame(frame *Frame) bool {
	if vm.frameTop+1 == maxFrame {
		return false
	}
	vm.frameTop++
	vm.stackFrame[vm.frameTop] = frame
	return true

}

/****Pop functions****/
func (frame *Frame) pop() interface{} {
	if frame.opStackTop == -1 {
		return 0
	}
	val := frame.operandStack[frame.opStackTop]
	frame.opStackTop--
	return val

}

func (vm *VM) popFrame() *Frame {
	if vm.frameTop == -1 {
		return nil
	}
	fr := vm.stackFrame[vm.frameTop]
	vm.frameTop--
	return fr

}

func (vm *VM) runStatic(pByteCode []uint8, pPC *int, pCA *AbstractApplet, maxStack uint8, params uint8, maxLocals uint8) {
	currentFrame := vm.stackFrame[vm.frameTop]
	if vm.frameTop >= 1 {
		invokerFrame := vm.stackFrame[vm.frameTop-1]
		for i := params; i > 0; i-- {
			currentFrame.localvariables[i-1] = invokerFrame.pop()
		}
	}

	for {
		bytecode := int(readU1(pByteCode, pPC))
		switch bytecode {
		case 0x0:
			// no operation
		case 0x01:
			//aconstnull
			aconstNull(currentFrame)
		case 0x09:
			iconst(currentFrame, -1)
		case 0x0A:
			iconst(currentFrame, 0)
		case 0x0B:
			iconst(currentFrame, 1)
		case 0x0C:
			iconst(currentFrame, 2)
		case 0x0D:
			iconst(currentFrame, 3)
		case 0x0E:
			iconst(currentFrame, 4)
		case 0x0F:
			iconst(currentFrame, 5)
		case 0x12:
			byte1 := readU1(pByteCode, pPC)
			bipush(currentFrame, byte1)
		case 0x13:
			sValue := readS2(pByteCode, pPC)
			sipush(currentFrame, sValue)
		case 0x15:
			bIndex := readU1(pByteCode, pPC)
			aload(currentFrame, bIndex)
		case 0x17:
			bIndex := readU1(pByteCode, pPC)
			iload(currentFrame, bIndex)
		case 0x18:
			aload(currentFrame, 0) //aload_0
		case 0x19:
			aload(currentFrame, 1) //aload_1
		case 0x1A:
			aload(currentFrame, 2) //aload_2
		case 0x1B:
			aload(currentFrame, 3) //aload_3
		case 0x20:
			iload(currentFrame, 0) //iload_0
		case 0x21:
			iload(currentFrame, 1) //iload_1
		case 0x22:
			iload(currentFrame, 2) //iload_2
		case 0x23:
			iload(currentFrame, 3) //iload_3
		case 0x24:
			aaload(currentFrame)
		case 0x25:
			baload(currentFrame)
		case 0x26:
			saload(currentFrame)
		case 0x28:
			bValue := readU1(pByteCode, pPC)
			astore(currentFrame, bValue)
		case 0x2A:
			bIndex := readU1(pByteCode, pPC)
			istore(currentFrame, bIndex)
		case 0x2B:
			astore(currentFrame, 0) //astore_0
		case 0x2C:
			astore(currentFrame, 1) //astore_1
		case 0x2D:
			astore(currentFrame, 2) //astore_2
		case 0x2E:
			astore(currentFrame, 3) //astore_3
		case 0x33:
			istore(currentFrame, 0) //istore_0
		case 0x34:
			istore(currentFrame, 1) //istore_1
		case 0x35:
			istore(currentFrame, 2) //istore_2
		case 0x36:
			istore(currentFrame, 3) //istore_3
		case 0x37:
			aastore(currentFrame)
		case 0x38:
			bastore(currentFrame)
		case 0x39:
			sastore(currentFrame)
		case 0x3B:
			popVal := popBytecode(currentFrame)
		case 0x3D:
			dup(currentFrame)
		case 0x3E:
			dup2(currentFrame)
		case 0x3F:
			dupX(currentFrame) //todo
		case 0x42:
			iadd(currentFrame)
		case 0x44:
			isub(currentFrame)
		case 0x46:
			imul(currentFrame)
		case 0x48:
			idiv(currentFrame)
		case 0x4a:
			irem(currentFrame)
		case 0x4E:
			ishl(currentFrame)
		case 0x52:
			iushr(currentFrame)
		case 0x54:
			iand(currentFrame)
		case 0x56:
			ior(currentFrame)
		case 0x58:
			ixor(currentFrame)
		case 0x5D:
			i2b(currentFrame)
		case 0x5E:
			i2s(currentFrame)
		case 0x60:
			bValue := int8(readU1(pByteCode, pPC))
			ifeq(currentFrame, bValue, pPC)
		case 0x61:
			bValue := int8(readU1(pByteCode, pPC))
			ifne(currentFrame, bValue, pPC)
		case 0x62:
			bValue := int8(readU1(pByteCode, pPC))
			iflt(currentFrame, bValue, pPC)
		case 0x63:
			bValue := int8(readU1(pByteCode, pPC))
			ifge(currentFrame, bValue, pPC)
		case 0x64:
			bValue := int8(readU1(pByteCode, pPC))
			ifgt(currentFrame, bValue, pPC)
		case 0x65:
			bValue := int8(readU1(pByteCode, pPC))
			ifle(currentFrame, bValue, pPC)
		case 0x66:
			bValue := int8(readU1(pByteCode, pPC))
			ifnull(currentFrame, bValue, pPC)
		case 0x67:
			bValue := int8(readU1(pByteCode, pPC))
			ifnnnull(currentFrame, bValue, pPC)
		case 0x70:
			bValue := int8(readU1(pByteCode, pPC))
			goTo(currentFrame, bValue, pPC)
		case 0x77:
			invokerFrame := vm.stackFrame[vm.frameTop-1]
			areturn(currentFrame, invokerFrame)
			vm.frameTop--
			return
		case 0x79:
			invokerFrame := vm.stackFrame[vm.frameTop-1]
			ireturn(currentFrame, invokerFrame)
			vm.frameTop--
			return
		case 0x7A:
			vm.frameTop--
			return
		case 0x8B:
			index := readU2(pByteCode, pPC)
			invokevirtual(currentFrame, index, pCA, vm)
		case 0x8C:
		}

	}
}