package share

import (
	"JCVM/core"
	"JCVM/jcre/api"
	"fmt"
	"log"
	"net"
	"reflect"
)

var (
	appletTab          = make(map[*api.AID]*core.CardApplet)
	currentSelectedApp *api.AID
)

func LaunchServer() {
	hostName := "localhost"
	portNum := "6000"
	service := hostName + ":" + portNum

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	ln, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Card is up and listening on port 6000")
	defer ln.Close()
	for {
		// wait for UDP client to connect
		handleAPDU(ln)

	}
}
func handleAPDU(conn *net.UDPConn) {
	buffer := make([]byte, 128)
	response := make([]byte, 4)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}
	switch buffer[0] {
	case 0x00: //install or select applet  or begin process command
		switch buffer[1] { //if it is an install applet command
		case 0x01:
			if len(appletTab) >= 16 {
				// there is already 16 applets in memory
				response[0] = 0x00
				response[1] = 0x00
				// NOTE : Need to specify client address in WriteToUDP() function
				//        otherwise, you will get this error message
				//        write udp : write: destination address required if you use Write() function instead of WriteToUDP()
				_, err = conn.WriteToUDP(response, addr)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				abspath := string(buffer[5:n])
				appletBuffer := core.ReadInBuffer(abspath)
				capp := core.BuildApplet(appletBuffer)
				aidbytes := capp.AbsA.PHeader.PThisPackage.AID
				aid := api.InitAID(aidbytes, 0, int16(len(aidbytes)))
				appletTab[aid] = capp
				vm := initVM()
				capp.Install(vm)

				response[0] = 0x90
				response[1] = 0x00
				_, err = conn.WriteToUDP(response, addr)
				if err != nil {
					log.Fatal(err)
				}
			}

		case 0x02: //if it is a select applet command
			aidBytes := buffer[5:n]
			aid := api.InitAID(aidBytes, 0, int16(len(aidBytes)))
			if !appletExists(aid) {
				log.Println("Error : this applet doesn't exists")
				response[0] = 0x55
				response[1] = 0x55
				_, err = conn.WriteToUDP(response, addr)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				setCurrentlySelectedApp(aid)
				/*vm := initVM()
				capp := appletTab[aid]
				capp.Process(vm) //TODO change this*/
				response[0] = 0x90
				response[1] = 0x00
				_, err = conn.WriteToUDP(response, addr)
				if err != nil {
					log.Fatal(err)
				}
			}
		case 0x03: //if it is a begin process command
			vm := initVM()
			preProcessedVM(vm)
			capp := appletTab[currentSelectedApp]
			capp.Process(vm)
		}

	default: //process APDU command or any other command
		//

	}

}

func initVM() *core.VM {
	vm := &core.VM{}
	vm.FrameTop = -1
	vm.StackFrame = make([]*core.Frame, 30)
	f := &core.Frame{}
	vm.PushFrame(f)
	return vm
}
func appletExists(aid *api.AID) bool {
	for searchAaid, _ := range appletTab {
		if reflect.DeepEqual(aid, searchAaid) {
			return true
		}
	}
	return false
}
func setCurrentlySelectedApp(aid *api.AID) {
	currentSelectedApp = aid
	//log.Println("Currently selected applet", currentSelectedApp)
}
func preProcessedVM(vm *core.VM) {
	vm.StackFrame[vm.FrameTop].Localvariables[0] = core.Reference(core.InstanceRefHeap[currentSelectedApp])
	vm.StackFrame[vm.FrameTop].Localvariables[1] = core.Reference(6000)
}
