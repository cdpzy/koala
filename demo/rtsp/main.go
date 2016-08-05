package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	//"encoding/binary"
)

const (
	PREFIX_SEI_NUT uint32 = 39
	SUFFIX_SEI_NUT        = 40
)

// import(
//     // "github.com/doublemo/koala"
//     // "log"
//     // "reflect"
//     //"fmt"
// )

//func main() {
// s := koala.NewRTSPServer()
// log.Println(s.GetRequest().String())
// log.Println(s.Serve(koala.SERVER_NET_PROTO_TCP, ":554"))
// rtsp := koala.NewRTSPServer()
// rtsp.HandlerFunc(Handler2)
// go rtsp.Serve(":554")

// http := koala.NewHTTPServer()
// http.HandlerFunc(Handler)
// log.Println(http.Serve(":9106"))
//}

// func Handler( req koala.Request, resp koala.Response ) {
//     method := koala.NewHandleMethod( req, resp )
//     v := reflect.ValueOf(method).MethodByName(req.GetMethod())
//     if !v.IsValid() {
//         resp.NotSupported( koala.AllowedMethod )
//         return
//     }

//     v.Call([]reflect.Value{})f
// }

// func Handler2( req koala.Request, resp koala.Response ) {
//     method := koala.NewHandleMethod( req, resp )
//     ch := req.GetInputChan()
//     for {
//         select {
//             case b,ok := <- ch:
//                  if !ok || b != nil {
//                      break
//                  }

//                  v := reflect.ValueOf(method).MethodByName(req.GetMethod())
//                  if !v.IsValid() {
//                      resp.NotSupported( koala.AllowedMethod )
//                     return
//                  }

//                 v.Call([]reflect.Value{})
//         }
//     }
// }

//http://www.2cto.com/kf/201506/408667.html
func main() {
	data := []byte{0x42, 0x00, 0x1E, 0xF1, 0x61, 0x62, 0x62}
	b := NewBitVector(data, 0, 50)
	fmt.Println(b.getBits(8), b.getBits(8), b.getBits(8), b.get_expGolomb())
	fid, err := os.Open("test.264")
	if err != nil {
		panic(err)
	}

	var (
		m     uint32
		issps bool
		sps   []byte
		pps   []byte
		ispps bool
	)
	for {
		to := make([]byte, 20000)
		_, err := fid.Read(to)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}

		//fmt.Printf("%v\n", to[0])
		reader := bytes.NewReader(to)

		for {
			h := make([]byte, 4)
			_, err = reader.Read(h)
			if err != nil {
				break
			}

			if toUint(h) != 0x00000001 && toUint(h)&0xFFFFFF00 != 0x00000100 {
				if issps {
					sps = append(sps, h...)
				} else if ispps {
					pps = append(pps, h...)
				}
				continue
			}

			if toUint(h)&0xFFFFFF00 == 0x00000100 {
				reader.UnreadByte()
			}

			issps = false
			ispps = false

			h = make([]byte, 4)
			_, err = reader.Read(h)
			if err != nil {
				break
			}

			m = toUint(h)
			m = m >> 24
			//nalRefIdc   := (m & 0x60) >> 5
			nalUnitType := m & 0x1F
			// fmt.Println("IDC/TYPE:", nalRefIdc, "/", nalUnitType)

			if is264SPS(nalUnitType) {
				issps = true
				sps = make([]byte, 0)
				sps = append(sps, h[1:]...)
			} else if is264PPS(nalUnitType) {
				ispps = true
				pps = append(pps, h[1:]...)
			}

		}

		break
	}

	bv := NewH264SPSPaser()
	fmt.Printf("SPS: %x\n", sps)
	//data = sps
	fmt.Println(bv.U(data, 8, 0), bv.U(data, 1, 8), bv.U(data, 1, 9), bv.U(data, 1, 10))
	fmt.Println(bv.U(data, 1, 11), bv.U(data, 1, 12), bv.U(data, 1, 13), bv.U(data, 2, 14), bv.U(data, 8, 16))
	fmt.Println("seq_parameter_set_id = ", bv.UE(data, 24))
	fmt.Println("log2_max_frame_num_minus4 = ", bv.UE(data, bv.GetStartBit()))
	fmt.Println("pic_order_cnt_type = ", bv.UE(data, bv.GetStartBit()))
	fmt.Println("log2_max_pic_order_cnt_lsb_minus4 = ", bv.UE(data, bv.GetStartBit()))
	fmt.Println("max_num_ref_frames = ", bv.UE(data, bv.GetStartBit()))
	s := bv.GetStartBit()
	fmt.Println("gaps_in_frame_num_value_allowed_flag  = ", bv.U(data, 1, s))
	fmt.Println("pic_width_in_mbs_minus1 = ", bv.UE(data, s+1))
	fmt.Println("pic_height_in_map_units_minus1  = ", bv.UE(data, bv.GetStartBit()))
	s = bv.GetStartBit()
	fmt.Println("frame_mbs_only_flag   = ", bv.U(data, 1, s))
	fmt.Printf("PPS: %x\n", pps)
}

func toUint(ptr []byte) uint32 {
	return (uint32(ptr[0]) << 24) | (uint32(ptr[1]) << 16) | (uint32(ptr[2]) << 8) | uint32(ptr[3])
}

func toUint16(ptr []byte) uint16 {
	return uint16(ptr[0]<<8 | ptr[1])
}

func toUint24(ptr []byte) uint32 {
	return uint32(ptr[0])<<16 | uint32(ptr[1])<<8 | uint32(ptr[2])
}

//二进制转十六进制
func btox(b string) string {
	base, _ := strconv.ParseInt(b, 2, 10)
	return strconv.FormatInt(base, 16)
}

//十六进制转二进制
func xtob(x string) string {
	base, _ := strconv.ParseInt(x, 16, 10)
	return strconv.FormatInt(base, 2)
}

func is264SPS(nal_unit_type uint32) bool {
	return nal_unit_type == 7
}

func is264PPS(nal_unit_type uint32) bool {
	return nal_unit_type == 8
}

func is264VCL(nal_unit_type uint32) bool {
	return nal_unit_type <= 5 && nal_unit_type > 0
}

func is264SEI(nal_unit_type uint32) bool {
	return nal_unit_type == 6
}

func is264EOF(nal_unit_type uint32) bool {
	return nal_unit_type == 10 || nal_unit_type == 11
}

// only 265
func is265VPS(nal_unit_type uint32) bool {
	return nal_unit_type == 32
}

func is265SPS(nal_unit_type uint32) bool {
	return nal_unit_type == 33
}

func is265PPS(nal_unit_type uint32) bool {
	return nal_unit_type == 34
}

func is265VCL(nal_unit_type uint32) bool {
	return nal_unit_type <= 31
}

func is265SEI(nal_unit_type uint32) bool {
	return nal_unit_type == PREFIX_SEI_NUT || nal_unit_type == SUFFIX_SEI_NUT
}

func is265EOF(nal_unit_type uint32) bool {
	return nal_unit_type == 36 || nal_unit_type == 37
}

func usuallyBeginsAccessUnit(code int, nal_unit_type uint32) bool {
	if code == 264 {
		return (nal_unit_type >= 6 && nal_unit_type <= 9) || (nal_unit_type >= 14 && nal_unit_type <= 18)
	}

	return (nal_unit_type >= 32 && nal_unit_type <= 35) ||
		(nal_unit_type == 39) || (nal_unit_type >= 41 && nal_unit_type <= 44) ||
		(nal_unit_type >= 48 && nal_unit_type <= 55)
}
