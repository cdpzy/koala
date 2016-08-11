package h264

import "github.com/doublemo/koala/helper"

// PicTiming Picture timing SEI message syntax
type PicTiming struct {
	CpbRemovalDelay    uint
	DpbOutputDelay     uint
	PicStruct          uint
	ClockTimestampFlag []uint
	CtType             uint
	NuitFieldBasedFlag uint
	CountingType       uint
	FullTimestampFlag  uint
	DiscontinuityFlag  uint
	CntDroppedFlag     uint
	NFrames            uint
	SecondsValue       uint
	MinutesValue       uint
	HoursValue         uint
	SecondsFlag        uint
	MinutesFlag        uint
	HoursFlag          uint
	TimeOffset         int
}

// SupplementalEnhancementInformation Supplemental enhancement information (SEI)
type SupplementalEnhancementInformation struct {
	PicTiming                   *PicTiming
	CpbDpbDelaysPresentFlag     bool
	CpbRemovalDelayLengthMinus1 uint
	DpbOutputDelayLengthMinus1  uint
	PicStructPresentFlag        bool
	TimeOffsetLength            uint
}

// NewSupplementalEnhancementInformation Supplemental enhancement information (SEI)
func NewSupplementalEnhancementInformation(vuiParameters *VUIParameters, hrdParameters *HRDParameters) *SupplementalEnhancementInformation {
	sei := new(SupplementalEnhancementInformation)
	sei.CpbRemovalDelayLengthMinus1 = 23
	sei.DpbOutputDelayLengthMinus1 = 23
	if vuiParameters != nil {
		sei.CpbDpbDelaysPresentFlag = vuiParameters.NalHrdParametersPresentFlag > 0 || vuiParameters.VclHrdParametersPresentFlag > 0
		sei.PicStructPresentFlag = vuiParameters.PicStructPresentFlag > 0
	}

	if hrdParameters != nil {
		sei.CpbRemovalDelayLengthMinus1 = hrdParameters.CpbRemovalDelayLengthMinus1
		sei.DpbOutputDelayLengthMinus1 = hrdParameters.DpbOutputDelayLengthMinus1
		sei.TimeOffsetLength = hrdParameters.TimeOffsetLength
	}
	return sei
}

// ParseSEIMessageBytes Supplemental enhancement information message syntax
// sei_message( )
func (supplementalEnhancementInformation *SupplementalEnhancementInformation) ParseSEIMessageBytes(b []byte) error {
	var (
		payloadType uint
		payloadSize uint
		i           int
	)

	seiSize := len(b)
	for i < seiSize {
		payloadType = 0
		payloadSize = 0
	LoopT:
		payloadType += uint(b[i])
		i++
		if b[i] == 0xFF && i < seiSize {
			goto LoopT
		}

		if i > seiSize {
			break
		}

	LoopS:
		payloadSize += uint(b[i])
		i++
		if b[i] == 0xFF && i < seiSize {
			goto LoopS
		}

		if i > seiSize {
			break
		}

		supplementalEnhancementInformation.setPayload(payloadType, payloadSize, b[i:])
		i += int(payloadSize)
	}
	return nil
}

func (supplementalEnhancementInformation *SupplementalEnhancementInformation) setPayload(payloadType, payloadSize uint, playload []byte) {
	switch payloadType {
	case 1:
		supplementalEnhancementInformation.parsePicTiming(payloadSize, playload)
	}
}

func (supplementalEnhancementInformation *SupplementalEnhancementInformation) parsePicTiming(payloadSize uint, playload []byte) {
	eg := helper.NewExpGolombReader(playload)
	picTiming := &PicTiming{}
	if supplementalEnhancementInformation.CpbDpbDelaysPresentFlag {
		picTiming.CpbRemovalDelay = handleParseError(eg.ReadBits(int(supplementalEnhancementInformation.CpbRemovalDelayLengthMinus1 + 1)))
		picTiming.DpbOutputDelay = handleParseError(eg.ReadBits(int(supplementalEnhancementInformation.DpbOutputDelayLengthMinus1 + 1)))
	}

	if supplementalEnhancementInformation.PicStructPresentFlag {
		picTiming.PicStruct = handleParseError(eg.ReadBits(4))
		numClockTS := 0
		switch picTiming.PicStruct {
		case 0:
			numClockTS = 1
		case 1:
			numClockTS = 1
		case 2:
			numClockTS = 1
		case 3:
			numClockTS = 2
		case 4:
			numClockTS = 2
		case 5:
			numClockTS = 3
		case 6:
			numClockTS = 3
		case 7:
			numClockTS = 2
		case 8:
			numClockTS = 3
		}

		for i := 0; i < numClockTS; i++ {
			clockTimestampFlag := handleParseError(eg.ReadBit())
			picTiming.ClockTimestampFlag = append(picTiming.ClockTimestampFlag, clockTimestampFlag)
			if clockTimestampFlag > 0 {
				picTiming.CtType = handleParseError(eg.ReadBits(2))
				picTiming.NuitFieldBasedFlag = handleParseError(eg.ReadBit())
				picTiming.CountingType = handleParseError(eg.ReadBits(5))
				picTiming.FullTimestampFlag = handleParseError(eg.ReadBit())
				picTiming.DiscontinuityFlag = handleParseError(eg.ReadBit())
				picTiming.CntDroppedFlag = handleParseError(eg.ReadBit())
				picTiming.NFrames = handleParseError(eg.ReadBits(8))

				if picTiming.FullTimestampFlag > 0 {
					picTiming.SecondsValue = handleParseError(eg.ReadBits(6))
					picTiming.MinutesValue = handleParseError(eg.ReadBits(6))
					picTiming.HoursValue = handleParseError(eg.ReadBits(5))
				} else {
					picTiming.SecondsFlag = handleParseError(eg.ReadBit())
					if picTiming.SecondsFlag > 0 {
						picTiming.SecondsValue = handleParseError(eg.ReadBits(6))
						picTiming.MinutesFlag = handleParseError(eg.ReadBit())
						if picTiming.MinutesFlag > 0 {
							picTiming.MinutesValue = handleParseError(eg.ReadBits(6))
							picTiming.HoursFlag = handleParseError(eg.ReadBit())
							if picTiming.HoursFlag > 0 {
								picTiming.HoursValue = handleParseError(eg.ReadBits(6))
							}
						}
					}

					// if time_offset_length > 0
					// time_offset
					if supplementalEnhancementInformation.TimeOffsetLength > 0 {
						picTiming.TimeOffset = handleParseSEError(eg.ReadIBits(int(supplementalEnhancementInformation.TimeOffsetLength)))
					}
				}
			}

		}
	}

	supplementalEnhancementInformation.PicTiming = picTiming
}
