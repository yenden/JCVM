package api

import (
	"errors"
)

type OwnerPIN struct {
	/**
	 * Try limit--the maximum number of times an incorrect PIN can be presented
	 * before the PIN is blocked. When the PIN is blocked, it cannot be
	 * validated even on valid PIN presentation.
	 */
	tryLimit byte
	/**
	 * max PIN size, the maximum length of PIN allowed
	 */
	maxPINSize byte
	/**
	 * PIN value
	 */
	pinValue []byte
	/**
	 * the current size of PIN array.
	 */
	pinSize byte
	/**
	 * validated flag, true if a valid PIN has been presented. This flag is
	 * reset on every card reset.
	 */
	flags     []bool // default null
	triesLeft []byte
}

const (
	validated = byte(0)
	numFlags  = byte(validated + 1)
)

func (own *OwnerPIN) GetValidatedFlag() bool {
	result := own.flags[validated]
	return result
}
func (own *OwnerPIN) SetValidatedFlag(value bool) {
	own.flags[validated] = value
}

func (own *OwnerPIN) ResetTriesRemaining() {
	ArrayfillNonAtomic(own.triesLeft, 0, 1, own.tryLimit)
}
func (own *OwnerPIN) DecrementTriesRemaining() {
	ArrayfillNonAtomic(own.triesLeft, 0, 1, byte(own.triesLeft[0]-1))
}

func (own *OwnerPIN) InitOwnerPIN(tryLimit byte, maxPINSize byte) (uint16, error) {
	if (tryLimit < 1) || (maxPINSize < 1) {
		return IllegalValue, errors.New("Illegal use of PIN")
	}
	own.pinValue = make([]byte, maxPINSize) // default value 0
	own.pinSize = maxPINSize                // default
	own.maxPINSize = maxPINSize
	own.tryLimit = tryLimit

	own.triesLeft = make([]byte, 1)
	own.ResetTriesRemaining()
	own.flags = make([]bool, numFlags)
	own.SetValidatedFlag(false)
	return 0, nil
}

func (own *OwnerPIN) GetTriesRemaining() byte {
	result := own.triesLeft[0]
	return result
}

func (own *OwnerPIN) Check(pin interface{}, offset int16, length byte) bool {
	noMoreTries := false
	own.SetValidatedFlag(false)
	if own.GetTriesRemaining() == 0 {
		noMoreTries = true
	} else {
		own.DecrementTriesRemaining()
	}
	if length > 0 {
		if (length != own.pinSize) || noMoreTries {

			return false
		}
	}
	switch pin.(type) {
	case []byte:
		n, err := ArrayCompare(pin.([]byte), offset, own.pinValue, 0, int16(length))
		if err != nil && n == 0 && length == own.pinSize {
			own.SetValidatedFlag(true)
			own.ResetTriesRemaining()

			return true
		}
	}
	return false
}

func (own *OwnerPIN) IsValidated() bool {
	result := own.GetValidatedFlag()
	return result
}

func (own *OwnerPIN) Reset() {
	if own.IsValidated() {
		own.ResetAndUnblock()
	}
}

func (own *OwnerPIN) Update(pin interface{}, offset int16, length byte) (uint16, error) {
	if length > own.maxPINSize {
		//nativeMethods.SetStatus(IllegalValue)
		return IllegalValue, errors.New("Illegal use of PIN")
	}
	switch pin.(type) {
	case []uint8:
		ArrayCopy(pin.([]uint8), offset, own.pinValue, 0, int16(length))
		own.pinSize = length
		own.triesLeft[0] = own.tryLimit
		own.SetValidatedFlag(false)
	}
	return 0, nil
}

func (own *OwnerPIN) ResetAndUnblock() {
	own.ResetTriesRemaining()
	own.SetValidatedFlag(false)
}
