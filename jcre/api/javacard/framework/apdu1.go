package framework

import (
	"JCVM/jcre/api/nativeMethods"
	"log"
)

const (
	InsSelect = byte(0xA4)
)

type Apdu struct {
	buffer []byte
}

func (apdu *Apdu) complete(status int16) {
	// Zero out APDU buffer
	var result int16
	var err error
	ArrayfillNonAtomic(apdu.buffer, 0, BufferSize, 0)
	/*	apdu.buffer[0] = byte(status >> 8)
		apdu.buffer[1] = byte(status)*/
	if status == 0 {
		result, err = nativeMethods.T0RcvCommand()
		if err != nil {
			log.Println(err)
		} else {
			nativeMethods.CopyToApdubuffer(apdu.buffer, len(apdu.buffer))
		}
	} else {
		nativeMethods.T0SetStatus(int(status))
		result = nativeMethods.T0SndStatusRcvCommand()
		nativeMethods.CopyToApdubuffer(apdu.buffer, len(apdu.buffer))
	}
	if result != 0 {
		log.Println("imput/output error in complete method")
	}
}
