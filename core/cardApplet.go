package core

import (
	"fmt"
)

type AbstractApplet struct {
	PHeader      *HeaderComponent
	PDir         *DirectoryComponent
	PImport      *ImportComponent
	PClass       *ClassComponent
	PStaticField *StaticFieldComponent
	PMethod      *MethodComponent
	PRefLoc      *ReferenceLocationComponent
	PConstPool   *ConstantPoolComponent
	PDescriptor  *DescriptorComponent
	PExport      *ExportComponent
}

func (abs *AbstractApplet) isThisLibrary(pPI *PackageInfo) bool {
	/*minV := abs.PHeader.pThisPackage.MinorVersion
	majV := abs.PHeader.pThisPackage.MajorVersion*/
	aidL := abs.PHeader.PThisPackage.AIDLength
	if /*minV == pPI.MinorVersion && majV == pPI.MajorVersion &&*/ aidL == pPI.AIDLength {
		for i := 0; i < int(aidL); i++ {
			if abs.PHeader.PThisPackage.AID[i] != pPI.AID[i] {
				return false
			}
		}
		return true
	}
	return false
}

type CardApplet struct {
	AbsA    *AbstractApplet
	PApplet *AppletComponent
}

func (cl *CardApplet) cloneLibrary() *AbstractApplet {
	pcl := &AbstractApplet{cl.AbsA.PHeader,
		cl.AbsA.PDir,
		cl.AbsA.PImport,
		cl.AbsA.PClass,
		cl.AbsA.PStaticField,
		cl.AbsA.PMethod,
		cl.AbsA.PRefLoc,
		cl.AbsA.PConstPool,
		cl.AbsA.PDescriptor,
		cl.AbsA.PExport,
	}
	cl.AbsA.PHeader = nil
	cl.AbsA.PDir = nil
	cl.AbsA.PClass = nil
	cl.AbsA.PConstPool = nil
	cl.AbsA.PDescriptor = nil
	cl.AbsA.PImport = nil
	cl.AbsA.PMethod = nil
	cl.AbsA.PRefLoc = nil
	cl.AbsA.PStaticField = nil
	cl.AbsA.PExport = nil

	return pcl
}

func (cl *CardApplet) Install(vm *VM) {
	fmt.Println("Start installing...")
	if cl.PApplet == nil {
		fmt.Println("Not an applet!")
		return
	}
	fmt.Printf("Install command from %d\r\n", int(cl.PApplet.applets[0].installMethodOffset))
	offset := uint16(cl.PApplet.applets[0].installMethodOffset)
	cl.AbsA.PMethod.executeByteCode(offset, cl.AbsA, vm, false)
	fmt.Println("Install finished!")
}

func (cl *CardApplet) Process(vm *VM) {
	fmt.Println("Start processing...")
	if cl.PApplet == nil {
		fmt.Println("Not an applet!")
		return
	}
	instanceclassref := vm.StackFrame[vm.FrameTop].Localvariables[0]
	instance := heap[instanceclassref.(Reference)]
	classref := instance.(*JavaClass).classref
	count := cl.AbsA.PDescriptor.classCount
	var processMethodOf uint16
	for i := 0; i < int(count); i++ {
		if cl.AbsA.PDescriptor.classes[i].thisClassRef == classref {
			class := cl.AbsA.PDescriptor.classes[i]
			for j := 0; j < int(class.methodCount); j++ {
				if class.methods[j].token == 7 {
					processMethodOf = class.methods[j].methodOffset
				}
			}
			break
		}
	}
	/*if processMethodOf != 7 {
		fmt.Println("Didn't found process method!")
		return
	}*/
	cl.AbsA.PMethod.executeByteCode(processMethodOf, cl.AbsA, vm, false)
	fmt.Println("Process finished!")
}
