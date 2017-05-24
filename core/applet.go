package core

/*Applet contains Applet information */
type Applet struct {
	aidLength           uint8
	pAID                []uint8
	installMethodOffset uint16
}

/*AppletComponent ...*/
type AppletComponent struct {
	count   uint8
	applets []*Applet
}
