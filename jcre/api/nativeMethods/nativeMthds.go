package nativeMethods

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	sendRcvCycleStarted = false
	conn                *net.UDPConn
	bufferRcv           = make([]byte, 128) //scratch buffer
	command             = make([]byte, 5)
	invalidGetResponse  = false
	LC, LE              byte
	sw                  int
	Cond                = false //one outgoing in process?TODO
)

func T0RcvCommand() (int16, error) {
	if sendRcvCycleStarted {
		return 0, errors.New("Error: T0RcvCommand has beel already called")
	}
	//send receive cycle started
	//receive apdu and copy its command in command buffer
	sendRcvCycleStarted = true
	receive()
	copy(command[0:], bufferRcv[0:5])
	return 0, nil
}
func protocolServer() *net.UDPConn {
	atr := []byte{59, 00, 17, 0, 00}
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
		conn.Write(atr)
	}
	return conn
}
func CopyToApdubuffer(apduBuffer []byte, Len int) {
	copy(apduBuffer[0:Len], bufferRcv[0:Len])
}
func receive() {
	n, err := protocolServer().Read(bufferRcv)
	if err != nil {
		log.Fatal(err)
	}
	setParam(n)
}
func sendStatus(sw int) {
	if !Cond { //we did'nt have a send in the last process function execution
		bufsw := new(bytes.Buffer)
		binary.Write(bufsw, binary.LittleEndian, sw)
		_, err := protocolServer().Write(bufsw.Bytes())
		if err != nil {
			log.Fatal(err)
		}
	}

}

func T0SndStatusRcvCommand() int16 {
	sendStatus(sw)
	receive()
	copy(command[0:], bufferRcv[0:5])
	return 0
}
func T0SetStatus(status int) {
	sw = status
}
func setParam(n int) {
	if n < 4 {
		log.Fatal("Error: Apdu len must be >4")
	} else if n == 4 { //apdu case 1 ---CLA|INS|P1|P2---
		LE = 0
		LC = 0
	} else if n > 4 {
		if bufferRcv[4] == 0 { //apdu case 2 ---CLA|INS|P1|P2|LC---
			LC = 0
			LE = 0
		} else {
			LC = bufferRcv[4]
			dataLen := bufferRcv[5]
			if n == int(5+dataLen) { //apdu case 3 ---CLA|INS|P1|P2|LC|Data---
				LE = 0
			} else { //apdu case 4 ---CLA|INS|P1|P2|LC|Data|Le---
				LE = bufferRcv[n-1]
			}

		}
	}
}
func T0RcvData(apduBuffer []byte, offset int16) int16 {
	receiveLen := int16(command[4] & 0xFF)
	copy(apduBuffer[offset:receiveLen], bufferRcv[offset:receiveLen])
	return receiveLen

}
