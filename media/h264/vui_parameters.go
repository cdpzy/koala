package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

const (
	ASPECT_RATIO_IDC_Extended_SAR uint = 0xFF
)

// VUIParameters VUI parameters syntax
type VUIParameters struct {
	AspectRatioInfoPresentFlag         uint
	AspectRatioIdc                     uint
	SarWidth                           uint
	SarHeight                          uint
	OverscanInfoPresentFlag            uint
	OverscanAppropriateFlag            uint
	VideoSignalTypePresentFlag         uint
	VideoFormat                        uint
	VideoFullRangeFlag                 uint
	ColourDescriptionPresentFlag       uint
	ColourPrimaries                    uint
	TransferCharacteristics            uint
	MatrixCoefficients                 uint
	ChromaLocInfoPresentFlag           uint
	ChromaSampleLocTypeTopField        uint
	ChromaSampleLocTypeBottomField     uint
	TimingInfoPresentFlag              uint
	NumUnitsInTick                     uint
	TimeScale                          uint
	FixedFrameRateFlag                 uint
	NalHrdParametersPresentFlag        uint
	HRDParameters                      *HRDParameters
	VclHrdParametersPresentFlag        uint
	LowDelayHrdFlag                    uint
	PicStructPresentFlag               uint
	BitstreamRestrictionFlag           uint
	MotionVectorsOverPicBoundariesFlag uint
	MaxBytesPerPicDenom                uint
	MaxBitsPerMbDenom                  uint
	Log2MaxMvLengthHorizontal          uint
	Log2MaxMvLengthVertical            uint
	MaxNumReorderFrames                uint
	MaxDecFrameBuffering               uint
}

// NewVUIParameters VUI parameters syntax
func NewVUIParameters() *VUIParameters {
	return new(VUIParameters)
}

func (vuiParameters *VUIParameters) ParseExpGolombReader(eg *helper.ExpGolombReader) (err error) {
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

	vuiParameters.AspectRatioInfoPresentFlag = handleParseError(eg.ReadBit())

	if vuiParameters.AspectRatioInfoPresentFlag > 0 {
		vuiParameters.AspectRatioIdc = handleParseError(eg.ReadBits(8))

		if vuiParameters.AspectRatioIdc == ASPECT_RATIO_IDC_Extended_SAR {
			vuiParameters.SarWidth = handleParseError(eg.ReadBits(16))
			vuiParameters.SarHeight = handleParseError(eg.ReadBits(16))
		}
	}

	vuiParameters.OverscanInfoPresentFlag = handleParseError(eg.ReadBit())
	if vuiParameters.OverscanInfoPresentFlag > 0 {
		vuiParameters.OverscanAppropriateFlag = handleParseError(eg.ReadBit())
	}

	vuiParameters.VideoSignalTypePresentFlag = handleParseError(eg.ReadBit())

	if vuiParameters.VideoSignalTypePresentFlag > 0 {
		vuiParameters.VideoFormat = handleParseError(eg.ReadBits(3))
		vuiParameters.VideoFullRangeFlag = handleParseError(eg.ReadBit())
		vuiParameters.ColourDescriptionPresentFlag = handleParseError(eg.ReadBit())

		if vuiParameters.ColourDescriptionPresentFlag > 0 {
			vuiParameters.ColourPrimaries = handleParseError(eg.ReadBits(8))
			vuiParameters.TransferCharacteristics = handleParseError(eg.ReadBits(8))
			vuiParameters.MatrixCoefficients = handleParseError(eg.ReadBits(8))
		}
	}

	vuiParameters.ChromaLocInfoPresentFlag = handleParseError(eg.ReadBit())

	if vuiParameters.ChromaLocInfoPresentFlag > 0 {
		vuiParameters.ChromaSampleLocTypeTopField = handleParseError(eg.ReadUE())
		vuiParameters.ChromaSampleLocTypeBottomField = handleParseError(eg.ReadUE())
	}

	vuiParameters.TimingInfoPresentFlag = handleParseError(eg.ReadBit())

	if vuiParameters.TimingInfoPresentFlag > 0 {
		vuiParameters.NumUnitsInTick = handleParseError(eg.ReadBits(32))
		vuiParameters.TimeScale = handleParseError(eg.ReadBits(32))
		vuiParameters.FixedFrameRateFlag = handleParseError(eg.ReadBit())
	}

	vuiParameters.NalHrdParametersPresentFlag = handleParseError(eg.ReadBit())
	if vuiParameters.NalHrdParametersPresentFlag > 0 {
		hrd := NewHRDParameters()
		err = hrd.ParseExpGolombReader(eg)
		if err != nil {
			return
		}

		vuiParameters.HRDParameters = hrd
	}

	vuiParameters.VclHrdParametersPresentFlag = handleParseError(eg.ReadBit())
	if vuiParameters.NalHrdParametersPresentFlag > 0 || vuiParameters.VclHrdParametersPresentFlag > 0 {
		vuiParameters.LowDelayHrdFlag = handleParseError(eg.ReadBit())
	}

	vuiParameters.PicStructPresentFlag = handleParseError(eg.ReadBit())
	vuiParameters.BitstreamRestrictionFlag = handleParseError(eg.ReadBit())

	if vuiParameters.BitstreamRestrictionFlag > 0 {
		vuiParameters.MotionVectorsOverPicBoundariesFlag = handleParseError(eg.ReadBit())
		vuiParameters.MaxBytesPerPicDenom = handleParseError(eg.ReadUE())
		vuiParameters.MaxBitsPerMbDenom = handleParseError(eg.ReadUE())
		vuiParameters.Log2MaxMvLengthHorizontal = handleParseError(eg.ReadUE())
		vuiParameters.Log2MaxMvLengthVertical = handleParseError(eg.ReadUE())
		vuiParameters.MaxNumReorderFrames = handleParseError(eg.ReadUE())
		vuiParameters.MaxDecFrameBuffering = handleParseError(eg.ReadUE())
	}

	return
}
