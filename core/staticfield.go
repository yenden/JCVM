package core

var (
	imageCounter int
)

type ArrayInitInfo struct {
	typ     uint8
	count   uint16
	pValues []uint8
}
type StaticFieldComponent struct {
	imageSize            uint16
	referenceCount       uint16
	arrayInitCount       uint16
	pArrayInit           []*ArrayInitInfo
	defaultValueCount    uint16
	nonDefaultValueCount uint16
	pNonDefaultValues    []uint8
	pStaticFieldImage    []uint8
}

/*func (sfc *StaticFieldComponent) beginBuildNonDefaultValues(defaultvaluecount uint16, nondefaultvaluecount uint16) {
	sfc.defaultValueCount = defaultvaluecount
	sfc.nonDefaultValueCount = nondefaultvaluecount
	sfc.pNonDefaultValues = make([]uint8, nondefaultvaluecount)
}*/
func (sfc *StaticFieldComponent) buildStaticFieldImage() {
	refccount := sfc.referenceCount
	sfc.pStaticFieldImage = make([]uint8, sfc.imageSize)
	//Build segment 1 and segment 2 data.
	//Segment 1 - arrays of primitive types initialized by <clinit> methods.
	//Segment 2 - reference types initialized to null, including arrays.
	for imageCounter = 0; imageCounter < int(sfc.referenceCount)*2; imageCounter++ {
		sfc.pStaticFieldImage[imageCounter] = 0
	}
	//Update segment 3
	//Segment 3 - primitive types initialized to default values.
	for i := 0; i < int(sfc.defaultValueCount); i++ {
		imageCounter++
		sfc.pStaticFieldImage[imageCounter] = 0
	}
	//Update segment 4
	//Segment 4 - primitive types initialized to non-default values.
	for i := 0; i < int(sfc.nonDefaultValueCount); i++ {
		imageCounter++
		sfc.pStaticFieldImage[imageCounter] = sfc.pNonDefaultValues[imageCounter]
	}
}
