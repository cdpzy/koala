package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

// HRDParameters HRD parameters syntax
type HRDParameters struct {
	CpbCntMinus1                       uint
	BitRateScal                        uint
	CpbSizeScale                       uint
	BitRateValueMinus1                 []uint
	CpbSizeValueMinus1                 []uint
	CbrFlag                            []uint
	InitialCpbRemovalDelayLengthMinus1 uint
	CpbRemovalDelayLengthMinus1        uint
	DpbOutputDelayLengthMinus1         uint
	TimeOffsetLength                   uint
}

// HRDParameters HRD parameters syntax
func NewHRDParameters() *HRDParameters {
	return new(HRDParameters)

}

// ParseExpGolombReader Parse
func (hrdParameters *HRDParameters) ParseExpGolombReader(eg *helper.ExpGolombReader) (err error) {
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

	var (
		SchedSelIdx uint
	)

	hrdParameters.CpbCntMinus1 = handleParseError(eg.ReadUE())
	hrdParameters.BitRateScal = handleParseError(eg.ReadBits(4))
	hrdParameters.CpbSizeScale = handleParseError(eg.ReadBits(4))

	for SchedSelIdx = 0; SchedSelIdx < hrdParameters.CpbCntMinus1; SchedSelIdx++ {
		hrdParameters.BitRateValueMinus1 = append(hrdParameters.BitRateValueMinus1, handleParseError(eg.ReadUE()))
		hrdParameters.CpbSizeValueMinus1 = append(hrdParameters.CpbSizeValueMinus1, handleParseError(eg.ReadUE()))
		hrdParameters.CbrFlag = append(hrdParameters.CbrFlag, handleParseError(eg.ReadBit()))
	}

	hrdParameters.InitialCpbRemovalDelayLengthMinus1 = handleParseError(eg.ReadBits(5))
	hrdParameters.CpbRemovalDelayLengthMinus1 = handleParseError(eg.ReadBits(5))
	hrdParameters.DpbOutputDelayLengthMinus1 = handleParseError(eg.ReadBits(5))
	hrdParameters.TimeOffsetLength = handleParseError(eg.ReadBits(5))
	return
}
