package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

// SequenceParameterSetExtensionRBSP  Sequence parameter set extension RBSP semantics
// T-REC-H.264-201602-I!!PDF-E.pdf 7.4
type SequenceParameterSetExtensionRBSP struct {
	SeqParameterSetID       uint
	AuxFormatIdc            uint
	BitDepthAuxMinus8       uint
	AlphaIncrFlag           uint
	AlphaOpaqueValue        uint
	AlphaTransparentValue   uint
	AdditionalExtensionFlag uint
}

// NewSequenceParameterSetExtensionRBSP  Sequence parameter set extension RBSP semantics
func NewSequenceParameterSetExtensionRBSP() *SequenceParameterSetExtensionRBSP {
	return new(SequenceParameterSetExtensionRBSP)
}

func (sequenceParameterSetExtensionRBSP *SequenceParameterSetExtensionRBSP) ParseBytes(b []byte) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			} else {
				panic(r)
			}
		}
	}()

	var v int

	eg := helper.NewExpGolombReader(b)
	sequenceParameterSetExtensionRBSP.SeqParameterSetID = handleParseError(eg.ReadUE())
	sequenceParameterSetExtensionRBSP.AuxFormatIdc = handleParseError(eg.ReadUE())

	if sequenceParameterSetExtensionRBSP.AuxFormatIdc != 0 {
		sequenceParameterSetExtensionRBSP.BitDepthAuxMinus8 = handleParseError(eg.ReadUE())
		sequenceParameterSetExtensionRBSP.AlphaIncrFlag = handleParseError(eg.ReadBit())
		v = int(sequenceParameterSetExtensionRBSP.BitDepthAuxMinus8 + 9)
		sequenceParameterSetExtensionRBSP.AlphaOpaqueValue = handleParseError(eg.ReadBits(v))      // u(v)
		sequenceParameterSetExtensionRBSP.AlphaTransparentValue = handleParseError(eg.ReadBits(v)) // u(v)
	}

	sequenceParameterSetExtensionRBSP.AdditionalExtensionFlag = handleParseError(eg.ReadBit())
	return
}
