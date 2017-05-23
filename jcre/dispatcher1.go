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

/*MainLoop is the main function of the JCRE
* It receives apdu and dispatches it to
* the corresponding applet
 */
func MainLoop() {
	/*if !isCardInitFlag {
	cardInit() // card initialization (first time only)
		}*/
	cardInit() // card initialization (first time only)
	//cardReset() // session initialization (each card reset)

	sw := uint16(0) //status word
	// main loop
	for {
		resetSelectingAppletFlag()
		TheApdu.Complete(sw)   // respond to previous APDU and get next
		core.SetStatus(0x9000) //reset status

		if processAndForward() { // Dispatcher handles the SELECT
			// APDU
			// dispatch to the currently selected applet
			var selectedApplet *core.CardApplet
			for i, j := range core.AppletTable {
				if reflect.DeepEqual(i, currSelectedApplet) {
					selectedApplet = j
					selectedApplet.Process(initProcess())
					sw = core.GetStatus()
					break
				}
			}
			if selectedApplet == nil {
				sw = api.SwUnknown
			}
		} else { //install meth
			sw = core.GetStatus()
		}
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

/*This function is used by the jcre to
* retrive aid in the received apdu and select
* the corresponding applet
 */
func selectApdu(apdu *api.Apdu) {
	Len := nativeMethods.T0RcvData(apdu.Buffer, 5)
	if Len == int16(apdu.Buffer[4]) {
		aidBytes := apdu.Buffer[5 : 5+Len]
		currSelectedApplet = api.InitAID(aidBytes, 0, int16(len(aidBytes)))
	}
	setSelectingAppletFlag()
}

/*Install is used by the JCRE to install
* the currently selected applet
 */
func install() {
	vm := initVM()
	core.ConstantApplet.Install(vm)
}

/* This function is supposed to be called
* the first time we init the card
 */

func cardInit() {
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

/*to reset selecting applet Flag */
func resetSelectingAppletFlag() {
	api.SelectingAppLetFlag = false
}

/*to set selecting applet Flag */
func setSelectingAppletFlag() {
	api.SelectingAppLetFlag = true
}

//create an instance of the vm
func initVM() *core.VM {
	vm := &core.VM{}
	vm.FrameTop = -1
	vm.StackFrame = make([]*core.Frame, 30)
	f := &core.Frame{}
	vm.PushFrame(f)
	return vm
}

/*init the vm and store the reference 60000
* in its local variables.
* Reference 6000 is the reference of the apdu buffer on the heap
 */
func initProcess() *core.VM {
	vm := initVM()
	vm.StackFrame[vm.FrameTop].Localvariables = make([]interface{}, 200)
	vm.StackFrame[vm.FrameTop].Localvariables[0] = core.InstanceRefHeap[currSelectedApplet]
	vm.StackFrame[vm.FrameTop].Localvariables[1] = core.Reference(6000)
	return vm
}
