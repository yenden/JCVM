package api

const (
	// Mnemonics for the SW1,SW2 error codes

	/**
	 * Response status : No Error = (short)0x9000
	 */
	SwNoError = uint16(0x9000)

	/**
	 * Response status : Wrong length = 0x6700
	 */
	SwWrongLength = uint16(0x6700)

	/**
	 * Response status : Security condition not satisfied = 0x6982
	 */
	SwSecurityStatusNotSatisfied = uint16(0x6982)

	/**
	 * Response status : Data invalid = 0x6984
	 */
	SwDataInvalid = uint16(0x6984)

	/**
	 * Response status : Conditions of use not satisfied = 0x6985
	 */
	SwConditionsNotStatisfied = uint16(0x6985)

	/**
	 * Response status : Applet selection failed = 0x6999;
	 */
	SwAppletSelectFailed = uint16(0x6999)

	/**
	 * Response status : Wrong data = 0x6A80
	 */
	SwWrongData = uint16(0x6A80)

	/**
	 * Response status : Correct Expected Length (Le) = 0x6C00
	 */
	SwCorrectLength00 = uint16(0x6C00)

	/**
	 * Response status : INS value not supported = 0x6D00
	 */
	SwInsNotSupported = uint16(0x6D00)

	/**
	 * Response status : CLA value not supported = 0x6E00
	 */
	SwClaNotSUpported = uint16(0x6E00)

	/**
	 * Response status : No precise diagnosis = 0x6F00
	 */
	SwUnknown = uint16(0x6F00)

	// Offsets into APDU header information

	/**
	 * APDU header offset : CLA = 0
	 */
	OffsetCLA = byte(0)

	/**
	 * APDU header offset : INS = 1
	 */
	OffsetIns = byte(1)

	/**
	 * APDU header offset : P1 = 2
	 */
	OfssetP1 = byte(2)

	/**
	 * APDU header offset : P2 = 3
	 */
	OfssetP2 = byte(3)

	/**
	 * APDU header offset : LC = 4
	 */
	OffsetLC = byte(4)

	/**
	 * APDU command data offset : CDATA = 5
	 */
	OffsetCData = byte(5)

	/**
	 * APDU command CLA : ISO 7816 = 0x00
	 */
	ClaISO7816 = byte(0x00)

	/**
	 * APDU command INS : SELECT = 0xA4
	 */
	InsSelect  = byte(0xA4)
	InsInstall = byte(0x01)

	/**
	 * APDU command INS : EXTERNAL AUTHENTICATE = 0x82
	 */
	InsExternalAuthenticate = byte(0x82)

	// PINException reason codes
	/**
	 * This reason code is used to indicate that one or more input parameters is
	 * out of allowed bounds.
	 */
	IllegalValue = uint16(1)
	/**
	 * This reason code is used to indicate a method has been invoked at an illegal
	 * or inappropriate time.
	 */
	IllegalState = uint16(2)
)
