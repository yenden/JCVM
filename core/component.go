package core

/*Different components in the cap file*/
const (
	TagHeaderComp            = 0x01
	TagDirComp               = 0x02
	TagAppletComp            = 0x03
	TagImportComp            = 0x04
	TagConstantPoolComp      = 0x05
	TagClassComp             = 0x06
	TagMethodComp            = 0x07
	TagStaticFieldComp       = 0x08
	TagReferenceLocationComp = 0x09
	TagExportComp            = 0x0A
	TagDescriptorComp        = 0x0B
	TagDebugComp             = 0x0C

	AccPublic    = 0x01
	AccPrivate   = 0x02
	AccProtected = 0x04
	AccStatic    = 0x08
	AccFinal     = 0x10
	AccAbstract  = 0x40
	AccInit      = 0x80
)

/*PackageInfo gathers information on the package*/
type PackageInfo struct {
	MinorVersion uint8
	MajorVersion uint8
	AIDLength    uint8
	AID          []uint8
}

func createPackageInfo(miorVersion uint8, MajorVersion uint8, AIDLength uint8) *PackageInfo {
	aidArr := make([]uint8, AIDLength)
	pi := &PackageInfo{miorVersion, MajorVersion, AIDLength, aidArr}
	return pi
}

/*CompareOperator to compare 2 PackageInfo*/
func (p *PackageInfo) CompareOperator(p2 *PackageInfo) bool {
	if p.MajorVersion != p2.MajorVersion {
		return false
	}
	if p.MinorVersion != p2.MinorVersion {
		return false
	}
	if p.AIDLength != p2.AIDLength {
		return false
	}
	var i uint8
	for i = 0; i < p.AIDLength; i++ {
		if p.AID[i] != p2.AID[i] {
			return false
		}
	}
	return true
}

/*ClassRef represent the reference to a class*/
type ClassRef struct {
	classref uint16
}

/*CompareOperator to compare 2 ClassRef*/
func (c *ClassRef) CompareOperator(cf *ClassRef) bool {
	return (c.classref == cf.classref)
}

/*AccessFlag ...*/
/*type AccessFlag struct {
	value uint8
}*/

/*IsPublic ..*/
func IsPublic(value uint8) bool {
	return (value & AccPublic) == AccPublic
}

/*IsPrivate ...*/
func IsPrivate(value uint8) bool {
	return (value & AccPrivate) == AccPrivate
}

/* IsProtected  ...*/
func IsProtected(value uint8) bool {
	return (value & AccProtected) == AccProtected
}

/*IsStatic ...*/
func IsStatic(value uint8) bool {
	return (value & AccStatic) == AccStatic
}

/*IsFinal ...*/
func IsFinal(value uint8) bool {
	return (value & AccFinal) == AccFinal
}

/*IsAbstract ...*/
func IsAbstract(value uint8) bool {
	return (value & AccAbstract) == AccAbstract
}

/*IsInit ...*/
func IsInit(value uint8) bool {
	return (value & AccInit) == AccInit
}
