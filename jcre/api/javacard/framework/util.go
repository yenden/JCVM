package framework

func ArrayCopy(srcArray []byte, srcoffset int16, destArray []byte, destoffset int16, length int16) {
	srclen := destoffset + length
	destlen := srcoffset + length
	copy(destArray[destoffset:srclen], srcArray[srcoffset:destlen])
}
