package framework

import (
	"JCVM/jcre/api/com/sun/javacard/impl"
	"JCVM/jcre/api/share"
	"errors"
	"log"
)

type APDU struct {
	/**
	 * The APDU will use the buffer byte[] to store data for input and output.
	 */
	buffer []byte
	//scrach buffer to receive incoming data
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
	CurrentState      = byte(0)
	LogicalChn        = byte(CurrentState + 1)
	RambyteVarsLength = byte(LogicalChn + 1)
	//apdu states
	StateOutgoing            = byte(3)
	StateOutgoingLengthKnown = byte(4)
	StatePartialOutgoing     = 5
	StateFullOutgoing        = 6
	StatePartialIncoming     = 1
	StateFullIncoming        = 2
	StateInitial             = byte(0)

	MaxLE              = int16(128)
	BufferSize         = 128
	InvalidGetResponse = int(0xC006)
	// procedure byte type constants
	AckNone     = 0
	AckIns      = 1
	AckNotIns   = 2
	EnvelopeIns = byte(0xC2)

	MaxXFerDataLength = int16(BufferSize - 5)

	ApduIsoClassMask     = byte(0x80)
	ApduTypeMask         = byte(0x4)
	LogicalChnMastType16 = byte(0x0F)
	ApduInvalidType4Msk  = byte(0x20)
	LogicalChnMaskType4  = byte(0x03)
)

var (
	incomingFlag, sendInProgressFlag, outgoingFlag byte
	outgoingLenSetFlag, noGetResponseFlag          byte
	thePackedBoolean                               *impl.PackedBoolean
)

func initAPDU() *APDU {
	apdu := &APDU{}
	/*	apdu.buffer = impl.InitAPDUBuffer()
		apdu.scratchBuffer = impl.T0InitScratchAPDUBuffer()*/
	apdu.ramByteVars = share.MakeTransientByteArray(int16(RambyteVarsLength), share.ClearOnReset)
	apdu.ramShortVars = share.MakeTransientShortArray(int16(RambyteVarsLength), share.ClearOnReset)
	thePackedBoolean.GetPackedBoolean()
	incomingFlag = thePackedBoolean.Allocate()
	sendInProgressFlag = thePackedBoolean.Allocate()
	outgoingFlag = thePackedBoolean.Allocate()
	outgoingLenSetFlag = thePackedBoolean.Allocate()
	noGetResponseFlag = thePackedBoolean.Allocate()
	//extendedAPDUFlag = thePackedBoolean.Allocate()
	//envelopeFlag = thePackedBoolean.Allocate()
	//extendedSupportFlag = thePackedBoolean.Allocate()
	//	determinedLE = thePackedBoolean.Allocate()

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
	dataReceived, err := apdu.ReceiveBytes(OffsetCData)
	if err != nil {
		log.Fatal(err)
	}
	return dataReceived
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
	//var err error
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
	Len := int16(0)
	if Lc != 0 {
		//	if !getEnvelopeFlag() {
		// Non-envelope case - call t0RcvData
		Len = 0 //NativeMethods.t0RcvData(bOff)
		if Len < 0 {
			return 0, errors.New("Error input/output in receivebytes")
		}
		/*	} else {
			// Envelope case - call t0RcvData first
			Len = 0 // NativeMethods.t0RcvData(bOff);

			if Len < 0 {
				return 0, errors.New("Error input/output in receivebytes")
			}
			if Len == 0 {
				// Envelope case - retrieve next ENVELOPE APDU if needed
				// NativeMethods.t0SetStatus((short) 0x9000);
				//NativeMethods.t0SndStatusRcvCommand()
				// verify that data is properly received in an
				// Envelope command
				Len, err = apdu.processEnvelopeData(bOff)
				if err != nil {
					log.Fatal(err)
				}
			}
		}*/

		// Move from scratch buffer to APDU buffer
		//NativeMethods.t0CopyToAPDUBuffer(bOff, len)
		Lc = Lc - Len
		/*if Lc < 0 {
			if Lc == -1 {
				// Partially received LE
				Lc = 0
				Len = Len - 1
				setDeterminedLEFlag()
			} else {
				Lc = 0
				Len = Len - 2
				setDeterminedLEFlag()
				resetEnvelopeFlag()
			}
		}*/
		apdu.setLc(Lc) // update RAM copy of Lc, the count remaining
		if Lc == 0 {
			apdu.setCurrentState(StateFullIncoming)
		} else {
			apdu.setCurrentState(StatePartialIncoming)
		}
		return Len, nil
	} /*else {
		// This is the case where LE is missing one byte
		if getEnvelopeFlag() && getDeterminedLEFlag() {
			Len = 0 //NativeMethods.t0RcvData(bOff);
			if Len < 0 {
				return 0, errors.New("Error input/output in receivebytes")
			}

			if Len == 0 {
				// Envelope case
				//  Util.arrayFillNonAtomic(buffer, (short) 0, BUFFERSIZE, (byte) 0);
				//   NativeMethods.t0SetStatus((short) 0x9000);
				//     NativeMethods.t0SndStatusRcvCommand();
				// verify that data is properly received in an
				// Envelope command if it is not...
				Len, err = apdu.processEnvelopeData(bOff)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Copy data to APDU buffer
			// NativeMethods.t0CopyToAPDUBuffer(bOff, len);
			// Process that one byte LE
			// comment out for bugfix No LE in T=0 Extended TPDU
			// setLe((short) (getLe() | ((buffer[(short) (bOff + len - (short) 1)] & (short) 0x00FF))));
			resetEnvelopeFlag()
		}
	}*/
	apdu.setCurrentState(StateFullIncoming)
	return 0, nil
}

func (apdu *APDU) getLc() int16 {
	return apdu.ramShortVars[LC]
}
func (apdu *APDU) setPreReadLength(data int16) {
	apdu.ramShortVars[PreReadLength] = data
} /*
func setDeterminedLEFlag() {
	thePackedBoolean.Set(determinedLE)
}
func getDeterminedLEFlag() bool {
	return thePackedBoolean.Get(determinedLE)
}
func getEnvelopeFlag() bool {
	return thePackedBoolean.Get(envelopeFlag)
}*/
/*
func (apdu *APDU) processEnvelopeData(bOffset int16) (int16, error) {
	// This method is called when an ENVELOPE command is expected
	// in order to retrieve its payload.
	// If this is not an ENVELOPE command... send I/O Error.
	if (isISOInterindustryCLA(apdu.scratchBuffer)) && (apdu.scratchBuffer[OffsetIns] == EnvelopeIns) {
		// retrieve the data and put it where it is requested
		Len := 0 //NativeMethods.t0RcvData(bOffset)
		if Len < 0 {
			return 0, errors.New("Error input/output in processEnvelopeData")
		}
		return int16(Len), nil
	} else {
		return 0, errors.New("Error input/output in processEnvelopeData")
	}
	return -1, nil

}
func isISOInterindustryCLA(aBuffer []byte) bool {
	theCLAType := (aBuffer[OffsetCLA] & ApduIsoClassMask)
	if theCLAType == 0x00 {
		return true
	}
	return false
}*/
/*
func resetEnvelopeFlag() {
	thePackedBoolean.Reset(envelopeFlag)
}*/
func (apdu *APDU) SendBytesLong(outData []byte, bOff int16, Len int16) {
	impl.CheckArrayArgs(outData, bOff, Len)
	sendLength := int16(len(apdu.buffer))
	for Len > 0 {
		if Len < sendLength {
			sendLength = Len
		}
		ArrayCopy(outData, bOff, apdu.buffer, 0, sendLength)
		apdu.SendBytes(0, sendLength)
		Len -= sendLength
		bOff += sendLength
	}
}
func (apdu *APDU) complete(status int16) error {
	// Zero out APDU buffer
	var result int16
	ArrayfillNonAtomic(apdu.buffer, 0, BufferSize, 0)
	if !getNoGetResponseFlag() && getSendInProgressFlag() {
		Le := apdu.getLe()
		sendLen := MaxXFerDataLength
		for Le > 0 {
			if Le < MaxXFerDataLength {
				sendLen = Le
			}
			result = 0 //NativeMethods.t0SndData(buffer, (byte) 0, sendLen, ACK_NONE)
			if result != 0 {
				return errors.New("imput/output error in complete method")
			}
			Le -= sendLen
		}
	}
	apdu.buffer[0] = byte(status >> 8)
	apdu.buffer[1] = byte(status)
	if status == 0 {
		result = 0 //NativeMethods.t0RcvCommand()
	} else {
		// NativeMethods.t0SetStatus(status)
		result = 0 //NativeMethods.t0SndStatusRcvCommand()
	}
	if result != 0 {
		return errors.New("imput/output error in complete method")
	}
	apdu.resetAPDU()
	apdu.preProcessAPDU()
	return nil
}
func (apdu *APDU) resetAPDU() {
	resetIncomingFlag()
	resetOutgoingFlag()
	resetOutgoingLenSetFlag()
	resetSendInProgressFlag()
	resetNoGetResponseFlag()
	apdu.setPreReadLength(0)

}
func resetNoGetResponseFlag() {
	thePackedBoolean.Reset(noGetResponseFlag)
}
func resetOutgoingFlag() {
	thePackedBoolean.Reset(outgoingFlag)
}
func resetOutgoingLenSetFlag() {
	thePackedBoolean.Reset(outgoingLenSetFlag)
}

func (apdu *APDU) preProcessAPDU() {
	// as described earlier, we assume case 1 or 2S, so Le=P3 and Lc=0
	apdu.setLe(int16(apdu.scratchBuffer[OffsetLC] & 0x00FF))
	apdu.setLc(0)
	apdu.setIL(int16(apdu.scratchBuffer[OffsetLC] & 0x00ff))
	apdu.resetCurrState()

	// set logical channel information as specified in the APDU CLA byte
	apdu.setLogicalChannel(apdu.getChannelInfo())
	apdu.setLr(0)
}

func (apdu *APDU) resetCurrState() {
	apdu.ramByteVars[CurrentState] = StateInitial
}
func (apdu *APDU) setLogicalChannel(data byte) {
	apdu.ramByteVars[LogicalChn] = data
}
func (apdu *APDU) getChannelInfo() byte {
	var theAPDUChannel byte
	if int16(apdu.scratchBuffer[OffsetCLA]&0x00FF) == 0x00FF {
		theAPDUChannel = 0
	} else if isType16CLA(apdu.scratchBuffer) {
		theAPDUChannel = byte(apdu.scratchBuffer[OffsetCLA] & LogicalChnMastType16)
		theAPDUChannel = byte(theAPDUChannel + 4)
	} else {
		if byte(apdu.scratchBuffer[OffsetCLA]&ApduInvalidType4Msk) != 0 {
			theAPDUChannel = 0
		} else {
			theAPDUChannel = byte(apdu.scratchBuffer[OffsetCLA] & LogicalChnMaskType4)
		}
	}

	return theAPDUChannel
}
func isType16CLA(aBuffer []byte) bool {
	var theCLAType byte
	theCLAType = byte(aBuffer[OffsetCLA] & ApduTypeMask)
	if theCLAType == 0x40 {
		return true
	}
	return false
}
