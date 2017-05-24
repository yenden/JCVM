package core

/*ReferenceLocationComponent of CAP file*/
type ReferenceLocationComponent struct {
	byteIndexCount        uint16
	offsetsToByteIndices  []uint8
	byte2IndexCount       uint16
	offsetsToByte2Indices []uint8
}
