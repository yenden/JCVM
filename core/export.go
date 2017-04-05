package core

type ClassExportInfo struct {
	classOffset          uint16
	staticFieldCount     uint8
	staticMethodCount    uint8
	pStaticFieldOffsets  []uint16
	pStaticMethodOffsets []uint16
}
type ExportComponent struct {
	classCount   uint8
	pClassExport []*ClassExportInfo
}
