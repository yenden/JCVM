package main

import (
	"JCVM/core"
	"JCVM/jcre"
	"JCVM/jcre/nativeMethods"
)

func main() {
	args := []string{ /*`../test/framework.ijc`,*/ `../ijcFiles/lang.ijc`}
	for i := 0; i < len(args); i++ {
		dataBuffer := core.ReadInBuffer(args[i])
		core.Lst.PushBack(core.BuildApplet(dataBuffer))
	}
	nativeMethods.PowerUP()
	jcre.MainLoop()
}
