package core

/*HeaderComponent describes the header comp of JCVM spec*/
type HeaderComponent struct {
	MinorVersion uint8
	MajorVersion uint8
	flags        uint8
	pThisPackage *PackageInfo
	pNameInfor   *PackageNameInfo
}
type PackageNameInfo struct {
	nameLength uint8
	name       []uint8
}
