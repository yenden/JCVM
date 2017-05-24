package core

import "JCVM/jcre/api"

/*This file has been created because of golang architecture*/
var (
	AppletTable    = make(map[*api.AID]*CardApplet)
	ConstantApplet = BuildApplet(ReadInBuffer(`../ijcFiles/cha8.ijc`))
)
