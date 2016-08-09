package main

import "math"

type H264SPSPaser struct {
	startBit uint
}

func NewH264SPSPaser() *H264SPSPaser {
	return new(H264SPSPaser)
}

// http://guoh.org/lifelog/2013/10/h-264-bit-stream-sps-pps-idr-nalu/

/*
 * 从数据流data中第StartBit位开始读，读bitCnt位，以无符号整形返回
 */
func (h264SPSPaser *H264SPSPaser) U(data []byte, bitCnt int, startBit uint) uint {
	var ret uint
	start := startBit

	for i := 0; i < bitCnt; i++ {
		ret <<= 1
		if (data[start/8] & (0x80 >> (start % 8))) != 0 {
			ret += 1
		}

		start++
	}

	//h264SPSPaser.startBit += start

	return ret
}

/*
 * 无符号指数哥伦布编码
 * leadingZeroBits = ?1;
 * for( b = 0; !b; leadingZeroBits++ )
 *    b = read_bits( 1 )
 * 变量codeNum 按照如下方式赋值：
 * codeNum = 2^leadingZeroBits ? 1 + read_bits( leadingZeroBits )
 * 这里read_bits( leadingZeroBits )的返回值使用高位在先的二进制无符号整数表示。
 */
func (h264SPSPaser *H264SPSPaser) UE(data []byte, startBit uint) uint {
	var (
		ret          uint = 0
		tempStartBit uint = 0
		b            uint
	)

	leadingZeroBits := -1
	tempStartBit = startBit

	for b = 0; b != 1; leadingZeroBits++ {
		b = h264SPSPaser.U(data, 1, tempStartBit)
		tempStartBit++
	}

	// fmt.Printf("H264SPSPaser ue leadingZeroBits = %d \n", leadingZeroBits)
	// fmt.Printf("Math.pow(2, leadingZeroBits) = %f \n", math.Pow(2, float64(leadingZeroBits)))
	// fmt.Printf("tempStartBit = %d \n", tempStartBit)
	ret = uint(math.Pow(2, float64(leadingZeroBits))) - 1 + h264SPSPaser.U(data, leadingZeroBits, tempStartBit)
	h264SPSPaser.startBit = tempStartBit + uint(leadingZeroBits)

	//fmt.Printf("ue startBit = %d \n", h264SPSPaser.startBit)
	return ret
}

/*
 * 有符号指数哥伦布编码
 * 9.1.1 有符号指数哥伦布编码的映射过程
 *按照9.1节规定，本过程的输入是codeNum。
 *本过程的输出是se(v)的值。
 *表9-3中给出了分配给codeNum的语法元素值的规则，语法元素值按照绝对值的升序排列，负值按照其绝对
 *值参与排列，但列在绝对值相等的正值之后。
 *表 9-3－有符号指数哥伦布编码语法元素se(v)值与codeNum的对应
 *codeNum 语法元素值
 *	0 		0
 *	1		1
 *	2		?1
 *	3		2
 *	4		?2
 *	5		3
 *	6		?3
 *	k		(?1)^(k+1) Ceil( k÷2 )
 */

func (h264SPSPaser *H264SPSPaser) SE(data []byte, startBit uint) int {

	codeNum := h264SPSPaser.UE(data, startBit)
	ret := (math.Pow(-1, float64(codeNum)+1) * math.Ceil(float64(codeNum)/2))
	return int(ret)
}

func (h264SPSPaser *H264SPSPaser) GetStartBit() uint {
	//fmt.Println("GetStartBit = ", h264SPSPaser.startBit)
	return h264SPSPaser.startBit
}
