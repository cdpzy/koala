// @author Randy Ma <435420057@qq.com>

/**
 *
 * 用于处理H264 PPS SPS 解码
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
 * te(v):   舍位指数哥伦布码编码语法元素，左位在先
 * ce(v)：CAVLC
 * ae(v)：CABAC。
 * 参考 : https://jordicenzano.name/2014/08/31/the-source-code-of-a-minimal-h264-encoder-c/
 *        http://www.2cto.com/kf/201506/408667.html
 *        http://guoh.org/lifelog/2013/10/h-264-bit-stream-sps-pps-idr-nalu/
 */

package helper

import (
	"bytes"
	"errors"
	"math"
)

const (
	BUFFER_SIZE_BITS                    = 24
	BUFFER_SIZE_BYTES                   = (24 / 8)
	H264_EMULATION_PREVENTION_BYTE uint = 0x03
)

// ExpGolombReader 指数哥伦布编码流进行处理,主要针对H264 PPS,SPS等
type ExpGolombReader struct {
	buffer   []byte
	startBit uint
}

// ExpGolombWriter 指数哥伦布编码写入
type ExpGolombWriter struct {
	bytes  []byte
	buffer *bytes.Buffer

	m_nLastbitinbuffer int // bit 计数器

	m_nStartingbyte int // 起始位置记录
}

// NewExpGolombReader 创建指数哥伦布编码流解析
func NewExpGolombReader(d []byte) *ExpGolombReader {
	return &ExpGolombReader{d, 0}
}

// NewExpGolombWriter 创建指数哥伦布编码流编码
func NewExpGolombWriter() *ExpGolombWriter {
	return &ExpGolombWriter{
		make([]byte, BUFFER_SIZE_BITS),
		bytes.NewBuffer(make([]byte, 0)),
		0,
		0,
	}
}

// ReadBits 从数据流buffer中第StartBit位开始读，读numBits位，以无符号整形返回
func (expGolombReader *ExpGolombReader) ReadBits(numBits int) (uint, error) {
	var ret uint
	start := expGolombReader.startBit
	for i := 0; i < numBits; i++ {
		ret <<= 1
		if (start/8 + 1) > uint(len(expGolombReader.buffer)) {
			return 0, errors.New("EOF")
		}

		if (expGolombReader.buffer[start/8] & (0x80 >> (start % 8))) != 0 {
			ret++
		}

		start++
	}

	expGolombReader.startBit += uint(numBits)
	return ret, nil
}

// ReadAtBits 根据指定起始位置从数据流中读numBits位，以无符号整形返回
// numBits
// startBit 起始位置
func (expGolombReader *ExpGolombReader) ReadAtBits(numBits int, startBit uint) (uint, error) {
	var ret uint

	start := startBit
	for i := 0; i < numBits; i++ {
		ret <<= 1

		if (start/8 + 1) > uint(len(expGolombReader.buffer)) {
			return 0, errors.New("EOF")
		}

		if (expGolombReader.buffer[start/8] & (0x80 >> (start % 8))) != 0 {
			ret++
		}

		start++
	}

	return ret, nil
}

// ReadBit 读取1bit，以无符号整形返回
func (expGolombReader *ExpGolombReader) ReadBit() (uint, error) {
	return expGolombReader.ReadBits(1)
}

// ReadIBits 从数据流buffer中第StartBit位开始读，读numBits位，以有符号整形返回
func (expGolombReader *ExpGolombReader) ReadIBits(numBits int) (int, error) {
	var ret int
	start := expGolombReader.startBit
	for i := 0; i < numBits; i++ {
		ret <<= 1
		if (start/8 + 1) > uint(len(expGolombReader.buffer)) {
			return 0, errors.New("EOF")
		}

		if (expGolombReader.buffer[start/8] & (0x80 >> (start % 8))) != 0 {
			ret++
		}

		start++
	}

	expGolombReader.startBit += uint(numBits)
	return ret, nil
}

// ReadUE 无符号指数哥伦布编码
/*
 * leadingZeroBits = ?1;
 * for( b = 0; !b; leadingZeroBits++ )
 *    b = ReadBits( 1 )
 * 变量codeNum 按照如下方式赋值：
 * codeNum = 2^leadingZeroBits ? 1 + ReadBits( leadingZeroBits )
 * 这里ReadBits( leadingZeroBits )的返回值使用高位在先的二进制无符号整数表示。
 */
func (expGolombReader *ExpGolombReader) ReadUE() (uint, error) {
	var (
		idx uint
		b   uint
		err error
	)

	leadingZeroBits := -1
	idx = expGolombReader.startBit
	for b = 0; b != 1; leadingZeroBits++ {
		b, err = expGolombReader.ReadBit()
		if err != nil {
			return 0, err
		}
		idx++
	}

	n, err := expGolombReader.ReadAtBits(leadingZeroBits, idx)
	if err != nil {
		return 0, err
	}
	ret := uint(math.Pow(2, float64(leadingZeroBits))) - 1 + n
	expGolombReader.startBit = idx + uint(leadingZeroBits)
	return ret, nil
}

// ReadAtUE 无符号指数哥伦布编码, 指定起始位置
func (expGolombReader *ExpGolombReader) ReadAtUE(startBit uint) (uint, uint, error) {
	var (
		idx uint
		b   uint
		err error
	)

	leadingZeroBits := -1
	idx = startBit
	for b = 0; b != 1; leadingZeroBits++ {
		b, err = expGolombReader.ReadAtBits(1, idx)
		if err != nil {
			return 0, 0, err
		}
		idx++
	}

	n, err := expGolombReader.ReadAtBits(leadingZeroBits, idx)
	if err != nil {
		return 0, 0, err
	}

	ret := uint(math.Pow(2, float64(leadingZeroBits))) - 1 + n
	startBit = idx + uint(leadingZeroBits)
	return startBit, ret, nil
}

// ReadSE 有符号指数哥伦布编码
/*
 * T-REC-H.264-201602-I!!PDF-E.pdf
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
func (expGolombReader *ExpGolombReader) ReadSE() (int, error) {
	codeNum, err := expGolombReader.ReadUE()
	if err != nil {
		return 0, err
	}
	ret := (math.Pow(-1, float64(codeNum)+1) * math.Ceil(float64(codeNum)/2))
	return int(ret), nil
}

// ReadAtSE 有符号指数哥伦布编码,指定起始位置
func (expGolombReader *ExpGolombReader) ReadAtSE(startBit uint) (uint, int, error) {
	s, codeNum, err := expGolombReader.ReadAtUE(startBit)
	if err != nil {
		return 0, 0, err
	}
	ret := (math.Pow(-1, float64(codeNum)+1) * math.Ceil(float64(codeNum)/2))
	return s, int(ret), nil
}

// SkipBits 跳过
func (expGolombReader *ExpGolombReader) SkipBits(numBits uint) {
	expGolombReader.startBit += numBits
}

// GetStartBit 获取当前起始位置
func (expGolombReader *ExpGolombReader) GetStartBit() uint {
	return expGolombReader.startBit
}

// WriteByte
func (expGolombWriter *ExpGolombWriter) WriteByte(v uint) error {
	if (expGolombWriter.m_nLastbitinbuffer % 8) == 0 {
		return expGolombWriter.addbytetostream(int(v))
	}

	return expGolombWriter.WriteBits(v, 8)
}

// WriteBits 写入指定位的值
func (expGolombWriter *ExpGolombWriter) WriteBits(v uint, numBits int) error {
	if numBits <= 0 || numBits > 64 {
		return errors.New("numbits must be between 1 ... 64")
	}

	nBit := 0
	n := numBits - 1
	for n >= 0 {
		nBit = expGolombWriter.getBitNum(v, n)
		expGolombWriter.addbittostream(nBit)
		n--
	}
	return nil
}

// 无符号指数哥伦布编码
func (expGolombWriter *ExpGolombWriter) WriteUE(v uint) error {
	lvalint := v + 1
	nnumbits := int(math.Log2(float64(lvalint)) + 1)

	for n := 0; n < (nnumbits - 1); n++ {
		expGolombWriter.WriteBits(0, 1)
	}

	expGolombWriter.WriteBits(lvalint, nnumbits)
	return nil
}

// 有符号指数哥伦布编码
func (expGolombWriter *ExpGolombWriter) WriteSE(v int) error {
	lvalint := uint(math.Abs(float64(v))*2 - 1)
	if v <= 0 {
		lvalint = uint(2 * math.Abs(float64(v)))
	}

	expGolombWriter.WriteUE(lvalint)
	return nil
}

// Write4bytesNoEmulationPrevention 写入4个byte
func (expGolombWriter *ExpGolombWriter) Write4bytesNoEmulationPrevention(nVal uint, bDoAlign bool) error {
	if bDoAlign {
		expGolombWriter.dobytealign()
	}

	if (expGolombWriter.m_nLastbitinbuffer % 8) != 0 {
		return errors.New("Error: Save to file must be byte aligned")
	}

	for expGolombWriter.m_nLastbitinbuffer != 0 {
		expGolombWriter.savebufferbyte(true)
	}

	expGolombWriter.buffer.WriteByte(byte((nVal & 0xFF000000) >> 24))
	expGolombWriter.buffer.WriteByte(byte((nVal & 0xFF000000) >> 16))
	expGolombWriter.buffer.WriteByte(byte((nVal & 0xFF000000) >> 8))
	expGolombWriter.buffer.WriteByte(byte(nVal & 0xFF000000))
	return nil
}

// getBitNum 获取bit位mask 计算参考
func (expGolombWriter *ExpGolombWriter) getBitNum(v uint, numBits int) int {
	mask := uint(math.Pow(2, float64(numBits)))
	if (v & mask) > 0 {
		return 1
	}

	return 0
}

// addbittostream 写入bit 流
func (expGolombWriter *ExpGolombWriter) addbittostream(nVal int) error {
	var nValTmp uint

	if expGolombWriter.m_nLastbitinbuffer >= BUFFER_SIZE_BITS {
		expGolombWriter.savebufferbyte(true)
	}

	nBytePos := (expGolombWriter.m_nStartingbyte + (expGolombWriter.m_nLastbitinbuffer / 8)) % BUFFER_SIZE_BYTES
	nBitPosInByte := 7 - expGolombWriter.m_nLastbitinbuffer%8
	nValTmp = uint(expGolombWriter.bytes[nBytePos])

	if nVal > 0 {
		nValTmp = (nValTmp | uint(math.Pow(2, float64(nBitPosInByte))))
	} else {
		nValTmp = (nValTmp & ^uint(math.Pow(2, float64(nBitPosInByte))))
	}

	expGolombWriter.bytes[nBytePos] = byte(nValTmp)
	expGolombWriter.m_nLastbitinbuffer++
	return nil
}

// addbittostream 写入byte 到bit 流
func (expGolombWriter *ExpGolombWriter) addbytetostream(nVal int) error {
	if expGolombWriter.m_nLastbitinbuffer >= BUFFER_SIZE_BITS {
		expGolombWriter.savebufferbyte(true)
	}

	nBytePos := (expGolombWriter.m_nStartingbyte + (expGolombWriter.m_nLastbitinbuffer / 8)) % BUFFER_SIZE_BYTES
	nBitPosInByte := 7 - expGolombWriter.m_nLastbitinbuffer%8
	if nBitPosInByte != 7 {
		return errors.New("Error: inserting not aligment byte")
	}

	expGolombWriter.bytes[nBytePos] = byte(uint(nVal))
	expGolombWriter.m_nLastbitinbuffer = expGolombWriter.m_nLastbitinbuffer + 8
	return nil
}

// savebufferbyte 写入byte 到buffer
func (expGolombWriter *ExpGolombWriter) savebufferbyte(bemulationprevention bool) error {
	var (
		bemulationpreventionexecuted bool
	)

	if (expGolombWriter.m_nLastbitinbuffer % 8) != 0 {
		return errors.New("Error: Save to file must be byte aligned")
	}

	if (expGolombWriter.m_nLastbitinbuffer / 8) <= 0 {
		return errors.New("Error: NO bytes to save")
	}

	if bemulationprevention {
		// Emulation prevention will be used:
		/*As per h.264 spec,
		rbsp_data shouldn't contain
				- 0x 00 00 00
				- 0x 00 00 01
				- 0x 00 00 02
				- 0x 00 00 03

		rbsp_data shall be in the following way
				- 0x 00 00 03 00
				- 0x 00 00 03 01
				- 0x 00 00 03 02
				- 0x 00 00 03 03
		*/
		// Check if emulation prevention is needed (emulation prevention is byte align defined)
		byte1 := expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + 0) % BUFFER_SIZE_BYTES)]
		byte2 := expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + 1) % BUFFER_SIZE_BYTES)]
		byte3 := expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + 2) % BUFFER_SIZE_BYTES)]

		if byte1 == 0x00 && byte2 == 0x00 && (byte3 == 0x00 || byte3 == 0x01 || byte3 == 0x02 || byte3 == 0x03) {
			nbuffersaved := 0
			expGolombWriter.buffer.WriteByte(expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + nbuffersaved) % BUFFER_SIZE_BYTES)])

			nbuffersaved++
			expGolombWriter.buffer.WriteByte(expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + nbuffersaved) % BUFFER_SIZE_BYTES)])

			nbuffersaved++
			expGolombWriter.buffer.WriteByte(byte(H264_EMULATION_PREVENTION_BYTE))

			for nbuffersaved < BUFFER_SIZE_BYTES {
				expGolombWriter.buffer.WriteByte(expGolombWriter.bytes[((expGolombWriter.m_nStartingbyte + nbuffersaved) % BUFFER_SIZE_BYTES)])
				nbuffersaved++
			}

			expGolombWriter.resetBytes()
			bemulationpreventionexecuted = true
		}

	}

	if !bemulationpreventionexecuted {
		expGolombWriter.buffer.WriteByte(expGolombWriter.bytes[expGolombWriter.m_nStartingbyte])
		expGolombWriter.bytes[expGolombWriter.m_nStartingbyte] = 0
		expGolombWriter.m_nStartingbyte++
		expGolombWriter.m_nStartingbyte = expGolombWriter.m_nStartingbyte % BUFFER_SIZE_BYTES
		expGolombWriter.m_nLastbitinbuffer = expGolombWriter.m_nLastbitinbuffer - 8
	}

	return nil
}

// clearbuffer 重置bytes
func (expGolombWriter *ExpGolombWriter) resetBytes() {
	expGolombWriter.bytes = make([]byte, BUFFER_SIZE_BITS)
	expGolombWriter.m_nLastbitinbuffer = 0
	expGolombWriter.m_nStartingbyte = 0
}

// dobytealign 对齐
func (expGolombWriter *ExpGolombWriter) dobytealign() {
	nr := expGolombWriter.m_nLastbitinbuffer % 8
	if (nr % 8) != 0 {
		expGolombWriter.m_nLastbitinbuffer = expGolombWriter.m_nLastbitinbuffer + (8 - nr)
	}
}

// Bytes bytes
func (expGolombWriter *ExpGolombWriter) Bytes() []byte {
	expGolombWriter.dobytealign()
	for expGolombWriter.m_nLastbitinbuffer != 0 {
		expGolombWriter.savebufferbyte(true)
	}
	return expGolombWriter.buffer.Bytes()
}
