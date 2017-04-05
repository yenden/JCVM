package core

var imageCounter int

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
	//	pStaticFieldImage    []uint8
}

/*func (sfc *StaticFieldComponent) beginBuildNonDefaultValues(defaultvaluecount uint16, nondefaultvaluecount uint16) {
	sfc.defaultValueCount = defaultvaluecount
	sfc.nonDefaultValueCount = nondefaultvaluecount
	sfc.pNonDefaultValues = make([]uint8, nondefaultvaluecount)
}*/
