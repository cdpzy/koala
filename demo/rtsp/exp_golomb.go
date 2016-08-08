/**
 * 指数哥伦布编码
 * 规定语法元素的编解码模式的描述符如下：
 * 比特串：
 * b(8):    任意形式的8比特字节（就是为了说明语法元素是为8个比特，没有语法上的含义）
 * f(n):     n位固定模式比特串（其值固定，如forbidden_zero_bit的值恒为0）
 * i(n):     使用n比特的有符号整数（语法中没有采用此格式）
 * u(n):    n位无符号整数
 *
 * 指数哥伦布编码：
 * ue(v):   无符号整数指数哥伦布码编码的语法元素
 * se(v):   有符号整数指数哥伦布编码的语法元素，位在先
 * te(v):    舍位指数哥伦布码编码语法元素，左位在先
 * ce(v)：CAVLC
 * ae(v)：CABAC。
 * https://jordicenzano.name/2014/08/31/the-source-code-of-a-minimal-h264-encoder-c/
 */
//  m := math.Floor(math.Log2(float64(v) + 1))
//  i := float64(v) + 1 - math.Pow(2, m)

//  fmt.Println(m, "->", i)
//  fmt.Printf("%b\n", uint(i))
package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

type ExpGolomb struct {
	data     []byte
	buffer   *bytes.Buffer
	startBit uint
}

/*
 * 从数据流data中第StartBit位开始读，读bitCnt位，以无符号整形返回
 */
func (expGolomb *ExpGolomb) ReadBits(numBits int) uint {
	var (
		ret uint
	)

	start := expGolomb.startBit
	for i := 0; i < numBits; i++ {
		ret <<= 1
		if (expGolomb.data[start/8] & (0x80 >> (start % 8))) != 0 {
			ret += 1
		}

		start++
	}

	expGolomb.startBit += uint(numBits)
	return ret
}

func (expGolomb *ExpGolomb) ReadAtBits(numBits int, startBit uint) uint {
	var (
		ret uint
	)

	start := startBit
	for i := 0; i < numBits; i++ {
		ret <<= 1
		if (expGolomb.data[start/8] & (0x80 >> (start % 8))) != 0 {
			ret += 1
		}

		start++
	}

	return ret
}

func (expGolomb *ExpGolomb) ReadBit() uint {
	return expGolomb.ReadBits(1)
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
func (expGolomb *ExpGolomb) ReadUE() uint {
	var (
		idx uint
		b   uint
	)

	leadingZeroBits := -1
	idx = expGolomb.startBit
	for b = 0; b != 1; leadingZeroBits++ {
		b = expGolomb.ReadBits(1)
		idx++
	}

	ret := uint(math.Pow(2, float64(leadingZeroBits))) - 1 + expGolomb.ReadAtBits(leadingZeroBits, idx)
	expGolomb.startBit = idx + uint(leadingZeroBits)
	return ret
}

func (expGolomb *ExpGolomb) ReadAtUE(startBit uint) (uint, uint) {
	var (
		idx uint
		b   uint
	)

	leadingZeroBits := -1
	idx = startBit
	for b = 0; b != 1; leadingZeroBits++ {
		b = expGolomb.ReadBits(1)
		idx++
	}

	ret := uint(math.Pow(2, float64(leadingZeroBits))) - 1 + expGolomb.ReadAtBits(leadingZeroBits, idx)
	startBit = idx + uint(leadingZeroBits)
	return startBit, ret
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
func (expGolomb *ExpGolomb) ReadSE() int {
	codeNum := expGolomb.ReadUE()
	ret := (math.Pow(-1, float64(codeNum)+1) * math.Ceil(float64(codeNum)/2))
	return int(ret)
}

func (expGolomb *ExpGolomb) ReadAtSE(startBit uint) (uint, int) {
	s, codeNum := expGolomb.ReadAtUE(startBit)
	ret := (math.Pow(-1, float64(codeNum)+1) * math.Ceil(float64(codeNum)/2))
	return s, int(ret)
}

func (expGolomb *ExpGolomb) GetStartBit() uint {
	return expGolomb.startBit
}

func (expGolomb *ExpGolomb) WriteBits(v uint, numBits int) error {
	if numBits <= 0 || numBits > 64 {
		return errors.New("numbits must be between 1 ... 64")
	}

	nBit := 0
	n := numBits - 1
	for {
		nBit = expGolomb.getBitNum(v, n)
		n--
		if n < 0 {
			break
		}
		expGolomb.addbittostream(nBit)
	}
	return nil
}

func (expGolomb *ExpGolomb) addbittostream(nVal int) error {
	var (
		nValTmp uint
		v       uint
	)

	v = 77
	b := []uint{v >> 24, v >> 16, v >> 8, v}
	nBytePos := expGolomb.buffer.Len() % len(b)
	nBitPosInByte := 7 - expGolomb.buffer.Len()%8
	nValTmp = b[nBytePos]
	//fmt.Println("nBytePos:", nBytePos, nValTmp)
	if nVal > 0 {
		nValTmp = (nValTmp | uint(math.Pow(2, float64(nBitPosInByte))))
	} else {
		nValTmp = (nValTmp | ^uint(math.Pow(2, float64(nBitPosInByte))))
	}

	expGolomb.buffer.WriteByte(byte(nValTmp))

	fmt.Println(expGolomb.buffer.Bytes())
	return nil
}

func (expGolomb *ExpGolomb) getBitNum(v uint, numBits int) int {
	mask := uint(math.Pow(2, float64(numBits)))
	if (v & mask) > 0 {
		return 1
	}

	return 0
}

func (expGolomb *ExpGolomb) SkipBits(numBits uint) {
	expGolomb.startBit += numBits
}

func NewExpGolomb(data []byte) *ExpGolomb {
	expGolomb := new(ExpGolomb)
	expGolomb.data = data
	expGolomb.buffer = bytes.NewBuffer(make([]byte, 0))
	expGolomb.startBit = 0

	return expGolomb
}
