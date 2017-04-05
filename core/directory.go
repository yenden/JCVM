package core

type StaticFieldSizeInfo struct {
	imageSize      uint16
	arrayInitCount uint16
	arrayInitSize  uint16
}

/*type CustomComponentInfo struct {
	componentTag uint8
	size         uint16
	aidLength    uint8
	aid          []uint8
}*/

type DirectoryComponent struct {
	componentSizes   [12]uint16
	pStaticFieldSize *StaticFieldSizeInfo
	importCount      uint8
	appletCount      uint8
	//customCount      uint8
	//pCustomComponents []*CustomComponentInfo
}
