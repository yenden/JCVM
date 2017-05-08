package main
/*
import (
	"JCVM/core"
	"JCVM/jcre/cadClient"
)

func main() {
	args := []string{`../test/framework.ijc`, `../test/lang.ijc`, `../test/helloword.ijc`}
	for i := 0; i < len(args)-1; i++ {
		dataBuffer := core.ReadInBuffer(args[i])
		core.Lst.PushBack(core.BuildApplet(dataBuffer))
	}
	cadClient.LaunchServer()
}


import "JCVM/core"

func main() {
	//args := os.Args[0:]
		//if args == nil || len(args) == 0 {
		fmt.Println("Usage: \n\tjcvm.exe library1.ijc library2.ijc ..... yourApplet.ijc ")
//	} else {
	var i int
	args := []string{`../test/framework.ijc`, `../test/lang.ijc`, `../test/helloword.ijc`}
	for i = 0; i < len(args)-1; i++ {
		dataBuffer := core.ReadInBuffer(args[i])
		core.Lst.PushBack(core.BuildApplet(dataBuffer))
	}
	appletBuffer := core.ReadInBuffer(args[i])
	capp := core.BuildApplet(appletBuffer)
	vm := initVM()
	capp.Install(vm)
	//}
}
func initVM() *core.VM {
	vm := &core.VM{}
	vm.FrameTop = -1
	vm.StackFrame = make([]*core.Frame, 256)
	f := &core.Frame{}
	vm.PushFrame(f)
	return vm
}*/
