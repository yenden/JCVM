package api

import (
	"errors"
)

/*ArrayCopy copies an array in another */
func ArrayCopy(srcArray []byte, srcoffset int16, destArray []byte, destoffset int16, length int16) int16 {
	srclen := destoffset + length
	destlen := srcoffset + length
	copy(destArray[destoffset:srclen], srcArray[srcoffset:destlen])
	return destoffset + length
}

/*ArrayCopyNonAtomic copies an array in another */
func ArrayCopyNonAtomic(src interface{}, srcOff int16, dest interface{}, destOff int16, length int16) int16 {
	srclen := destOff + length
	destlen := srcOff + length
	switch src.(type) {
	case []uint8:
		copy(dest.([]uint8)[destOff:srclen], src.([]uint8)[srcOff:destlen])
	}
	return destOff + length
}

/*ArrayfillNonAtomic fiels an array with a value */
func ArrayfillNonAtomic(bArray []byte, bOff int16, bLen int16, bValue byte) (int16, error) {
	if bLen < 0 {
		return 0, errors.New("Error in ArrayFillNonatomic")
	}

	for ; bLen > 0; bLen-- {
		bArray[bOff] = bValue
		bOff++
	}
	return bOff + bLen, nil
}

/*ArrayCompare compare two parts of two arrays */
func ArrayCompare(src []byte, srcOff int16, dest []byte, destOff int16, length int16) (int, error) {
	if length < 0 {
		return 0, errors.New("Error length<0 in array compare")
	}
	for i := 0; i < int(length); i++ {
		if src[int(srcOff)+i] != dest[int(destOff)+i] {
			if src[int(srcOff)+i] < dest[int(destOff)+i] {
				return -1, nil
			}
			return 1, nil
		}
	}
	return 0, nil
}

/*MakeShort makes short from 2 bytes*/
func MakeShort(b1 byte, b2 byte) int16 {
	return int16(b1)<<8 + int16(b2)&0xFF
}

/*GetShort gets short from
* 2  byte elements of an array
 */
func GetShort(bArray []byte, bOff int16) int16 {
	return int16(bArray[bOff])<<8 + int16(bArray[bOff+1])&0x0FF
}

/*SetShort sets short in
* 2  byte elements of an array
 */
func SetShort(bArray []byte, bOff int16, sValue int16) int16 {
	bArray[bOff] = byte(sValue >> 8)
	bArray[bOff+1] = byte(sValue)
	return int16(bOff + 2)
}

/*CheckArrayArgs checks the offset and the length of an array*/
func CheckArrayArgs(bArray []byte, offset int16, length int16) error {
	ln := int16(len(bArray) - 1)
	if offset < 0 && offset > ln && length <= ln-offset {
		return nil
	}
	return errors.New("Array arguments exception")
}
