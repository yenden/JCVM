package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ReadInBuffer copies IJC file to a byte slice
func ReadInBuffer(fileName string) []byte {
	fi, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}
	baComp := make([]byte, fi.Size())
	baComp, err = ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return baComp

}

//BuildApplet parse the ijc file and fill components
func BuildApplet(dataBuffer []byte, dataLength int) *CardApplet {

	var iPos int
	newapplet := &CardApplet{}
	var i uint8
	c1 := readU1(dataBuffer, &iPos)

	/* Read Header Component */
	if c1 != TagHeaderComp {
		log.Fatal("Header tag is not correct :This is not an applet!")
		return nil
	}

	iPos += 6
	appletMinVer := readU1(dataBuffer, &iPos)
	appletMajVer := readU1(dataBuffer, &iPos)
	appletFlags := readU1(dataBuffer, &iPos)

	minorVersion := readU1(dataBuffer, &iPos)
	majorVersion := readU1(dataBuffer, &iPos)
	aIDLength := readU1(dataBuffer, &iPos)
	pi := createPackageInfo(minorVersion, majorVersion, aIDLength)
	for i = 0; i < pi.AIDLength; i++ {
		pi.AID[i] = readU1(dataBuffer, &iPos)
	}
	nameLength := readU1(dataBuffer, &iPos)
	nameInfo := make([]uint8, nameLength)
	for i = 0; i < nameLength; i++ {
		nameInfo[i] = readU1(dataBuffer, &iPos)
	}

	ipHeader := &HeaderComponent{appletMinVer, appletMajVer, appletFlags, pi, &PackageNameInfo{nameLength, nameInfo}}

	/*Read Directory Component*/
	c1 = dataBuffer[iPos]
	if c1 != TagDirComp {
		log.Fatal("Directory tag is not  correct")
		return nil
	}
	iPos += 3
	var sizes [12]uint16
	for i = 0; i < 12; i++ {
		sizes[i] = readU2(dataBuffer, &iPos)
	}
	fmt.Println(sizes)
	pSfsi := &StaticFieldSizeInfo{}
	pSfsi.imageSize = readU2(dataBuffer, &iPos)
	pSfsi.arrayInitCount = readU2(dataBuffer, &iPos)
	pSfsi.arrayInitSize = readU2(dataBuffer, &iPos)
	importCount := readU1(dataBuffer, &iPos)
	appletCount := readU1(dataBuffer, &iPos)

	/*The subset of jcvm we consider doesn't have customize ClassComponent
	*We thus need to step over that data
	 */
	customcount := readU1(dataBuffer, &iPos)
	for i = 0; i < customcount; i++ {
		iPos = iPos + 3
		aidL := readU1(dataBuffer, &iPos)
		iPos = iPos + int(aidL)
	}

	ipDir := &DirectoryComponent{sizes, pSfsi, importCount, appletCount}

	/*Split the IJC file into component blocks
	*to facilitate component data manipulation.
	 */

	pAppletComponent := make([]uint8, sizes[TagAppletComp-1])
	pImportComponent := make([]uint8, sizes[TagImportComp-1])
	pConstantPoolComponent := make([]uint8, sizes[TagConstantPoolComp-1])
	pClassComponent := make([]uint8, sizes[TagClassComp-1])
	pMethodComponent := make([]uint8, sizes[TagMethodComp-1])
	pStaticFieldComponent := make([]uint8, sizes[TagStaticFieldComp-1])
	pReferenceLocationComponent := make([]uint8, sizes[TagReferenceLocationComp-1])
	pExportComponent := make([]uint8, sizes[TagExportComp-1])
	pDescriptorComponent := make([]uint8, sizes[TagDescriptorComp-1])
	pDebugComponent := make([]uint8, sizes[TagDebugComp-1])

	var compLength uint16
	for iPos < dataLength-1 {
		c1 = readU1(dataBuffer, &iPos)
		compLength = readU2(dataBuffer, &iPos)
		switch c1 {
		case TagImportComp:
			for j := 0; j < int(compLength); j++ {
				pImportComponent[j] = readU1(dataBuffer, &iPos)
			}

		case TagAppletComp:
			for j := 0; j < int(compLength); j++ {
				pAppletComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagConstantPoolComp:
			for j := 0; j < int(compLength); j++ {
				pConstantPoolComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagClassComp:
			for j := 0; j < int(compLength); j++ {
				pClassComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagMethodComp:
			for j := 0; j < int(compLength); j++ {
				pMethodComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagStaticFieldComp:
			for j := 0; j < int(compLength); j++ {
				pStaticFieldComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagReferenceLocationComp:
			for j := 0; j < int(compLength); j++ {
				pReferenceLocationComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagExportComp:
			for j := 0; j < int(compLength); j++ {
				pExportComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagDescriptorComp:
			for j := 0; j < int(compLength); j++ {
				pDescriptorComponent[j] = readU1(dataBuffer, &iPos)
			}
		case TagDebugComp:
			for j := 0; j < int(compLength); j++ {
				pDebugComponent[j] = readU1(dataBuffer, &iPos)
			}

		default:
			//nothing
		}
	}

	/*Parse applet component*/
	var iPosa int
	if sizes[TagAppletComp-1] == 0 {
		return nil
	}
	appletcount := readU1(pAppletComponent, &iPosa)
	applets := make([]*Applet, appletCount)
	for ia := 0; ia < int(appletcount); ia++ {
		aidLength := readU1(pAppletComponent, &iPosa)
		paidArray := make([]uint8, aidLength)
		for i = 0; i < aidLength; i++ {
			paidArray[i] = readU1(pAppletComponent, &iPosa)
		}
		installoffset := readU2(pAppletComponent, &iPosa)
		applets[ia] = &Applet{aidLength, paidArray, installoffset}
	}

	newapplet.PApplet = &AppletComponent{appletcount, applets}

	/*Parse import component*/
	iPosimp := 0
	if sizes[TagImportComp-1] == 0 {
		fmt.Println("import nil")
		return nil
	}
	pcount := readU1(pImportComponent, &iPosimp)
	packages := make([]*PackageInfo, pcount)
	for ia := 0; ia < int(pcount); ia++ {
		packageMinVer := readU1(pImportComponent, &iPosimp)
		packageMajVer := readU1(pImportComponent, &iPosimp)
		aidLength := readU1(pImportComponent, &iPosimp)
		aid := make([]uint8, aidLength)
		for i = 0; i < aidLength; i++ {
			aid[i] = readU1(pImportComponent, &iPosimp)
		}
		packages[ia] = &PackageInfo{packageMinVer, packageMajVer, aidLength, aid}
	}
	ipImport := &ImportComponent{pcount, packages}

	/*Parse class component*/

	iPosl := 0
	classCompLength := sizes[TagClassComp-1]
	nbinterface := 0
	nbclass := 0

	/*signature pool for remote methods. In our subset we don't use remote methods*/
	sPoolLenght := readU2(pClassComponent, &iPosl)
	sigPool := make([]*TypeDescriptor, sPoolLenght)
	for ia := 0; ia < int(sPoolLenght); ia++ {
		nibCount := readU1(pClassComponent, &iPosl)
		typeArray := make([]uint8, nibCount)
		l := 0
		for i := 0; i < int((nibCount+1)/2); i++ {
			val := readU1(pClassComponent, &iPosl)
			typeArray[l] = readHighShift(val)
			l++
			if l == int(nibCount) {
				continue
			}
			typeArray[l] = readLow(val)
			l++
		}
		sigPool[ia] = &TypeDescriptor{nibCount, typeArray}
	}

	//in the spec there will be at most 255 classes and interfaces
	//here we consider a subset so we will say 50 classes and 50 interfaces
	interfaces := make([]*InterfaceInfo, 50)
	classes := make([]*ClassInfo, 50)
	for iPosl < int(classCompLength-1) {
		bitfield := readU1(pClassComponent, &iPosl)
		isInterf, interfacecount := checkBitField(bitfield)
		if isInterf {
			superinterfaces := make([]uint16, interfacecount)
			for ia := 0; ia < int(interfacecount); ia++ {
				superinterfaces[ia] = readU2(pClassComponent, &iPosl)
			}
			interfaces[nbinterface] = &InterfaceInfo{bitfield, superinterfaces}
			nbinterface++
		} else {
			superClassRef := readU2(pClassComponent, &iPosl)
			declaredInstanceSize := readU1(pClassComponent, &iPosl)
			firstReferenceToken := readU1(pClassComponent, &iPosl)
			referenceCount := readU1(pClassComponent, &iPosl)
			publicMethodTableBase := readU1(pClassComponent, &iPosl)
			publicMethodTableCount := readU1(pClassComponent, &iPosl)
			packageMethodTableBase := readU1(pClassComponent, &iPosl)
			packageMethodTableCount := readU1(pClassComponent, &iPosl)
			puVMT := make([]uint16, publicMethodTableCount)
			for ia := 0; ia < int(publicMethodTableCount); ia++ {
				puVMT[ia] = readU2(pClassComponent, &iPosl)
			}
			paVMT := make([]uint16, packageMethodTableCount)
			for ia := 0; ia < int(packageMethodTableCount); ia++ {
				paVMT[ia] = readU2(pClassComponent, &iPosl)
			}
			impInterfaces := make([]*ImplementedInterfaceInfo, interfacecount)
			for ia := 0; ia < int(interfacecount); ia++ {
				interff := readU2(pClassComponent, &iPosl)
				count := readU1(pClassComponent, &iPosl)
				index := make([]uint8, count)
				for i = 0; i < count; i++ {
					index[i] = readU1(pClassComponent, &iPosl)
				}
				impInterfaces[ia] = &ImplementedInterfaceInfo{interff, count, index}
			}
			classes[nbclass] = &ClassInfo{bitfield,
				superClassRef,
				declaredInstanceSize,
				firstReferenceToken,
				referenceCount,
				publicMethodTableBase,
				publicMethodTableCount,
				packageMethodTableBase,
				packageMethodTableCount,
				puVMT,
				paVMT,
				impInterfaces}
			nbclass++
		}
	}
	ipClass := &ClassComponent{interfaces, classes}

	/*Parse constant pool component*/
	iPosc := 0
	count := readU2(pConstantPoolComponent, &iPosc)
	pCPC := make([]*CpInfo, count)
	for ia := 0; ia < int(count); ia++ {
		tag := readU1(pConstantPoolComponent, &iPosc)
		info := make([]uint8, 3)
		info[0] = readU1(pConstantPoolComponent, &iPosc)
		info[1] = readU1(pConstantPoolComponent, &iPosc)
		info[2] = readU1(pConstantPoolComponent, &iPosc)
		pCPC[ia] = &CpInfo{tag, info}
	}
	ipConstantPool := &ConstantPoolComponent{count, pCPC}

	/*Parse reference location component*/
	iPosr := 0
	byteIndexCount := readU2(pReferenceLocationComponent, &iPosr)
	offsetsToByteIndices := make([]uint8, byteIndexCount)
	for ia := 0; ia < int(byteIndexCount); ia++ {
		offsetsToByteIndices[ia] = readU1(pReferenceLocationComponent, &iPosr)
	}
	byteIndex2Count := readU2(pReferenceLocationComponent, &iPosr)
	offsetsToByte2Indices := make([]uint8, byteIndex2Count)
	for ia := 0; ia < int(byteIndex2Count); ia++ {
		offsetsToByte2Indices[ia] = readU1(pReferenceLocationComponent, &iPosr)
	}
	ipRefLoc := &ReferenceLocationComponent{byteIndexCount, offsetsToByteIndices, byteIndex2Count, offsetsToByte2Indices}

	/*Parse static image component*/
	iPoss := 0
	imageSize := readU2(pStaticFieldComponent, &iPoss)
	referenceCount := readU2(pStaticFieldComponent, &iPoss)
	arrCount := readU2(pStaticFieldComponent, &iPoss)
	pArrayInit := make([]*ArrayInitInfo, arrCount)
	for ia := 0; ia < int(arrCount); ia++ {
		typ := readU1(pStaticFieldComponent, &iPoss)
		valcount := readU2(pStaticFieldComponent, &iPoss)
		pValues := make([]uint8, valcount)
		for i := 0; i < int(valcount); i++ {
			pValues[i] = readU1(pStaticFieldComponent, &iPoss)
		}
		pArrayInit[ia] = &ArrayInitInfo{typ, valcount, pValues}
	}
	defaultValueCount := readU2(pStaticFieldComponent, &iPoss)
	nnvalcount := readU2(pStaticFieldComponent, &iPoss)
	pNonDefaultValues := make([]uint8, nnvalcount)
	for ia := 0; ia < int(nnvalcount); ia++ {
		pNonDefaultValues[ia] = readU1(pStaticFieldComponent, &iPoss)
	}
	ipStaticField := &StaticFieldComponent{imageSize,
		referenceCount,
		arrCount,
		pArrayInit,
		defaultValueCount,
		nnvalcount,
		pNonDefaultValues}

	/*Parse method component*/
	var iPosm int
	handlerCount := readU1(pMethodComponent, &iPosm)
	pExceptionHandlers := make([]*ExceptionHandlerInfo, handlerCount)
	for ia := 0; ia < int(handlerCount); ia++ {
		startOffset := readU2(pMethodComponent, &iPosm)
		activeLength := readU2(pMethodComponent, &iPosm)
		handlerOffset := readU2(pMethodComponent, &iPosm)
		catchTypeIndex := readU2(pMethodComponent, &iPosm)
		pExceptionHandlers[ia] = &ExceptionHandlerInfo{startOffset, activeLength, handlerOffset, catchTypeIndex}
	}

	msize := len(pMethodComponent) - int(8*handlerCount+1)
	pMethodInfo := make([]uint8, msize)
	for ia := 0; ia < msize; ia++ {
		pMethodInfo[ia] = readU1(pMethodComponent, &iPosm)
	}
	ipMethod := &MethodComponent{handlerCount, pExceptionHandlers, pMethodInfo}

	/*Export Component*/
	var iPose int
	ipExport := &ExportComponent{}
	if len(pExportComponent) != 0 {

		classCount := readU1(pExportComponent, &iPose)
		pClassExport := make([]*ClassExportInfo, classCount)
		for ia := 0; ia < int(classCount); ia++ {
			classOffset := readU2(pExportComponent, &iPose)
			sfc := readU1(pExportComponent, &iPose)
			smc := readU1(pExportComponent, &iPose)
			pStaticFieldOffsets := make([]uint16, sfc)
			pStaticMethodOffsets := make([]uint16, smc)
			for i := 0; i < int(sfc); i++ {
				pStaticFieldOffsets[i] = readU2(pExportComponent, &iPose)
			}
			for i := 0; i < int(smc); i++ {
				pStaticMethodOffsets[i] = readU2(pExportComponent, &iPose)
			}
			pClassExport[ia] = &ClassExportInfo{classOffset,
				sfc,
				smc,
				pStaticFieldOffsets,
				pStaticMethodOffsets}
		}
		ipExport = &ExportComponent{classCount, pClassExport}
	}

	/*descripor comp*/

	var iPosd int
	dclength := sizes[TagDescriptorComp-1]
	ccount := readU1(pDescriptorComponent, &iPosd)
	pcDis := make([]*ClassDescriptorInfo, count)
	for ia := 0; ia < int(ccount); ia++ {
		token := readU1(pDescriptorComponent, &iPosd)
		accessFlags := readU1(pDescriptorComponent, &iPosd)
		classRef := readU2(pDescriptorComponent, &iPosd)
		interCount := readU1(pDescriptorComponent, &iPosd)
		fcount := readU2(pDescriptorComponent, &iPosd)
		mcount := readU2(pDescriptorComponent, &iPosd)
		interfaces := make([]uint16, interCount)
		for i := 0; i < int(interCount); i++ {
			interfaces[i] = readU2(pDescriptorComponent, &iPosd)
		}
		fields := make([]*FieldDescriptorInfo, fcount)
		var fieldRef [3]uint8
		for i := 0; i < int(fcount); i++ {
			token := readU1(pDescriptorComponent, &iPosd)
			pAF := readU1(pDescriptorComponent, &iPosd)
			fieldRef[0] = readU1(pDescriptorComponent, &iPosd)
			fieldRef[1] = readU1(pDescriptorComponent, &iPosd)
			fieldRef[2] = readU1(pDescriptorComponent, &iPosd)
			pFieldRef := &FieldRef{fieldRef}
			pFieldtype := readU2(pDescriptorComponent, &iPosd)
			fields[i] = &FieldDescriptorInfo{token, pAF, pFieldRef, pFieldtype}
		}
		methods := make([]*MethodDescriptorInfo, mcount)
		for i := 0; i < int(mcount); i++ {
			token := readU1(pDescriptorComponent, &iPosd)
			pAF := readU1(pDescriptorComponent, &iPosd)
			methodOffset := readU2(pDescriptorComponent, &iPosd)
			typeOffset := readU2(pDescriptorComponent, &iPosd)
			bytecodeCount := readU2(pDescriptorComponent, &iPosd)
			exceptionHandlerCount := readU2(pDescriptorComponent, &iPosd)
			exceptionHandlerIndex := readU2(pDescriptorComponent, &iPosd)
			methods[i] = &MethodDescriptorInfo{token, pAF, methodOffset, typeOffset, bytecodeCount, exceptionHandlerCount, exceptionHandlerIndex}
		}

		pcDis[ia] = &ClassDescriptorInfo{token, accessFlags, classRef, interCount, fcount, mcount, interfaces, fields, methods}
	}

	cpCount := readU2(pDescriptorComponent, &iPosd)
	pConstantPoolTypes := make([]uint16, cpCount)
	for ia := 0; ia < int(cpCount); ia++ {
		pConstantPoolTypes[ia] = readU2(pDescriptorComponent, &iPosd)
	}
	pTypeDesc := make([]*TypeDescriptor, int(dclength)-iPosd-1)
	for ia := 0; iPosd < int(dclength); ia++ {
		nibCount := readU1(pDescriptorComponent, &iPosd)
		typeArray := make([]uint8, nibCount)
		fmt.Println(nibCount, iPosd)
		l := 0
		for i := 0; i < int((nibCount+1)/2); i++ {
			val := readU1(pDescriptorComponent, &iPosd)
			typeArray[l] = readHighShift(val)
			l++
			if l == int(nibCount) {
				continue
			}
			typeArray[l] = readLow(val)
			l++
		}
		pTypeDesc[ia] = &TypeDescriptor{nibCount, typeArray}
	}
	ptdi := &TypeDescriptorInfo{cpCount, pConstantPoolTypes, pTypeDesc}
	ipDescriptor := &DescriptorComponent{ccount, pcDis, ptdi}

	newapplet.AbsA = &AbstractApplet{ipHeader, ipDir, ipImport, ipClass, ipStaticField, ipMethod, ipRefLoc,
		ipConstantPool,
		ipDescriptor,
		ipExport}
	return newapplet
}
