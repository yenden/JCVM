package core

/*ClassExportInfo is an exported class informations */
type ClassExportInfo struct {
	classOffset          uint16
	staticFieldCount     uint8
	staticMethodCount    uint8
	pStaticFieldOffsets  []uint16
	pStaticMethodOffsets []uint16
}

/*ExportComponent of the CAP file*/
type ExportComponent struct {
	classCount   uint8
	pClassExport []*ClassExportInfo
}
