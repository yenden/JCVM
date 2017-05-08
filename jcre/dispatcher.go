package jcre

import (
	"JCVM/core"
	"JCVM/jcre/api"
	"JCVM/jcre/nativeMethods"
	"reflect"
)

var (
	isCardInitFlag      = false
	AppletTable         = make(map[*api.AID]*core.CardApplet)
	currSelectedApplet  = &api.AID{}
	selectingAppletFlag = false
	processMethodFlag   = false
	TheApdu             *api.Apdu
	//appletTab          = make(map[*AID]*core.CardApplet)
)

func MainLoop() {
	if !isCardInitFlag {
		cardInit() // card initialization (first time only)
	}
	cardReset() // session initialization (each card reset)
	sw := uint16(0)
	// main loop
	for {
		resetSelectingAppletFlag()
		resetProcessMethodFlag()
		TheApdu.Complete(sw) // respond to previous APDU and get next
		// Process channel information
		if processAndForward() { // Dispatcher handles the SELECT
			// APDU
			// dispatch to the currently selected applet
			//selectedApplet := AppletTable[currSelectedApplet]
			var selectedApplet *core.CardApplet
			for j, val := range AppletTable {
				if reflect.DeepEqual(j, currSelectedApplet) {
					selectedApplet = val
				}
			}
			selectedApplet.Process(initProcess())
		}
		sw = uint16(0x9000)
	}

}
func processAndForward() bool {
	switch TheApdu.Command[api.OffsetIns] {
	case api.InsSelect: // ISO 7816-4 SELECT FIlE command
		selectApdu(TheApdu)
	default:
		//nothing .... We don't support manage channel command
	}
	return true

}
func selectApdu(apdu *api.Apdu) {
	Len := nativeMethods.T0RcvData(apdu.Command, apdu.Buffer, 5)
	if Len == int16(apdu.Command[4]) {
		aidBytes := apdu.Buffer[0:Len]
		currSelectedApplet = api.InitAID(aidBytes, 0, int16(len(aidBytes)))
	}
}

// This function is callhe first time we init the card
//It installs all the existing applets in appletTable
func cardInit() {
	for i := range AppletTable {
		vm := initVM()
		capp := AppletTable[i]
		capp.Install(vm)
	}
	TheApdu = &api.Apdu{}
	TheApdu.Buffer = make([]byte, 128-5)
	TheApdu.Command = make([]byte, 5)
	currSelectedApplet.TheAID = []byte{0, 0, 0, 0, 0, 0}
	core.InitApduArr()
	isCardInitFlag = true
}
func cardReset() {
	currSelectedApplet = &api.AID{[]byte{0, 0, 0, 0, 0, 0}}

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
func initProcess() *core.VM {
	vm := initVM()
	apduBuff := append(TheApdu.Command, TheApdu.Buffer...)
	core.FillApduArr(apduBuff, core.Reference(6000))
	vm.StackFrame[vm.FrameTop].Localvariables = make([]interface{}, 256)
	vm.StackFrame[vm.FrameTop].Localvariables[0] = core.InstanceRefHeap[currSelectedApplet]
	vm.StackFrame[vm.FrameTop].Localvariables[1] = core.Reference(6000)
	return vm
}
