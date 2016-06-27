package codec

type H264FUAFragmenter struct {
	saveNumTruncatedBytes        uint
	maxOutputPacketSize          uint
	numValidDataBytes            uint
	inputBufferSize              uint
	curDataOffset                uint
	inputBuffer                  []byte
	lastFragmentCompletedNALUnit bool
}

func NewH264FUAFragmenter(inputSource IFramedSource, inputBufferMax uint) *H264FUAFragmenter {
	fragment := new(H264FUAFragmenter)
	fragment.numValidDataBytes = 1
	fragment.inputBufferSize = inputBufferMax + 1
	fragment.inputBuffer = make([]byte, fragment.inputBufferSize)
	fragment.InitFramedFilter(inputSource)
	fragment.InitFramedSource(fragment)
	return fragment
}

func (this *H264FUAFragmenter) doGetNextFrame() {
	fmt.Println(fmt.Sprintf("H264FUAFragmenter::doGetNextFrame -> %p", this.inputSource))
	if this.numValidDataBytes == 1 {
		this.inputSource.getNextFrame(this.buffTo, this.maxSize, this.afterGettingFunc, this.onCloseFunc)
	} else {
		if this.maxSize < this.maxOutputPacketSize {
		} else {
			this.maxSize = this.maxOutputPacketSize
		}

		if this.curDataOffset == 1 {
			if this.numValidDataBytes-1 <= this.maxSize { // case 1
				this.buffTo = this.inputBuffer[1 : this.numValidDataBytes-1]
				this.frameSize = this.numValidDataBytes - 1
				this.curDataOffset = this.numValidDataBytes
			} else { // case 2
				// We need to send the NAL unit data as FU-A packets.  Deliver the first
				// packet now.  Note that we add FU indicator and FU header bytes to the front
				// of the packet (reusing the existing NAL header byte for the FU header).
				this.inputBuffer[0] = (this.inputBuffer[1] & 0xE0) | 28   // FU indicator
				this.inputBuffer[1] = 0x80 | (this.inputBuffer[1] & 0x1F) // FU header (with S bit)
				this.buffTo = this.inputBuffer[:this.maxSize]
				this.frameSize = this.maxSize
				this.curDataOffset += this.maxSize - 1
				this.lastFragmentCompletedNALUnit = false
			}
		} else {
			this.inputBuffer[this.curDataOffset-2] = this.inputBuffer[0]         // FU indicator
			this.inputBuffer[this.curDataOffset-1] = this.inputBuffer[1] &^ 0x80 // FU header (no S bit)
			numBytesToSend := 2 + this.numValidDataBytes - this.curDataOffset
			if numBytesToSend > this.maxSize {
				// We can't send all of the remaining data this time:
				numBytesToSend = this.maxSize
				this.lastFragmentCompletedNALUnit = false
			} else {
				// This is the last fragment:
				this.inputBuffer[this.curDataOffset-1] |= 0x40 // set the E bit in the FU header
				this.numTruncatedBytes = this.saveNumTruncatedBytes
			}
			this.buffTo = this.inputBuffer[this.curDataOffset-2 : numBytesToSend]
			this.frameSize = numBytesToSend
			this.curDataOffset += numBytesToSend - 2
		}
	}

	if this.curDataOffset >= this.numValidDataBytes {
		// We're done with this data.  Reset the pointers for receiving new data:
		this.numValidDataBytes = 1
		this.curDataOffset = 1
	}

	// Complete delivery to the client:
	this.inputSource.afterGetting()
}

func (this *H264FUAFragmenter) afterGettingFrame(frameSize uint) {
	this.numValidDataBytes += frameSize
}
