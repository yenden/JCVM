package core

const (
	maxlocals  = 256
	maxOpStack = 256
	maxFrame   = 256
)

/*ReturnAddress type*/
type ReturnAddress int

/*Reference type*/
type Reference uint16

/*NullType type*/
type NullType int

/*Frame represent an executing method*/
type Frame struct {
	Localvariables []interface{}
	operandStack   []interface{}
	opStackTop     int
}

/*VM : The main VM structure
--- Vm is a set of Frames
*/
type VM struct {
	StackFrame []*Frame //creer apres maxframe
	FrameTop   int
}

var (
	status           = uint16(0x9000)
	leaveVM          = false
	slookupswitchMap = make(map[int16]int16)
)

/*GetStatus helps the jcre to get the status modified by JCVM*/
func GetStatus() uint16 {
	return status
}

/*SetStatus is used when the VM modifies SW*/
func SetStatus(st uint16) {
	status = st
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

/*PushFrame pushes frame in the VM*/
func (vm *VM) PushFrame(frame *Frame) bool {
	if vm.FrameTop+1 == maxFrame {
		return false
	}
	vm.FrameTop++
	vm.StackFrame[vm.FrameTop] = frame
	return true

}

/****Pop functions****/
func (frame *Frame) pop() interface{} {
	if frame.opStackTop == -1 {
		return 0
	}
	val := frame.operandStack[frame.opStackTop]
	i := frame.opStackTop
	frame.operandStack = append(frame.operandStack[:i], frame.operandStack[i+1:]...)
	frame.opStackTop--
	return val

}

func (vm *VM) popFrame() *Frame {
	if vm.FrameTop == -1 {
		return nil
	}
	fr := vm.StackFrame[vm.FrameTop]
	i := vm.FrameTop
	vm.StackFrame = append(vm.StackFrame[:i], vm.StackFrame[i+1:]...)
	vm.FrameTop--
	return fr

}

func (vm *VM) runStatic(pByteCode []uint8, pPC *int, pCA *AbstractApplet, params uint8) {
	currentFrame := vm.StackFrame[vm.FrameTop]
	leaveVM = false
	for !leaveVM {
		bytecode := int(readU1(pByteCode, pPC))
		switch bytecode {
		case 0x0:
			// no operation
		case 0x01:
			//aconstnull
			aconstNull(currentFrame)
		case 0x02:
			sconst(currentFrame, -1)
		case 0x03:
			sconst(currentFrame, 0)
		case 0x04:
			sconst(currentFrame, 1)
		case 0x05:
			sconst(currentFrame, 2)
		case 0x06:
			sconst(currentFrame, 3)
		case 0x07:
			sconst(currentFrame, 4)
		case 0x08:
			sconst(currentFrame, 5)
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
		case 0x10:
			byte1 := readU1(pByteCode, pPC)
			bspush(currentFrame, byte1)
		case 0x11:
			byte1 := readU2(pByteCode, pPC)
			sspush(currentFrame, byte1)
		case 0x12:
			byte1 := readU1(pByteCode, pPC)
			bipush(currentFrame, byte1)
		case 0x13:
			sValue := readS2(pByteCode, pPC)
			sipush(currentFrame, sValue)
		case 0x15:
			bIndex := readU1(pByteCode, pPC)
			aload(currentFrame, bIndex)
		case 0x16:
			bIndex := readU1(pByteCode, pPC)
			sload(currentFrame, bIndex)
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
		case 0x1C:
			sload(currentFrame, 0) //sload_0
		case 0x1D:
			sload(currentFrame, 1) //sload_1
		case 0x1E:
			sload(currentFrame, 2) //sload_2
		case 0x1F:
			sload(currentFrame, 3) //sload_3
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
		case 0x29:
			bValue := readU1(pByteCode, pPC)
			sstore(currentFrame, bValue)
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
		case 0x2F:
			sstore(currentFrame, 0)
		case 0x30:
			sstore(currentFrame, 1)
		case 0x31:
			sstore(currentFrame, 2)
		case 0x32:
			sstore(currentFrame, 3)
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
			popBytecode(currentFrame)
		case 0x3D:
			dup(currentFrame)
		case 0x3E:
			dup2(currentFrame)
		case 0x3F:
			mn := readU1(pByteCode, pPC)
			dupX(currentFrame, mn)
		case 0x41:
			sadd(currentFrame)
		case 0x42:
			iadd(currentFrame)
		case 0x43:
			ssub(currentFrame)
		case 0x44:
			isub(currentFrame)
		case 0x45:
			smul(currentFrame)
		case 0x46:
			imul(currentFrame)
		case 0x48:
			idiv(currentFrame)
		case 0x49:
			srem(currentFrame)
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
		case 0x59:
			index := readU1(pByteCode, pPC)
			constant := readS1(pByteCode, pPC)
			sinc(currentFrame, index, constant)
		case 0x5A:
			index := readU1(pByteCode, pPC)
			constant := readS1(pByteCode, pPC)
			iinc(currentFrame, index, constant)
		case 0x5B:
			s2b(currentFrame)
		case 0x5C:
			s2i(currentFrame)
		case 0x5D:
			i2b(currentFrame)
		case 0x5E:
			i2s(currentFrame)
		case 0x5F:
			icmp(currentFrame)
		case 0x60:
			bValue := readS1(pByteCode, pPC)
			ifeq(currentFrame, bValue, pPC)
		case 0x61:
			bValue := readS1(pByteCode, pPC)
			ifne(currentFrame, bValue, pPC)
		case 0x62:
			bValue := readS1(pByteCode, pPC)
			iflt(currentFrame, bValue, pPC)
		case 0x63:
			bValue := readS1(pByteCode, pPC)
			ifge(currentFrame, bValue, pPC)
		case 0x64:
			bValue := readS1(pByteCode, pPC)
			ifgt(currentFrame, bValue, pPC)
		case 0x65:
			bValue := readS1(pByteCode, pPC)
			ifle(currentFrame, bValue, pPC)
		case 0x66:
			bValue := readS1(pByteCode, pPC)
			ifnull(currentFrame, bValue, pPC)
		case 0x67:
			bValue := readS1(pByteCode, pPC)
			ifnnnull(currentFrame, bValue, pPC)
		case 0x6A:
			bValue := readS1(pByteCode, pPC)
			ifScmpeq(currentFrame, bValue, pPC)
		case 0x6B:
			bValue := readS1(pByteCode, pPC)
			ifScmpne(currentFrame, bValue, pPC)
		case 0x70:
			bValue := readS1(pByteCode, pPC)
			goTo(currentFrame, bValue, pPC)
		case 0x75:
			deflt := readS2(pByteCode, pPC)
			npairs := readU2(pByteCode, pPC)
			for k := 0; k < int(npairs); k++ {
				match := readS2(pByteCode, pPC)
				pairs := readS2(pByteCode, pPC)
				slookupswitchMap[match] = pairs
			}
			slookupswitch(currentFrame, deflt, npairs, pPC)
		case 0x77:
			invokerFrame := vm.StackFrame[vm.FrameTop-1]
			areturn(currentFrame, invokerFrame)
			vm.popFrame()
			return
		case 0x78:
			invokerFrame := vm.StackFrame[vm.FrameTop-1]
			sreturn(currentFrame, invokerFrame)
			vm.popFrame()
			return
		case 0x79:
			invokerFrame := vm.StackFrame[vm.FrameTop-1]
			ireturn(currentFrame, invokerFrame)
			vm.popFrame()
			return
		case 0x7A:
			vm.popFrame()
			return

		case 0x7B: //getstatic_a
			value := readU2(pByteCode, pPC)
			getstatic(currentFrame, value, pCA, 0x7B)
		case 0x7C: //getstatic_b
			value := readU2(pByteCode, pPC)
			getstatic(currentFrame, value, pCA, 0x7C)
		case 0x7D: //getstatic_s
			value := readU2(pByteCode, pPC)
			getstatic(currentFrame, value, pCA, 0x7D)
		case 0x7E: //getstatic_i
			value := readU2(pByteCode, pPC)
			getstatic(currentFrame, value, pCA, 0x7E)

		case 0x7F: //putstatic_a
			value := readU2(pByteCode, pPC)
			putstatic(currentFrame, value, pCA, 0x7F)
		case 0x80: //putstatic_b
			value := readU2(pByteCode, pPC)
			putstatic(currentFrame, value, pCA, 0x80)
		case 0x81: //putstatic_s
			value := readU2(pByteCode, pPC)
			putstatic(currentFrame, value, pCA, 0x81)
		case 0x82: //puttstatic_i
			value := readU2(pByteCode, pPC)
			putstatic(currentFrame, value, pCA, 0x82)

		case 0x87:
			index7 := readU1(pByteCode, pPC)
			putfield(currentFrame, index7, pCA)
		case 0x88:
			index8 := readU1(pByteCode, pPC)
			putfield(currentFrame, index8, pCA)
		case 0x89:
			index9 := readU1(pByteCode, pPC)
			putfield(currentFrame, index9, pCA)
		case 0x8a:
			indexa := readU1(pByteCode, pPC)
			putfield(currentFrame, indexa, pCA)

		case 0x8B:
			indexb := readU2(pByteCode, pPC)
			invokevirtual(currentFrame, indexb, pCA, vm)
		case 0x8C:
			indexc := readU2(pByteCode, pPC)
			invokespecial(currentFrame, indexc, pCA, vm)
		case 0x8D:
			indexd := readU2(pByteCode, pPC)
			invokestatic(currentFrame, indexd, pCA, vm)
		case 0x8E:
			nargs := readU1(pByteCode, pPC)
			indexe := readU2(pByteCode, pPC)
			methodToken := readU1(pByteCode, pPC)
			invokeinterface(currentFrame, pCA, vm, nargs, indexe, methodToken)
		case 0x8F:
			index2 := readU2(pByteCode, pPC)
			vmNew(currentFrame, index2, pCA)
		case 0x90:
			atype := readU1(pByteCode, pPC)
			newArray(currentFrame, atype)
		case 0x93:
			athrow(currentFrame)
			return
		case 0xAD:
			indexad := readU1(pByteCode, pPC)
			getFieldThis(currentFrame, indexad, pCA)
		case 0xAE:
			indexae := readU1(pByteCode, pPC)
			getFieldThis(currentFrame, indexae, pCA)
		case 0xAF:
			indexaf := readU1(pByteCode, pPC)
			getFieldThis(currentFrame, indexaf, pCA)

		case 0xB0:
			indexb0 := readU1(pByteCode, pPC)
			getFieldThis(currentFrame, indexb0, pCA)
		}

	}
}
