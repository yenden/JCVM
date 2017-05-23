package api

import (
	"JCVM/jcre/nativeMethods"
	"log"
)

const (
	//max size of apdubuffer stored in memory
	//this is not the size of the received apdu
	bufferSize = 37
)

var (
	//SelectingAppLetFlag ...
	SelectingAppLetFlag = false
	//Lc is the data Length send in the apdu
	Lc byte
	//Le expected Length in the response
	Le byte
	//Lr apdu response Length
	Lr byte
	//SendInProgressFlag flag whiLe sending is not finished
	sendInProgressFlag = false
)

/*Apdu represent a received apdu */
type Apdu struct {
	Buffer []byte
}

//Complete response to the previous adpu and send next
func (apdu *Apdu) Complete(status uint16) {
	// Zero out APDU buffer
	var result int
	var err error
	ArrayfillNonAtomic(apdu.Buffer, 0, bufferSize, 0)
	if status == 0 {
		result, err = nativeMethods.T0RcvCommand(apdu.Buffer)
		if err != nil {
			log.Println(err)
		}

	} else {
		nativeMethods.T0SetStatus(int(status))
		result = nativeMethods.T0SndStatusRcvCommand(apdu.Buffer)
	}
	if result == 0 {
		log.Println("imput/output error in compLete method")
	}
	setParam(result)
	nativeMethods.ApduSendPtr = 0
}

//GetSelectingAppLetFlag ...
func GetSelectingAppLetFlag() bool {
	return SelectingAppLetFlag
}

/*send sw for data remaining
* It send the status and waits for getResponse apdu
 */
func send61xx(Length int16) int16 {
	expLen := Length
	for ok := true; ok; ok = (expLen > Length) { //do... whiLe
		// Set SW1SW2 as 61xx (data remaining).
		nativeMethods.T0SetStatus(int(0x6100 + Length&0x00FF))
		newLen := nativeMethods.T0SndGetResponse()

		//0xC0xx <=>invalid getResponse apdu
		if newLen > 0 && (newLen>>8 != 0xC0) {
			Le = byte(newLen)
			expLen = int16(Le)
		}
	}
	sendInProgressFlag = false
	return expLen
}

/*SendBytes api framework function*/
func SendBytes(arr []byte, offset, Length int16) {
	for Length > 0 {
		temp := Length
		// Need to force GET RESPONSE for Case 4 & for partial blocks
		if Length != int16(Lr) || Lr != Le || sendInProgressFlag {
			temp = send61xx(Length) //send data remainig status
		}
		nativeMethods.T0SendData(arr, offset, temp)
		sendInProgressFlag = true
		offset += temp
		Length -= temp
		Lr -= byte(temp)
		Le = Lr
	}
	sendInProgressFlag = false
}

/*SendBytesLong : api function to send a long buffer*/
func SendBytesLong(Len, bOff int16, outData, apduBuff []byte) {
	CheckArrayArgs(outData, bOff, Len)
	sendLength := int16(len(apduBuff))
	for Len > 0 {
		if Len < sendLength {
			sendLength = Len
		}
		ArrayCopy(outData, bOff, apduBuff, 0, sendLength)
		SendBytes(apduBuff, 0, sendLength)
		Len -= sendLength
		bOff += sendLength
	}
}

/*SetOutgoingAndSend : api function to send outgoing bytes*/
func SetOutgoingAndSend(arr []byte, Len, bOff int16) {
	SetOutgoing()
	SetOutgoingLength(Len)
	SendBytes(arr, bOff, Len)
}

/*ReceiveBytes api framework function*/
func ReceiveBytes(arr []byte, offset int16) int16 {
	//Only APDUs case 3 and 4 are expected to call this method.
	Length := nativeMethods.T0RcvData(arr, offset)
	return Length
}

/*SetIncomingandreceive api framework function*/
func SetIncomingandreceive(arr []byte) int16 {
	Length := nativeMethods.T0RcvData(arr, int16(OffsetCData))
	return Length
}

/*SetOutgoing api framework function*/
func SetOutgoing() byte {
	return Le
}

/*SetOutgoingLength api framework function*/
func SetOutgoingLength(Len int16) {
	Lr = byte(Len)
}

/*GetBuffer get the apdu buffer*/
func GetBuffer(apduarray []byte) {
	copy(apduarray[0:], nativeMethods.BufferRcv[:len(apduarray)])
}

/* setParam set apdu parameter(Lr,Lc and Le)
* purpose: know if all data has been sent
* or there is data remaining
 */
func setParam(n int) {
	if n < 4 {
		log.Fatal("Error: Apdu Len must be >4")
	} else if n == 4 { //apdu case 1 ---CLA|INS|P1|P2---
		Le = 0
		Lc = 0
	} else if n == 5 { //apdu case 2 ---CLA|INS|P1|P2|Le---
		Lc = 0
		Le = nativeMethods.BufferRcv[4]
	} else if n == int(5+nativeMethods.BufferRcv[4]) { //apdu case 3 ---CLA|INS|P1|P2|Lc|Data---
		Lc = nativeMethods.BufferRcv[4]
		Le = 0
	} else { //apdu case 4 ---CLA|INS|P1|P2|Lc|Data|Le---
		Lc = nativeMethods.BufferRcv[4]
		Le = nativeMethods.BufferRcv[n-1]
	}
	if Le == 0x7F {
		Le = 0
	}
	Lr = Le
}
