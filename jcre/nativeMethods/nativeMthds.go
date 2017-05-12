package nativeMethods

import (
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	sendRcvCycleStarted = false
	conn                *net.UDPConn
	//BufferRcv used to receive the incoming
	//and outgoing data
	BufferRcv          = make([]byte, 128)
	bufferSend         = make([]byte, 128)
	command            = make([]byte, 5)
	invalidGetResponse = false
	//LC is the data length send in the apdu
	LC byte
	//LE expected length in the response
	LE byte
	//LR apdu response length
	LR byte
	sw int
	//Addr client application address
	Addr               *net.UDPAddr
	apduPtr            = int16(5)
	SendInProgressFlag = false
)

func protocolServer() *net.UDPConn {
	if conn == nil {
		hostName := "localhost"
		portNum := "6000"
		service := hostName + ":" + portNum
		udpAddr, err := net.ResolveUDPAddr("udp4", service)
		if err != nil {
			log.Fatal(err)
		}
		// setup listener for incoming UDP connection
		conn, err = net.ListenUDP("udp", udpAddr)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Card is up ")
	}
	return conn
}
func receive() *net.UDPAddr {
	apduPtr = 5
	n, addr, err := protocolServer().ReadFromUDP(BufferRcv)
	if err != nil {
		log.Fatal(err)
	}
	setParam(n)
	return addr
}
func sendStatus(sw int) {
	bs := make([]byte, 2)
	bs[0] = byte(sw >> 8)
	bs[1] = byte(sw)
	var send []byte
	if LR > 0 {
		send = append(bufferSend[0:LR], bs...)
	} else {
		send = bs
	}
	_, err := protocolServer().WriteToUDP(send, Addr)
	if err != nil {
		log.Fatal(err)
	}
}

//T0RcvCommand is to receive the command part of apdu
func T0RcvCommand(com []byte) (int16, error) {
	if sendRcvCycleStarted {
		return 0, errors.New("Error: T0RcvCommand has been already called")
	}
	//send receive cycle started
	//receive apdu and copy its command in command buffer
	sendRcvCycleStarted = true
	Addr = receive()
	copy(com[0:5], BufferRcv[0:5])
	copy(command[0:], BufferRcv[0:5])
	return 0, nil
}

//T0SndStatusRcvCommand send response and wait for next apdu
func T0SndStatusRcvCommand(com []byte) int16 {
	sendStatus(sw)
	Addr = receive()
	copy(com[0:5], BufferRcv[0:5])
	copy(command[0:], BufferRcv[0:5])
	return 0
}

//T0RcvData retrieves data form buffer
func T0RcvData(apduBuffer []byte, offset int16) int16 {
	receiveLen := int16(command[4] & 0xFF)
	receiveSpace := int16(len(apduBuffer)) - offset
	if receiveLen > receiveSpace {
		receiveLen = receiveSpace
	}
	copy(apduBuffer[offset:receiveLen+offset], BufferRcv[apduPtr:apduPtr+receiveLen])
	apduPtr += receiveLen
	command[4] -= byte(receiveLen)
	LC = command[4]
	return receiveLen
}

//T0SendData copy to the outgoing apdu buffer
func T0SendData(apduBuffer []byte, offset int16, length int16) {
	copy(bufferSend[0:length], apduBuffer[offset:length])
	//sendStatus()
}

//T0CopyToApdubuffer copy the content of buffer in another array
func T0CopyToApdubuffer(apduBuffer []byte, Len int) {
	copy(apduBuffer[0:Len], BufferRcv[0:Len])
}

//T0SetStatus set the status word to send
func T0SetStatus(status int) {
	sw = status
}

func T0SndGetResponse() int16 {
	T0SndStatusRcvCommand(command)
	return int16(command[4] & 0xFF)
}
func Send61xx(length int16) int16 {
	expLen := length
	for ok := true; ok; ok = (expLen > length) { //do... while
		// Set SW1SW2 as 61xx.
		T0SetStatus(int(0x6100 + length&0x00FF)) //61xx means data remaining
		newLen := T0SndGetResponse()
		if newLen > 0 && (newLen>>8 != 0xC0) { //0xC0xx
			LE = byte(newLen)
			expLen = int16(LE)
		}
	}
	SendInProgressFlag = false
	return expLen
}

func setParam(n int) {
	if n < 4 {
		log.Fatal("Error: Apdu len must be >4")
	} else if n == 4 { //apdu case 1 ---CLA|INS|P1|P2---
		LE = 0
		LC = 0
	} else if n == 5 { //apdu case 2 ---CLA|INS|P1|P2|Le---
		LC = 0
		LE = BufferRcv[4]
	} else if n == int(5+BufferRcv[4]) { //apdu case 3 ---CLA|INS|P1|P2|LC|Data---
		LC = BufferRcv[4]
		LE = 0
	} else { //apdu case 4 ---CLA|INS|P1|P2|LC|Data|Le---
		LC = BufferRcv[4]
		LE = BufferRcv[n-1]
	}
	LR = LE
}

//PowerUP launch the jcre and represent the first reset
func PowerUP() {
	var addr *net.UDPAddr
	var err error
	atr := []byte{59, 00, 17, 0, 00}
	for BufferRcv[0] == 0x00 && BufferRcv[1] == 0x00 {
		fmt.Println("wait for Powerup() for card activation")
		_, addr, err = protocolServer().ReadFromUDP(BufferRcv)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = protocolServer().WriteToUDP(atr, addr)
	if err != nil {
		log.Fatal(err)
	}
}
