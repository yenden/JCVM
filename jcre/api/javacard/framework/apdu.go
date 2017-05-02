package framework

import (
	"JCVM/jcre/api/com/sun/javacard"
	"errors"
	"log"
)

type APDU struct {
	/**
	 * The APDU will use the buffer byte[] to store data for input and output.
	 */
	buffer        []byte
	scratchBuffer []byte
	// Define two state variable buffers: short-length and byte-length.
	ramByteVars  []byte
	ramShortVars []int16
}

const (
	// State variable 16 bits in length
	LE                 = int8(0)
	LR                 = int8(LE + 1)
	LC                 = int8(LR + 1)
	PreReadLength      = int8(LC + 1)
	IncommingLength    = int8(PreReadLength + 1)
	RamShortVarsLength = int8(IncommingLength + 1)
	// State variables 8 bits in length
	CurrentState             = byte(0)
	LogicalChn               = byte(CurrentState + 1)
	RambyteVarsLength        = byte(LogicalChn + 1)
	StateOutgoing            = 3
	StateOutgoingLengthKnown = 4
	MaxLE                    = int16(256)
	BufferSize               = 133
	InvalidGetResponse       = int(0xC006)
	// procedure byte type constants
	AckNone              = 0
	AckIns               = 1
	AckNotIns            = 2
	StatePartialOutgoing = 5
	StateFullOutgoing    = 6
	StatePartialIncoming = 1
	StateFullIncoming    = 2
)

var (
	incomingFlag, sendInProgressFlag, outgoingFlag       int8
	outgoingLenSetFlag, noChainingFlag, extendedAPDUFlag int8
	envelopeFlag, extendedSupportFlag, determinedLE      int8
	noGetResponseFlag                                    int8
	thePackedBoolean                                     *javacard.PackedBoolean
)

func initAPDU() *APDU {
	apdu := &APDU{}
	/*	apdu.buffer = impl.InitAPDUBuffer()
		apdu.scratchBuffer = impl.T0InitScratchAPDUBuffer()
		apdu.ramByteVars = MakeTransientArray(RambyteVarsLength, ClearOnReset, 1)
		apdu.ramShortVars = MakeTransientArray(RambyteVarsLength, ClearOnReset, 2)
		thePackedBoolean = PrivAccess.getPackedBoolean()*/
	incomingFlag = thePackedBoolean.Allocate()
	sendInProgressFlag = thePackedBoolean.Allocate()
	outgoingFlag = thePackedBoolean.Allocate()
	outgoingLenSetFlag = thePackedBoolean.Allocate()
	noChainingFlag = thePackedBoolean.Allocate()
	noGetResponseFlag = thePackedBoolean.Allocate()
	extendedAPDUFlag = thePackedBoolean.Allocate()
	envelopeFlag = thePackedBoolean.Allocate()
	extendedSupportFlag = thePackedBoolean.Allocate()
	determinedLE = thePackedBoolean.Allocate()

	return apdu
}

func (apdu *APDU) getBuffer() []byte {
	return apdu.buffer
}
func (apdu *APDU) checkErrorState() error {
	if apdu.getCurrentState() < 0 { // error state
		return errors.New("error: APDU in error state")
	}
	return nil
}
func (apdu *APDU) getCurrentState() byte {
	return apdu.ramByteVars[CurrentState]
}
func (apdu *APDU) setCurrentState(data byte) {
	apdu.ramByteVars[CurrentState] = data
}

func getOutgoingFlag() bool {
	return thePackedBoolean.Get(outgoingFlag)
}
func setOutgoingFlag() {
	thePackedBoolean.Set(outgoingFlag)
}

func (apdu *APDU) setOutgoing() int16 {
	err := apdu.checkErrorState()
	if err == nil {
		log.Fatal(err)
	}
	// If we've previously called this method, then throw an exception
	if getOutgoingFlag() {
		log.Fatal("Error setOutgoing, Illegal use of APDU")
	}
	setOutgoingFlag()
	apdu.setCurrentState(StateOutgoing)
	return apdu.getLe()
}

func (apdu *APDU) getLe() int16 {
	if apdu.ramShortVars[LE] == int16(0) {
		return MaxLE
	}
	return apdu.ramShortVars[LE]
}
func (apdu *APDU) getLr() int16 {
	return apdu.ramShortVars[LR]
}

func (apdu *APDU) setLr(data int16) {
	apdu.ramShortVars[LR] = data
}

func getOutgoingLenSetFlag() bool {
	return thePackedBoolean.Get(outgoingLenSetFlag)
}
func setOutgoingLenSetFlag() {
	thePackedBoolean.Set(outgoingLenSetFlag)
}
func (apdu *APDU) SetOutgoingAndSend(bOff int16, len int16) {
	apdu.setOutgoing()
	err := apdu.setOutgoingLength(len)
	if err != nil {
		log.Fatal(err)
	}
	apdu.SendBytes(bOff, len)
}
func (apdu *APDU) setOutgoingLength(len int16) error {
	if !getOutgoingFlag() {
		return errors.New("ADPU Illegal use in setOutgoingLength")
	}
	// if we've previously called this method, then throw an exception
	if getOutgoingLenSetFlag() {
		return errors.New("ADPU Illegal use in setOutgoingLength")
	}
	// If the outbound length is being set to more than 256 and the applet
	// does not implement the ExtendedAPDU interface, throw an exception.
	if len > MaxLE {
		return errors.New("ADPU Bad Length use in setOutgoingLength")
	}
	if len < 0 {
		return errors.New("ADPU Bad Length use in setOutgoingLength")
	}
	setOutgoingLenSetFlag()
	apdu.setCurrentState(StateOutgoingLengthKnown)
	apdu.setLr(len)
	return nil
}
func getNoGetResponseFlag() bool {
	return thePackedBoolean.Get(noGetResponseFlag)
}
func getIncomingFlag() bool {
	return thePackedBoolean.Get(incomingFlag)
}
func getSendInProgressFlag() bool {
	return thePackedBoolean.Get(sendInProgressFlag)
}
func (apdu *APDU) getCLAChannel() byte {
	return apdu.getLogicalChannel()

}
func (apdu *APDU) getLogicalChannel() byte {
	return apdu.ramByteVars[LogicalChn]
}
func (apdu *APDU) send61xx(len int16) (int16, error) {
	//originChannel := apdu.getCLAChannel()
	expLen := len
	for ok := true; ok; ok = (expLen > len) { //do... while
		// Set SW1SW2 as 61xx. For Case 2E and 4E, set to 6100 if more
		// than 256 bytes are to be sent inthe outbound direction.
		if len >= MaxLE {
			//NativeMethods.t0SetStatus((ISO7816.SW_BYTES_REMAINING_00))
		} else {
			// NativeMethods.t0SetStatus((short) (ISO7816.SW_BYTES_REMAINING_00 + (short) (len & (short) 0x00FF)))
		}
		newLen := 0 //NativeMethods.t0SndGetResponse(originChannel)
		if newLen == InvalidGetResponse {
			// Get Response not received
			setNoGetResponseFlag()
			return 0, errors.New("Error get response not received")
		} else if newLen > 0 {
			apdu.setLe(int16(newLen))
			expLen = apdu.getLe()
		} else {
			return 0, errors.New("Error input/output in send61xx")
		}

	}
	resetSendInProgressFlag()
	return expLen, nil

}
func setNoGetResponseFlag() {
	thePackedBoolean.Set(noGetResponseFlag)
}
func (apdu *APDU) setLe(data int16) {
	apdu.ramShortVars[LE] = int16(data & int16(0x7FFF))
}

func resetSendInProgressFlag() {
	thePackedBoolean.Reset(sendInProgressFlag)
}
func setSendInProgressFlag() {
	thePackedBoolean.Set(sendInProgressFlag)
}

func (apdu *APDU) SendBytes(bOff int16, len int16) error {
	if bOff < 0 || len < 0 || int(bOff+len) > BufferSize {
		return errors.New("ADPU BufferBound error in sendBytes")
	}
	if !getOutgoingLenSetFlag() || getNoGetResponseFlag() {
		return errors.New("ADPU Illegal use in SendBytes")
	}
	if len == 0 {
		return nil
	}
	Lr := apdu.getLr()
	if len > Lr {
		return errors.New("ADPU Illegal use in SendBytes")
	}
	Le := apdu.getLe()
	for len > 0 {
		temp := len
		var err error
		// Need to force GET RESPONSE for Case 4 & for partial blocks
		if len != Lr || getIncomingFlag() || Lr != Le || getSendInProgressFlag() {
			temp, err = apdu.send61xx(len) // resets sendInProgressFlag
			if err != nil {
				log.Fatal(err)
			}
			resetIncomingFlag() // no more incoming->outgoing
			// switch.
		}
		result := 0 //NativeMethods.t0SndData(buffer, bOff, temp, AckIns )
		setSendInProgressFlag()
		if result != 0 {
			return errors.New("Error: input/output errors in sendbytes")
		}
		bOff += temp
		len -= temp
		Lr -= temp
		Le = Lr

	}
	if Lr == 0 {
		apdu.setCurrentState(StateFullOutgoing)
	} else {
		apdu.setCurrentState(StatePartialOutgoing)
	}
	apdu.setLe(Le) // update RAM copy of Le, the expected count remaining
	apdu.setLr(Lr) // update RAM copy of Lr, the response count remaining
	return nil
}
func resetIncomingFlag() {
	thePackedBoolean.Reset(incomingFlag)
}
func (apdu *APDU) SetIncomingAndReceive() int16 {
	err := apdu.setIncoming()
	if err != nil {
		log.Fatal(err)
	}
	return receiveBytes(OffsetCData)
}
func (apdu *APDU) setIncoming() error {
	// if Java Card runtime environment has undone a previous
	// setIncomingAndReceive ignore
	if apdu.getPreReadLength() != 0 {
		return nil
	}
	// if we've previously called this or setOutgoing() method,
	// then throw an exception
	if getIncomingFlag() || getOutgoingFlag() {
		return errors.New("Error illegal use in setIncoming")
	}
	setIncomingFlag()  // indicate that this method has been called
	Lc := apdu.getLe() // what we stored in Le was really Lc
	// If the incoming data is greater than 256
	if Lc > 256 {
		return errors.New("Error wrong status word length")
	}
	apdu.setLc(Lc)
	apdu.setIL(Lc)
	apdu.setLe(0) // in T=0, the real Le is now unknown (assume 256)
	return nil

}
func (apdu *APDU) getPreReadLength() int16 {
	return apdu.ramShortVars[PreReadLength]
}
func setIncomingFlag() {
	thePackedBoolean.Set(incomingFlag)
}
func (apdu *APDU) setLc(data int16) {
	apdu.ramShortVars[LC] = data
}
func (apdu *APDU) setIL(data int16) {
	apdu.ramShortVars[IncommingLength] = data
}

func (apdu *APDU) ReceiveBytes(bOff int16) (int16, error) {
	// Main receive method. This method will receive data from the CAD,
	// or it will issue a positive reply toretrieve the next APDU command.
	// Only APDUs case 3 and 4 are expected to call this method.
	if !getIncomingFlag() || getOutgoingFlag() {
		return 0, errors.New("error illegal use in receiveBytes")
	}
	Lc := apdu.getLc() & int16(0x7FFF)
	if bOff < 0 {
		return 0, errors.New("error buffer bound in receiveBytes")
	}
	pre := apdu.getPreReadLength() & int16(0x7FFF)
	if pre != 0 {
		apdu.setPreReadLength(0)
		if Lc == 0 {
			apdu.setCurrentState(StateFullIncoming)
		} else {
			apdu.setCurrentState(StatePartialIncoming)
		}
		return pre, nil
	}
	Len := int16(0) //one time :NativeMethods.t0RcvData(bOff)
	if Lc != 0 {
		Len = 0 //NativeMethods.t0RcvData(bOff)
		if Len < 0 {
			return 0, errors.New("Error input/output in receivebytes")
		}
		// Move from scratch buffer to APDU buffer
		//NativeMethods.t0CopyToAPDUBuffer(bOff, len)
		Lc = Lc - Len
		if Lc < 0 {
			if Lc == -1 {
				// Partially received LE
				Lc = 0
				Len = Len - 1
				setDeterminedLEFlag()
			} else {
				Lc = 0
				Len = Len - 2
				setDeterminedLEFlag()
			}
		}
		apdu.setLc(Lc) // update RAM copy of Lc, the count remaining
		if Lc == 0 {
			apdu.setCurrentState(StateFullIncoming)
		} else {
			apdu.setCurrentState(StatePartialIncoming)
		}
		return Len, nil
	}
	apdu.setCurrentState(StateFullIncoming)
	return 0, nil
}

func (apdu *APDU) getLc() int16 {
	return apdu.ramShortVars[LC]
}
func (apdu *APDU) setPreReadLength(data int16) {
	apdu.ramShortVars[PreReadLength] = data
}
func setDeterminedLEFlag() {
	thePackedBoolean.Set(determinedLE)
}
func getDeterminedLEFlag() bool {
	return thePackedBoolean.Get(determinedLE)
}
