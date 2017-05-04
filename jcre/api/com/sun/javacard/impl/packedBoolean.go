package impl

import "JCVM/jcre/api/share"

/**
 * The PackedBoolean manages booleans in volatile storage
 * space efficiently.
 */
type PackedBoolean struct {
	container []byte
	nextID    byte
}

/*
 * Constructor. Allocates an instance of PackedBoolean
 */
func initPackedBoolean(maxbytes byte) *PackedBoolean {
	pbl := &PackedBoolean{}
	pbl.container = share.MakeTransientByteArray(int16(maxbytes), share.ClearOnReset)
	pbl.nextID = 0
	return pbl
}

/**
 * Allocates a new boolean and returns the associated int8 identifier.
 */
func (pbl *PackedBoolean) Allocate() byte {
	pbl.nextID++
	return pbl.nextID
}

/*
 * Returns the state of identified boolean.
 */
func (pbl *PackedBoolean) Get(identifier byte) bool {
	return pbl.access(identifier, 0)
}

/**
 * Changes the state of the identified boolean to the specified value or
 * simply queries
 *
 * @param identifier
 *            of boolean flag
 * @param type
 *            1 set, -1 reset, 0 no change
 * @return value boolean value of specified flag
 */
func (pbl *PackedBoolean) access(identifier byte, typ int8) bool {
	bOff := byte(identifier >> 3)
	bitNum := byte(identifier & 0x7)
	interm := uint(bitNum)
	bitMask := byte(int16(0x80) >> interm)
	switch typ {
	case 1:
		pbl.container[bOff] |= bitMask
	case -1:
		pbl.container[bOff] &= (^bitMask)

	}
	return ((pbl.container[bOff] & bitMask) != 0)

}

/**
 * Sets the state of the identified boolean to true.
 *
 * @param boolean
 *            identifier
 */
func (pbl *PackedBoolean) Set(identifier byte) {
	pbl.access(identifier, 1)
}

/**
 * Resets the state of the identified boolean to false.
 *
 * @param boolean
 *            identifier
 */
func (pbl *PackedBoolean) Reset(identifier byte) {
	pbl.access(identifier, -1)
}

const (
	NumberSystemBools = 24
)

func (pckB *PackedBoolean) GetPackedBoolean() {
	if pckB == nil {
		pckB = initPackedBoolean(byte((NumberSystemBools-1)>>3) + 1)
	}
}