package core

/*
u1<--> uint8
u2<--> uint16
u3<--> uint32

s1<--> int8
s2<--> int16
s4<--> int32
*/

func readU1(dataBuffer []byte, iPos *int) uint8 {
	temp := dataBuffer[*iPos]
	(*iPos)++
	return temp
}
func readS1(dataBuffer []byte, iPos *int) int8 {
	temp := int8(dataBuffer[*iPos])
	(*iPos)++
	return temp
}

func readU2(dataBuffer []byte, iPos *int) uint16 {
	temp := ((uint16(dataBuffer[*iPos]) & 0x00FF) << 8) + uint16(dataBuffer[*iPos+1])
	(*iPos) += 2
	return temp
}

func readS2(dataBuffer []byte, iPos *int) int16 {
	temp := (int16(dataBuffer[*iPos]) << 8) + int16(dataBuffer[*iPos+1])
	(*iPos) += 2
	return temp
}

func readU4(dataBuffer []byte, iPos *int) uint32 {
	temp := (uint32(dataBuffer[*iPos]) << 24) + (uint32(dataBuffer[*iPos+1]) << 16) + (uint32(dataBuffer[*iPos+2]) << 8) + uint32(dataBuffer[*iPos+3])
	(*iPos) += 4
	return temp
}

func readS4(dataBuffer []byte, iPos *int) int32 {
	temp := (int32(dataBuffer[*iPos]) << 24) + (int32(dataBuffer[*iPos+1]) << 16) + (int32(dataBuffer[*iPos+2]) << 8) + int32(dataBuffer[*iPos+3])
	(*iPos) += 4
	return temp
}

func readHigh(data uint8) uint8 {
	return data & 0xF0
}

func readLow(data uint8) uint8 {
	return data & 0x0F
}

func readHighShift(data uint8) uint8 {
	return (data & 0xF0) >> 4
}

func makeU2(byte1 uint8, byte2 uint8) uint16 {
	return uint16(byte1)*0x100 + uint16(byte2)
}

func makeInt(short1 int16, short2 int16) int32 {
	return int32(short1)*0x10000 + int32(short2)
}

func getShortHigh(value int32) int16 {
	return int16(value / 0x10000)
}

func getShortLow(value int32) int16 {
	return int16(value % 0x10000)
}
