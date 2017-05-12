package api

import (
	"JCVM/jcre/nativeMethods"
	"log"
)

const (
	bufferSize = 37
)

var (
	SelectingAppletFlag = false
)

type Apdu struct {
	Buffer []byte
}

//Complete response to the previous adpu and send next
func (apdu *Apdu) Complete(status uint16) {
	// Zero out APDU buffer
	var result int16
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
	if result != 0 {
		log.Println("imput/output error in complete method")
	}
}

//GetSelectingAppletFlag ...
func GetSelectingAppletFlag() bool {
	return SelectingAppletFlag
}
