package core

import "container/list"

var Lst *list.List

func findLibrary(pPI *PackageInfo) *CardApplet {
	for e := Lst.Front(); e != nil; e = e.Next() {
		if e.Value.(*CardApplet).AbsA.isThisLibrary(pPI) {
			return e.Value.(*CardApplet)
		}
	}
	return nil
}
