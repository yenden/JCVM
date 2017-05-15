package jcre

import (
	"JCVM/core"
	"JCVM/jcre/api"
	"JCVM/jcre/nativeMethods"
	"reflect"
)

var (
	//flag to dertermine the first init
	isCardInitFlag = false
	//current selected applet aid
	currSelectedApplet = &api.AID{}
	//TheApdu Represent the apdu buffer
	TheApdu *api.Apdu
)

func MainLoop() {
	/*if !isCardInitFlag {
	cardInit() // card initialization (first time only)
		}*/
	cardInit() // card initialization (first time only)
	//cardReset() // session initialization (each card reset)
	sw := uint16(0)
	// main loop
	for {
		resetSelectingAppletFlag()
		TheApdu.Complete(sw)   // respond to previous APDU and get next
		core.SetStatus(0x9000) //reset status
		// Process channel information
		if processAndForward() { // Dispatcher handles the SELECT
			// APDU
			// dispatch to the currently selected applet
			//selectedApplet := AppletTable[currSelectedApplet]
			var selectedApplet *core.CardApplet
			for j, val := range core.AppletTable {
				if reflect.DeepEqual(j, currSelectedApplet) {
					selectedApplet = val
				}
			}
			selectedApplet.Process(initProcess())
		}
		sw = core.GetStatus()
	}
}
func processAndForward() bool {
	switch TheApdu.Buffer[api.OffsetIns] {
	case api.InsSelect: // ISO 7816-4 SELECT FIlE command
		selectApdu(TheApdu)
	case api.InsInstall: //install command
		install()
		return false
	default:
		//nothing .... We don't support manage channel command
	}
	return true

}
func selectApdu(apdu *api.Apdu) {
	Len := nativeMethods.T0RcvData(apdu.Buffer, 5)
	if Len == int16(apdu.Buffer[4]) {
		aidBytes := apdu.Buffer[5 : 5+Len]
		currSelectedApplet = api.InitAID(aidBytes, 0, int16(len(aidBytes)))
	}
	setSelectingAppletFlag()
}

func install() {
	vm := initVM()
	core.ConstantApplet.Install(vm)
}

// This function is call the first time we init the card
//It installs all the existing applets in appletTable
func cardInit() {
	/*for i := range core.AppletTable {
		vm := initVM()
		capp := core.AppletTable[i]
		capp.Install(vm)
	}
	*/
	TheApdu = &api.Apdu{}
	TheApdu.Buffer = make([]byte, 37)
	currSelectedApplet.TheAID = []byte{0, 0, 0, 0, 0, 0}
	core.InitApduArr()
	isCardInitFlag = true
}

/*
func cardReset() {

	currSelectedApplet.TheAID = []byte{0, 0, 0, 0, 0, 0}
}*/
func resetSelectingAppletFlag() {
	api.SelectingAppLetFlag = false
}
func setSelectingAppletFlag() {
	api.SelectingAppLetFlag = true
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
	//core.FillApduArr(TheApdu.Buffer, core.Reference(6000))
	vm.StackFrame[vm.FrameTop].Localvariables = make([]interface{}, 200)
	vm.StackFrame[vm.FrameTop].Localvariables[0] = core.InstanceRefHeap[currSelectedApplet]
	vm.StackFrame[vm.FrameTop].Localvariables[1] = core.Reference(6000)
	return vm
}
