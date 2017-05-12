package main

import (
	"JCVM/core"
	"JCVM/jcre"
	"JCVM/jcre/nativeMethods"
)

func main() {
	args := []string{ /*`../test/framework.ijc`,*/ `../test/lang.ijc`}
	for i := 0; i < len(args); i++ {
		dataBuffer := core.ReadInBuffer(args[i])
		core.Lst.PushBack(core.BuildApplet(dataBuffer))
	}
	nativeMethods.PowerUP()
	jcre.MainLoop()
	/*
		aidBytes := []byte{0xA0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}
		aid := api.InitAID(aidBytes, 0, int16(len(aidBytes)))
		jcre.AppletTable[aid] = core.BuildApplet(core.ReadInBuffer(`../test/shortnew.ijc`))
		fmt.Println(aid)
		jcre.MainLoop()*/

}
