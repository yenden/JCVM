package framework

import (
	"JCVM/core"
)

var (
	isCardInitFlag      = false
	AppletTable         = make(map[*AID]*core.CardApplet)
	currSelectedApplet  = &AID{[]byte{0, 0, 0, 0, 0, 0}}
	selectingAppletFlag = false
	processMethodFlag   = false
	TheApdu             *Apdu
	//appletTab          = make(map[*AID]*core.CardApplet)
)

func mainLoop() {
	if !isCardInitFlag {
		cardInit() // card initialization (first time only)
	}
	cardReset() // session initialization (each card reset)
	sw := int16(0)
	// main loop
	for {
		resetSelectingAppletFlag()
		resetProcessMethodFlag()
		TheApdu.complete(sw) // respond to previous APDU and get next

	}

}
func processAndForward() bool {
	switch TheApdu.buffer[OffsetIns] {
	case InsSelect: // ISO 7816-4 SELECT FIlE command
		TheApdu.selectApdu()
	default:
		//nothing .... We don't support manage channel command
	}
	return true

}
func (apdu *Apdu) selectApdu() {
	//	len = NativeMethods.t0RcvData(bOff)
}

// This function is callhe first time we init the card
//It installs all the existing applets in appletTable
func cardInit() {
	for i := range AppletTable {
		vm := initVM()
		capp := AppletTable[i]
		capp.Install(vm)
	}
	TheApdu = &Apdu{}
	TheApdu.buffer = make([]byte, 128)
	isCardInitFlag = true
}
func cardReset() {
	currSelectedApplet = &AID{[]byte{0, 0, 0, 0, 0, 0}}

}
func resetSelectingAppletFlag() {
	selectingAppletFlag = false
}
func setSelectingAppletFlag() {
	selectingAppletFlag = true
}
func resetProcessMethodFlag() {
	processMethodFlag = false
}
func setProcessMethodFlag() {
	processMethodFlag = true
}
func initVM() *core.VM {
	vm := &core.VM{}
	vm.FrameTop = -1
	vm.StackFrame = make([]*core.Frame, 30)
	f := &core.Frame{}
	vm.PushFrame(f)
	return vm
}
