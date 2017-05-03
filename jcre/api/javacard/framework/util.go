package framework

func ArrayCopy(srcArray []byte, srcoffset int16, destArray []byte, destoffset int16, length int16) int16 {
	srclen := destoffset + length
	destlen := srcoffset + length
	copy(destArray[destoffset:srclen], srcArray[srcoffset:destlen])
	return destoffset + length
}
func ArrayCopyNonAtomic(src []byte, srcOff int16, dest []byte, destOff int16, length int16) int16 {
	srclen := destOff + length
	destlen := srcOff + length
	copy(dest[destOff:srclen], src[srcOff:destlen])
	return destOff + length
}
