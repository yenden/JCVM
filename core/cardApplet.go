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
	minV := abs.PHeader.pThisPackage.MinorVersion
	majV := abs.PHeader.pThisPackage.MajorVersion
	aidL := abs.PHeader.pThisPackage.AIDLength
	if minV == pPI.MinorVersion && majV == pPI.MajorVersion && aidL == pPI.AIDLength {
		for i := 0; i < int(aidL); i++ {
			if abs.PHeader.pThisPackage.AID[i] != pPI.AID[i] {
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
