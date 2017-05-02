package javacard

/**
 * The PackedBoolean manages booleans in volatile storage
 * space efficiently.
 */
type PackedBoolean struct {
	container []int8
	nextID    int8
}

/*
 * Constructor. Allocates an instance of PackedBoolean
 */
func initPackedBoolean(maxbytes int8) *PackedBoolean {
	pbl := &PackedBoolean{}
	//pbl.container = framework.MakeTransientArray(int16(maxbytes), framework.ClearOnReset, 1)
	pbl.nextID = 0
	return pbl
}

/**
 * Allocates a new boolean and returns the associated int8 identifier.
 */
func (pbl *PackedBoolean) Allocate() int8 {
	pbl.nextID++
	return pbl.nextID
}

/*
 * Returns the state of identified boolean.
 */
func (pbl *PackedBoolean) Get(identifier int8) bool {
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
func (pbl *PackedBoolean) access(identifier int8, typ int8) bool {
	bOff := int8(identifier >> 3)
	bitNum := int8(identifier & 0x7)
	interm := uint(bitNum)
	bitMask := int8(int16(0x80) >> interm)
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
func (pbl *PackedBoolean) Set(identifier int8) {
	pbl.access(identifier, 1)
}

/**
 * Resets the state of the identified boolean to false.
 *
 * @param boolean
 *            identifier
 */
func (pbl *PackedBoolean) Reset(identifier int8) {
	pbl.access(identifier, int8(-1))
}
