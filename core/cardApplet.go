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

/*
func (abs *AbstractApplet) isThisLibrary(pPI *PackageInfo) bool {
	return (*(abs.pHeader.pThisPackage) == *pPI)
}*/

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
	/*cl.absA.pHeader = nil
	cl.absA.pDir = nil
	cl.absA.pClass = nil
	cl.absA.pConstPool = nil
	cl.absA.pDescriptor = nil
	cl.absA.pImport = nil
	cl.absA.pMethod = nil
	cl.absA.pRefLoc = nil
	cl.absA.pStaticField = nil
	cl.absA.pExport = nil*/

	return pcl
}

func (cl *CardApplet) install() {
	fmt.Println("Start installing...")
	if cl.PApplet == nil {
		fmt.Println("Not an applet!")
		return
	}
	fmt.Printf("Install command from %d", int(cl.PApplet.applets[0].installMethodOffset))
	//cl.AbsA.PMethod.executeByteCode(uint16(cl.PApplet.applets[0].installMethodOffset), cl.AbsA)
	fmt.Println("Install finished!")
}

func (cl *CardApplet) process() {
}
