package core

var (
	imageCounter int
)

/*ArrayInitInfo : informations on arrays that has been initialized in <clinit>*/
type ArrayInitInfo struct {
	typ     uint8
	count   uint16
	pValues []uint8
}

/*StaticFieldComponent of CAP file*/
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

/*this builds the staticfield Image in an array
 */
func (sfc *StaticFieldComponent) buildStaticFieldImage() {
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
		sfc.pStaticFieldImage[imageCounter] = 0
		imageCounter++

	}
	//Update segment 4
	//Segment 4 - primitive types initialized to non-default values.
	for i := 0; i < int(sfc.nonDefaultValueCount); i++ {
		sfc.pStaticFieldImage[imageCounter] = sfc.pNonDefaultValues[i]
		imageCounter++
	}
}
