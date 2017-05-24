package core

import "container/list"

//Lst is the api's packages list which can converted to CAP file
var Lst = list.New()

//*list.List
func findLibrary(pPI *PackageInfo) *CardApplet {
	for e := Lst.Front(); e != nil; e = e.Next() {
		if e.Value.(*CardApplet).AbsA.isThisLibrary(pPI) {
			return e.Value.(*CardApplet)
		}
	}
	return nil
}
