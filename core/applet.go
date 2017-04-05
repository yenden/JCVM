package core

type Applet struct {
	aidLength           uint8
	pAID                []uint8
	installMethodOffset uint16
}

type AppletComponent struct {
	count   uint8
	applets []*Applet
}
