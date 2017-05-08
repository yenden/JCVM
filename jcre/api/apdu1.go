package api

import (
	"JCVM/jcre/nativeMethods"
	"log"
)

const (
	InsSelect = byte(0xA4)
)

type Apdu struct {
	Command []byte
	Buffer  []byte
}

func (apdu *Apdu) Complete(status uint16) {
	// Zero out APDU buffer
	var result int16
	var err error
	ArrayfillNonAtomic(apdu.Buffer, 0, BufferSize-5, 0)
	/*	apdu.Buffer[0] = byte(status >> 8)
		apdu.Buffer[1] = byte(status)*/
	if status == 0 {
		result, err = nativeMethods.T0RcvCommand(apdu.Command)
		if err != nil {
			log.Println(err)
		} else {
			nativeMethods.CopyToApdubuffer(apdu.Buffer, len(apdu.Buffer))
		}
	} else {
		nativeMethods.T0SetStatus(int(status))
		result = nativeMethods.T0SndStatusRcvCommand(apdu.Command)
		nativeMethods.CopyToApdubuffer(apdu.Buffer, len(apdu.Buffer))
	}
	if result != 0 {
		log.Println("imput/output error in complete method")
	}
}
