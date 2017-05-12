package core

import "JCVM/jcre/api"

var (
	AppletTable    = make(map[*api.AID]*CardApplet)
	ConstantApplet = BuildApplet(ReadInBuffer(`../test/helloword.ijc`))
)
