package main

import (
	"fmt"
	"testing"
)

var _ = &testing.T{}

func expGolombReadBits() {
	data := []byte{0x42, 0x00, 0x1E, 0xF1, 0x61, 0x62, 0x62}
	eg := NewExpGolombReader(data)
	fmt.Println(eg.ReadBits(8), eg.ReadBits(1), eg.ReadBits(1), eg.ReadBits(1))
	fmt.Println(eg.ReadBits(1), eg.ReadBits(1), eg.ReadBits(1), eg.ReadBits(2), eg.ReadBits(8))
	fmt.Println("seq_parameter_set_id = ", eg.ReadUE())
	fmt.Println("log2_max_frame_num_minus4 = ", eg.ReadUE())
	fmt.Println("pic_order_cnt_type = ", eg.ReadUE())
	fmt.Println("log2_max_pic_order_cnt_lsb_minus4 = ", eg.ReadUE())
	fmt.Println("max_num_ref_frames = ", eg.ReadUE())
	fmt.Println("gaps_in_frame_num_value_allowed_flag  = ", eg.ReadBits(1))
	fmt.Println("pic_width_in_mbs_minus1 = ", eg.ReadUE())
	fmt.Println("pic_height_in_map_units_minus1  = ", eg.ReadUE())
	fmt.Println("frame_mbs_only_flag   = ", eg.ReadBits(1))
}

func expGolombWriteBits() {
	wg := NewExpGolombWriter()
	wg.WriteBits(66, 8)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 2)
	wg.WriteBits(30, 8)
	wg.WriteUE(0)
	wg.WriteUE(0)
	wg.WriteUE(0)
	wg.WriteUE(0)
	wg.WriteUE(10)
	wg.WriteBits(0, 1)
	wg.WriteUE(10)
	wg.WriteUE(8)
	wg.WriteBits(1, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(0, 1)
	wg.WriteBits(1, 1)
	wg.WriteBits(0, 1)
	fmt.Println(wg.Bytes())
}

func ExampleExpGolombReadBits() {
	expGolombReadBits()
	// Output:
	// 66 0 0 0
	// 0 0 0 0 30
	// seq_parameter_set_id =  0
	// log2_max_frame_num_minus4 =  0
	// pic_order_cnt_type =  0
	// log2_max_pic_order_cnt_lsb_minus4 =  0
	// max_num_ref_frames =  10
	// gaps_in_frame_num_value_allowed_flag  =  0
	// pic_width_in_mbs_minus1 =  10
	// pic_height_in_map_units_minus1  =  8
	// frame_mbs_only_flag   =  1
	//
}

func ExampleExpGolombWriteBits() {
	expGolombWriteBits()
	// Output:
	// [66 0 30 241 97 98 98]
	//
}

func BenchmarkReadUE(b *testing.B) {
	b.StopTimer()
	data := []byte{0x42, 0x00, 0x1E, 0xF1, 0x61, 0x62, 0x62}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		eg := NewExpGolombReader(data)
		eg.ReadAtUE(24)
	}
}

func BenchmarkReadBits(b *testing.B) {
	b.StopTimer()
	data := []byte{0x42, 0x00, 0x1E, 0xF1, 0x61, 0x62, 0x62}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		eg := NewExpGolombReader(data)
		eg.ReadBits(8)
	}
}

func BenchmarkWriterUE(b *testing.B) {
	b.StopTimer()

	wg := NewExpGolombWriter()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.WriteUE(10)
	}
}

func BenchmarkWriterSE(b *testing.B) {
	b.StopTimer()

	wg := NewExpGolombWriter()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.WriteSE(10)
	}
}

func BenchmarkWriterBits(b *testing.B) {
	b.StopTimer()

	wg := NewExpGolombWriter()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.WriteBits(77, 8)
	}
}
