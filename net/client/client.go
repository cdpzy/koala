package client

import (
	"crypto/rc4"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/doublemo/koala/helper"
)

// FlagClient 客户端状态
type FlagClient int32

const (
	FlagClientKeyexcg    FlagClient = 0x1 // 是否已经交换完毕KEY
	FlagClientEncrypt    FlagClient = 0x2 // 是否可以开始加密
	FlagClientKickedOut  FlagClient = 0x4 // 踢掉
	FlagClientAuthorized FlagClient = 0x8 // 已授权访问
)

// 错误信息
var (
	ErrorKickedOut = errors.New("ErrorKickedOut")
	ErrorChanFull  = errors.New("ErrorChanFull")
)

// 路由类型方法
type RouterFunc func(*Client, []byte) ([]byte, error)

// Client 客户端处理
type Client struct {
	IP             net.IP        // 客户端IP
	Port           string        //
	ID             string        // 客户端ID
	Encoder        *rc4.Cipher   // 加密器
	Decoder        *rc4.Cipher   // 解密器
	Flag           FlagClient    // 会话标记
	ConnectTime    time.Time     // TCP链接建立时间
	PacketTime     time.Time     // 当前包的到达时间
	LastPacketTime time.Time     // 前一个包到达时间
	CreateAt       time.Time     // 客户端连接时间
	RpmLimit       int           // 客户发包控制
	PacketCount    int           // 对收到的包进行计数，避免恶意发包
	PacketCountRPM int           //
	Params         *Params       //
	in             chan []byte   //
	pending        chan []byte   //
	inputReadyed   chan struct{} //
	outputReadyed  chan struct{} //
	cache          []byte        //
	die            chan struct{} // 会话关闭信号
	closed         chan struct{} // 已关闭, 要求退出
	conn           net.Conn      //
	RouteFunc      RouterFunc    //
	OnBeforeClose  *Evt          //
	OnAfertClose   *Evt          //
}

// WriteIn 写入
func (c *Client) WriteIn(b []byte) error {
	if c.Flag&FlagClientKickedOut != 0 || c.die == nil {
		return ErrorKickedOut
	}

	select {
	case c.in <- b:
	case <-c.closed:
		return ErrorKickedOut
	default:
		return ErrorChanFull
	}

	return nil
}

func (c *Client) input() {
	defer helper.RecoverStack()
	defer func() {
		log.Infoln("Client input closed")
	}()

	log.Infoln("Client input readyed")

	c.ConnectTime = time.Now()
	c.LastPacketTime = time.Now()
	heartbeaterTimer := time.After(time.Minute)
	close(c.inputReadyed)
	for {
		select {
		case msg, ok := <-c.in:
			if !ok {
				return
			}

			c.PacketCount++
			c.PacketCountRPM++
			c.PacketTime = time.Now()
			if r := c.call(msg); r != nil {
				c.Send(r)
			}

			c.LastPacketTime = c.PacketTime

		case <-heartbeaterTimer:
			if c.PacketCountRPM > c.RpmLimit {
				c.Flag |= FlagClientKickedOut
				log.WithFields(log.Fields{"ID": c.ID, "rpm": c.PacketCountRPM, "total": c.PacketCount}).Error("RPM")
			}

			c.PacketCountRPM = 0
			heartbeaterTimer = time.After(time.Minute)

		case <-c.die:
			c.Flag |= FlagClientKickedOut
		}

		// kicked out
		if c.Flag&FlagClientKickedOut != 0 {
			return
		}
	}
}

func (c *Client) output() {
	defer helper.RecoverStack()
	defer func() {
		log.Infoln("Client output closed")
	}()

	log.Infoln("Client output readyed")
	close(c.outputReadyed)
	for {

		select {
		case data := <-c.pending:
			c.send(data)

		case <-c.die:
			c.Flag |= FlagClientKickedOut
		}

		// kicked out
		if c.Flag&FlagClientKickedOut != 0 {
			return
		}
	}
}

// call 响应客户请求
func (c *Client) call(b []byte) []byte {
	stime := time.Now()
	defer helper.RecoverStack(c, b)

	// 解密
	if c.Flag&FlagClientEncrypt != 0 {
		c.Decoder.XORKeyStream(b, b)
	}

	ret, err := c.RouteFunc(c, b)
	if err == ErrorKickedOut {
		c.Flag |= FlagClientKickedOut
		return nil
	}

	if err != nil {
		log.Errorf("RouteFunc:%v", err)
		return nil
	}

	etime := time.Now().Sub(stime)
	log.WithFields(log.Fields{"cost": etime, "code": b}).Debug("REQ")
	return ret
}

// Send 发送
func (c *Client) Send(b []byte) error {
	if b == nil {
		return nil
	}

	log.Debug("Send :", b, c.Flag&FlagClientEncrypt != 0)

	if c.Flag&FlagClientEncrypt != 0 {
		c.Encoder.XORKeyStream(b, b)
	} else if c.Flag&FlagClientKeyexcg != 0 {
		c.Flag &^= FlagClientKeyexcg
		c.Flag |= FlagClientEncrypt
	}

	select {
	case c.pending <- b:
	case <-time.After(30 * time.Second):
		log.WithFields(log.Fields{"ID": c.ID, "ip": c.IP}).Warning("pending full")
		return ErrorChanFull
	}

	return nil
}

// send 发送到客户端
func (c *Client) send(b []byte) bool {
	size := len(b)
	binary.BigEndian.PutUint16(c.cache, uint16(size))
	copy(c.cache[2:], b)
	log.Debug("Send To Client :", b, c.cache[:size+2])
	n, err := c.conn.Write(c.cache[:size+2])
	if err != nil {
		log.Warningf("Error send reply data, bytes: %v reason: %v", n, err)
		return false
	}
	return true
}

// Close //
func (c *Client) Close() {
	if c.OnBeforeClose.Count() > 0 {
		c.OnBeforeClose.Iterator(func(k string, f EvtCallBack) bool {
			f(c)
			return true
		})
	}

	if c.die == nil {
		return
	}

	close(c.die)
	c.Flag |= FlagClientKickedOut
	if c.OnAfertClose.Count() > 0 {
		c.OnAfertClose.Iterator(func(k string, f EvtCallBack) bool {
			f(c)
			return true
		})
	}
}

// NewClient 创建客户端
func NewClient(conn net.Conn, op *Config) *Client {
	c := &Client{
		Params:        NewParams(),
		in:            make(chan []byte),
		pending:       make(chan []byte, op.PendingSize),
		cache:         make([]byte, 65535),
		die:           make(chan struct{}),
		conn:          conn,
		CreateAt:      time.Now(),
		RpmLimit:      op.RpmLimit,
		OnBeforeClose: NewEvt(),
		OnAfertClose:  NewEvt(),
		RouteFunc:     op.RouteFunc,
		closed:        op.Closed,
	}

	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err == nil {
		c.Port = port
		c.IP = net.ParseIP(host)
	}

	c.inputReadyed = make(chan struct{})
	go c.input()
	<-c.inputReadyed

	c.outputReadyed = make(chan struct{})
	go c.output()
	<-c.outputReadyed
	return c
}

// HandleClient 默认连接处理
func HandleClient(conn net.Conn, readDeadline, writeDeadline int, op *Config) {
	defer helper.RecoverStack()

	header := make([]byte, 2)
	client := NewClient(conn, op)

	Register(client)
	defer func() {
		client.Close()
		Unregister(client.ID)
		log.WithFields(log.Fields{"ID": client.ID, "IP": client.IP, "ONLINE:": Count()}).Debug("Client shutdown")
	}()

	log.WithFields(log.Fields{"host": client.IP.String(), "port": client.Port, "ONLINE:": Count()}).Debug("new connection from")

	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(readDeadline) * time.Second))
		conn.SetWriteDeadline(time.Now().Add(time.Duration(writeDeadline) * time.Second))

		n, err := io.ReadFull(conn, header)
		if err == io.EOF {
			return
		}

		if err != nil {
			log.Warningf("read header failed, ip:%v reason:%v size:%v", client.IP, err, n)
			return
		}

		size := binary.BigEndian.Uint16(header)
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warningf("read payload failed, ip:%v reason:%v size:%v", client.IP, err, n)
			return
		}

		err = client.WriteIn(payload)
		if err == ErrorKickedOut {
			return
		} else if err == ErrorChanFull {
			log.Warningf("HandleClient, ErrorChanFull: ID %v IP %s", client.Params, client.IP.String())
		}
	}
}
