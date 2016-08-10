package h264

import (
	"errors"

	"github.com/doublemo/koala/helper"
)

// SequenceParameterSetScaling Scaling list semantics
type SequenceParameterSetScaling struct {
	List                        []int
	useDefaultScalingMatrixFlag bool
}

// SequenceParameterSetRBSP : Sequence parameter set RBSP semantics
// T-REC-H.264-201602-I!!PDF-E.pdf 7.4
type SequenceParameterSetRBSP struct {
	ProfileIdc                      uint
	ConstraintSet0Flag              uint
	ConstraintSet1Flag              uint
	ConstraintSet2Flag              uint
	ConstraintSet3Flag              uint
	ConstraintSet4Flag              uint
	ConstraintSet5Flag              uint
	ReservedZero2bit                uint
	LevelIdc                        uint
	SeqParameterSetID               uint
	ChromaFormatIdc                 uint
	SeparateColourPlaneFlag         uint
	BitDepthLumaMinus8              uint
	BitDepthChromaMinus8            uint
	QpprimeYZeroTransformBypassFlag uint
	SeqScalingMatrixPresentFlag     uint
	SeqScalingListPresentFlag       []uint
	ScalingList4X4                  []*SequenceParameterSetScaling
	ScalingList8X8                  []*SequenceParameterSetScaling
	Log2MaxFrameNumMinus4           uint
	PicOrderCntType                 uint
	Log2MaxPicOrderCntLsbMinus4     uint
	DeltaPicOrderAlwaysZeroFlag     uint
	OffsetForNonRefPic              int
	OffsetForTopToBottomField       int
	NumRefFramesInPicOrderCntCycle  uint
	OffsetForRefFrame               []int
	MaxNumRefFrames                 uint
	GapsInFrameNumValueAllowedFlag  uint
	PicWidthInMbsMinus1             uint
	PicHeightInMapUnitsMinus1       uint
	FrameMbsOnlyFlag                uint
	MbAdaptiveFrameFieldFlag        uint
	Direct8x8InferenceFlag          uint
	FrameCroppingFlag               uint
	FrameCropLeftOffset             uint
	FrameCropRightOffset            uint
	FrameCropTopOffset              uint
	FrameCropBottomOffset           uint
	VuiParametersPresentFlag        uint
	VuiParameters                   *VUIParameters
	Extension                       *SequenceParameterSetExtensionRBSP //Sequence parameter set extension RBSP semantics
}

func NewSequenceParameterSetRBSP() *SequenceParameterSetRBSP {
	return new(SequenceParameterSetRBSP)
}

// ParseBytes Sequence parameter set RBSP syntax
// T-REC-H.264-201602-I!!PDF-E.pdf 7.3
func (sequenceParameterSetRBSP *SequenceParameterSetRBSP) ParseBytes(b []byte) (err error) {
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

	eg := helper.NewExpGolombReader(b)
	sequenceParameterSetRBSP.ProfileIdc = handleParseError(eg.ReadBits(8))
	sequenceParameterSetRBSP.ConstraintSet0Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ConstraintSet1Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ConstraintSet2Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ConstraintSet3Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ConstraintSet4Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ConstraintSet5Flag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.ReservedZero2bit = handleParseError(eg.ReadBits(2))
	sequenceParameterSetRBSP.LevelIdc = handleParseError(eg.ReadBits(8))
	sequenceParameterSetRBSP.SeqParameterSetID = handleParseError(eg.ReadUE())

	//
	if helper.InUint(sequenceParameterSetRBSP.ProfileIdc, []uint{100, 110, 122, 244, 44, 83, 86, 118, 128, 138, 139, 134, 135}) {
		sequenceParameterSetRBSP.ChromaFormatIdc = handleParseError(eg.ReadUE())
		if sequenceParameterSetRBSP.ChromaFormatIdc == 3 {
			sequenceParameterSetRBSP.SeparateColourPlaneFlag = handleParseError(eg.ReadBit())
		}

		sequenceParameterSetRBSP.BitDepthLumaMinus8 = handleParseError(eg.ReadUE())
		sequenceParameterSetRBSP.BitDepthChromaMinus8 = handleParseError(eg.ReadUE())
		sequenceParameterSetRBSP.QpprimeYZeroTransformBypassFlag = handleParseError(eg.ReadBit())
		sequenceParameterSetRBSP.SeqScalingMatrixPresentFlag = handleParseError(eg.ReadBit())

		if sequenceParameterSetRBSP.SeqScalingMatrixPresentFlag == 1 {
			x := 8
			if sequenceParameterSetRBSP.ChromaFormatIdc == 3 {
				x = 12
			}

			for i := 0; i < x; i++ {
				m := handleParseError(eg.ReadBit())
				sequenceParameterSetRBSP.SeqScalingListPresentFlag = append(sequenceParameterSetRBSP.SeqScalingListPresentFlag, m)
				if m > 0 {
					if i < 6 {
						scalingList, useDefaultScalingMatrixFlag := scalingList(eg, 16)
						sequenceParameterSetRBSP.ScalingList4X4 = append(sequenceParameterSetRBSP.ScalingList4X4, &SequenceParameterSetScaling{scalingList, useDefaultScalingMatrixFlag})
					} else {
						scalingList, useDefaultScalingMatrixFlag := scalingList(eg, 16)
						sequenceParameterSetRBSP.ScalingList8X8 = append(sequenceParameterSetRBSP.ScalingList8X8, &SequenceParameterSetScaling{scalingList, useDefaultScalingMatrixFlag})
					}
				}
			}
		}
	}

	sequenceParameterSetRBSP.Log2MaxFrameNumMinus4 = handleParseError(eg.ReadUE())
	sequenceParameterSetRBSP.PicOrderCntType = handleParseError(eg.ReadUE())

	if sequenceParameterSetRBSP.PicOrderCntType == 0 {
		sequenceParameterSetRBSP.Log2MaxPicOrderCntLsbMinus4 = handleParseError(eg.ReadUE())
	} else {
		sequenceParameterSetRBSP.DeltaPicOrderAlwaysZeroFlag = handleParseError(eg.ReadBit())
		sequenceParameterSetRBSP.OffsetForNonRefPic = handleParseSEError(eg.ReadSE())
		sequenceParameterSetRBSP.OffsetForTopToBottomField = handleParseSEError(eg.ReadSE())
		sequenceParameterSetRBSP.NumRefFramesInPicOrderCntCycle = handleParseError(eg.ReadUE())

		for i := 0; i < int(sequenceParameterSetRBSP.NumRefFramesInPicOrderCntCycle); i++ {
			sequenceParameterSetRBSP.OffsetForRefFrame = append(sequenceParameterSetRBSP.OffsetForRefFrame, handleParseSEError(eg.ReadSE()))
		}
	}

	sequenceParameterSetRBSP.MaxNumRefFrames = handleParseError(eg.ReadUE())
	sequenceParameterSetRBSP.GapsInFrameNumValueAllowedFlag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.PicWidthInMbsMinus1 = handleParseError(eg.ReadUE())
	sequenceParameterSetRBSP.PicHeightInMapUnitsMinus1 = handleParseError(eg.ReadUE())
	sequenceParameterSetRBSP.FrameMbsOnlyFlag = handleParseError(eg.ReadBit())

	if sequenceParameterSetRBSP.FrameMbsOnlyFlag < 1 {
		sequenceParameterSetRBSP.MbAdaptiveFrameFieldFlag = handleParseError(eg.ReadBit())
	}

	sequenceParameterSetRBSP.Direct8x8InferenceFlag = handleParseError(eg.ReadBit())
	sequenceParameterSetRBSP.FrameCroppingFlag = handleParseError(eg.ReadBit())

	if sequenceParameterSetRBSP.FrameCroppingFlag > 0 {
		sequenceParameterSetRBSP.FrameCropLeftOffset = handleParseError(eg.ReadUE())
		sequenceParameterSetRBSP.FrameCropRightOffset = handleParseError(eg.ReadUE())
		sequenceParameterSetRBSP.FrameCropTopOffset = handleParseError(eg.ReadUE())
		sequenceParameterSetRBSP.FrameCropBottomOffset = handleParseError(eg.ReadUE())
	}

	sequenceParameterSetRBSP.VuiParametersPresentFlag = handleParseError(eg.ReadBit())
	if sequenceParameterSetRBSP.VuiParametersPresentFlag > 0 {
		vui := NewVUIParameters()
		err = vui.ParseExpGolombReader(eg)
		if err != nil {
			return
		}

		sequenceParameterSetRBSP.VuiParameters = vui
	}

	err = nil
	return
}

// scalingList Scaling list syntax
func scalingList(eg *helper.ExpGolombReader, sizeOfScalingList int) ([]int, bool) {
	var (
		lastScale                   int = 8
		nextScale                   int = 8
		useDefaultScalingMatrixFlag bool
	)

	scalingList := make([]int, 0)

	for j := 0; j < sizeOfScalingList; j++ {
		if nextScale != 0 {
			delta_scale := handleParseSEError(eg.ReadSE())
			nextScale = (lastScale + delta_scale + 256) % 256
			if j == 0 && nextScale == 0 {
				useDefaultScalingMatrixFlag = true
			} else {
				useDefaultScalingMatrixFlag = false
			}
		}

		x := nextScale
		if nextScale == 0 {
			x = lastScale
		}
		scalingList = append(scalingList, x)
	}

	return scalingList, useDefaultScalingMatrixFlag
}

func handleParseError(n uint, e error) uint {
	if e != nil {
		panic(e)
	}
	return n
}

func handleParseSEError(n int, e error) int {
	if e != nil {
		panic(e)
	}
	return n
}
