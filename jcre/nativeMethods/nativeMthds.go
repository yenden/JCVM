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
	//	conn net.Conn
	//data with status flag --if there is just data or data+SW
	dataWithStatusFlag = false
	//BufferRcv is used to receive the incoming
	BufferRcv = make([]byte, 128)
	//bufferSend is used to store outgoing data
	bufferSend = make([]byte, 128)
	command    = make([]byte, 5)
	sw         int
	//Addr client application address
	Addr *net.UDPAddr
	//A pointer to the received apdu buffer
	apduRcvPtr = int16(5)
	//A pointer to sending apdu buffer
	ApduSendPtr = 0
	firsttime   = true
)

func protocolServer() *net.UDPConn { //net.Conn {
	/*	if firsttime {
			fmt.Println("Card is up ")
		}
		firsttime = false*/
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

	}
	return conn
}

//PowerUP launch the jcre and represent the first reset
func PowerUP( /*connect net.Conn*/ ) {
	var addr *net.UDPAddr
	var err error
	//	conn = connect
	atr := []byte{59, 00, 17, 0, 00}
	for BufferRcv[0] == 0x00 && BufferRcv[1] == 0x00 {
		fmt.Println("wait for Powerup() for card activation")
		/*_, err = protocolServer().Read(BufferRcv)
		if err != nil {
			log.Fatal(err)
		}*/

		_, addr, err = protocolServer().ReadFromUDP(BufferRcv)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Receive PowerUp signal ")
	fmt.Println("Sending ATR")
	/*_, err = protocolServer().Write(atr)
	if err != nil {
		log.Fatal(err)
	}*/
	_, err = protocolServer().WriteToUDP(atr, addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Card is up ")
}

func receive() (*net.UDPAddr, int) {
	apduRcvPtr = 5
	dataWithStatusFlag = false
	/*	n, err := protocolServer().Read(BufferRcv)
		if err != nil {
			log.Fatal(err)
		}*/
	n, addr, err := protocolServer().ReadFromUDP(BufferRcv)
	if err != nil {
		log.Fatal(err)
	}
	return addr, n
	//	return n
}
func sendStatus(sw int) {
	bs := make([]byte, 2)
	bs[0] = byte(sw >> 8)
	bs[1] = byte(sw)
	var send []byte

	//if there is data to send with SW
	if dataWithStatusFlag {
		send = append(bufferSend[0:ApduSendPtr], bs...)
	} else {
		send = bs
	}
	/*_, err := protocolServer().Write(send)
	if err != nil {
		log.Fatal(err)
	}*/
	_, err := protocolServer().WriteToUDP(send, Addr)
	if err != nil {
		log.Fatal(err)
	}
}

//T0RcvCommand is to receive the command part of apdu
func T0RcvCommand(com []byte) (int, error) {
	if sendRcvCycleStarted {
		return 0, errors.New("Error: T0RcvCommand has been already called")
	}
	//send receive cycle started
	//receive apdu and copy its command in command buffer
	sendRcvCycleStarted = true
	var n int
	//n = receive()
	Addr, n = receive()
	copy(com[0:5], BufferRcv[0:5])
	copy(command[0:], BufferRcv[0:5])
	return n, nil
}

//T0SndStatusRcvCommand send sw and wait for next apdu
func T0SndStatusRcvCommand(com []byte) int {
	sendStatus(sw)
	var n int
	//n = receive()
	Addr, n = receive()
	copy(com[0:5], BufferRcv[0:5])
	copy(command[0:], BufferRcv[0:5])
	return n
}

//T0RcvData retrieves data form buffer
func T0RcvData(apduBuffer []byte, offset int16) int16 {
	receiveLen := int16(command[4] & 0xFF)
	receiveSpace := int16(len(apduBuffer)) - offset
	if receiveLen > receiveSpace {
		receiveLen = receiveSpace
	}
	copy(apduBuffer[offset:receiveLen+offset], BufferRcv[apduRcvPtr:apduRcvPtr+receiveLen])
	apduRcvPtr += receiveLen
	command[4] -= byte(receiveLen)
	//LC = command[4]
	return receiveLen
}

//T0SendData copy to the outgoing apdu buffer
func T0SendData(apduBuffer []byte, offset int16, length int16) {
	if !dataWithStatusFlag {
		dataWithStatusFlag = true
	}
	copy(bufferSend[ApduSendPtr:ApduSendPtr+int(length)], apduBuffer[offset:length])
	ApduSendPtr += int(length)
	/*_, err := protocolServer().Write(apduBuffer[offset:length])
	if err != nil {
		log.Fatal(err)
	}*/

	/*
		_, err := protocolServer().WriteToUDP(apduBuffer[offset:length], Addr)
		if err != nil {
			log.Fatal(err)
		}*/
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
