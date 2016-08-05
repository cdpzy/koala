package main

const (
	NALU_TYPE_SLICE uint = iota + 1
	NALU_TYPE_DPA
	NALU_TYPE_DPB
	NALU_TYPE_DPC
	NALU_TYPE_IDR
	NALU_TYPE_SEI
	NALU_TYPE_SPS
	NALU_TYPE_PPS
	NALU_TYPE_AUD
	NALU_TYPE_EOSEQ
	NALU_TYPE_EOSTREAM
	NALU_TYPE_FILL
)

type SPS struct {
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
	Log2MaxFrameNumMinus4           uint
	PicOrderCntType                 uint
	Log2MaxPicOrderCntLsbMinus4     uint
	DeltaPicOrderAlwaysZeroFlag     uint
	OffsetForNonRefPic              uint
	OffsetForTopToBottomField       uint
	NumRefFramesInPicOrderCntCycle  uint
	OffsetForRefFrame               []uint
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
}

func DecodeSPS(data []byte) *SPS {
	b := NewBitVector(data, 0, 1000)
	sps := &SPS{
		ProfileIdc:         b.getBits(8),
		ConstraintSet0Flag: b.get1Bit(),
		ConstraintSet1Flag: b.get1Bit(),
		ConstraintSet2Flag: b.get1Bit(),
		ConstraintSet3Flag: b.get1Bit(),
		ConstraintSet4Flag: b.get1Bit(),
		ConstraintSet5Flag: b.get1Bit(),
		ReservedZero2bit:   b.getBits(2),
		LevelIdc:           b.getBits(8),
		SeqParameterSetID:  b.get_expGolomb(),
	}

	if sps.ProfileIdc == 100 ||
		sps.ProfileIdc == 110 ||
		sps.ProfileIdc == 122 ||
		sps.ProfileIdc == 244 ||
		sps.ProfileIdc == 44 ||
		sps.ProfileIdc == 83 ||
		sps.ProfileIdc == 86 ||
		sps.ProfileIdc == 118 ||
		sps.ProfileIdc == 128 ||
		sps.ProfileIdc == 138 ||
		sps.ProfileIdc == 139 ||
		sps.ProfileIdc == 134 ||
		sps.ProfileIdc == 135 {
		sps.ChromaFormatIdc = b.get_expGolomb()
	}

	if sps.ChromaFormatIdc == 3 {
		sps.SeparateColourPlaneFlag = b.get1Bit()
	}
	return sps
}

////b := NewBitVector(data, 0, 50)
//fmt.Println(b.getBits(8), b.getBits(8), b.getBits(8), b.get_expGolomb())

type SequenceParameterSetRBSP struct{}

type PictureParameterSetRBSP struct {
	PicParameterSetID                     uint
	SeqParameterSetID                     uint
	EntropyCodingModeFlag                 uint
	BottomFieldPicOrderInFramePresentFlag uint
	NumSliceGroupsMinus1                  uint
}
