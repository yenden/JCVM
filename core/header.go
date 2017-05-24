package core

/*HeaderComponent describes the header comp of JCVM spec*/
type HeaderComponent struct {
	MinorVersion uint8
	MajorVersion uint8
	flags        uint8
	PThisPackage *PackageInfo
	pNameInfor   *PackageNameInfo
}

/*PackageNameInfo is package's name informations*/
type PackageNameInfo struct {
	nameLength uint8
	name       []uint8
}
