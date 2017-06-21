package daemon

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 命令定义
const (
	CommandSetRestart int16 = iota + 1 // 要求重启
)

// DaemonBody 消息结构
type DaemonBody struct {
	Command int16 // 消息命令
}

// Daemon 进程守护系统
type Daemon struct {
	Data       *Data           // 数据
	DataDir    string          // 进程监控文件写入目录
	UpdateTime time.Duration   //
	pending    chan DaemonBody // 消息输入通道
	die        chan struct{}   //
}

// Init 初始化
func (d *Daemon) Init() error {
	if err := d.update(); err != nil {
		return err
	}
	return nil
}

func (d *Daemon) Serve() {
	d.die = make(chan struct{})
	timer := time.After(time.Second * d.UpdateTime)
	for {
		select {
		case frame, ok := <-d.pending:
			if !ok {
				return
			}

			switch frame.Command {
			case CommandSetRestart:
				d.Data.Command = CommandSetRestart
				d.update()
			}

		case <-timer:
			timer = time.After(time.Second * d.UpdateTime)
			d.update()

		case <-d.die:
			return
		}
	}
}

func (d *Daemon) update() error {
	d.Data.DateUinx = time.Now().Unix()
	path, err := filepath.Abs(d.DataDir)
	if err != nil {
		return err
	}

	if f, err := os.Stat(path); os.IsNotExist(err) || !f.IsDir() {
		return fmt.Errorf("Invalid path:%s", d.DataDir)
	}

	file, err := filepath.Abs(d.DataDir + "/" + d.name())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	b, err := d.Data.Encode()
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

func (d *Daemon) Send(m DaemonBody) {
	select {
	case d.pending <- m:
	default:
		log.Println("Daemon pending chan full.")
	}
}

func (d *Daemon) Close() {
	if d.die == nil {
		return
	}

	close(d.die)
}

func (d *Daemon) name() string {
	return fmt.Sprintf("%s", d.Data.NodeName)
}

func NewDaemon(op *Config) *Daemon {
	return &Daemon{
		Data:       &op.Data,
		DataDir:    op.DataDir,
		UpdateTime: time.Duration(op.UpdateTime),
		pending:    make(chan DaemonBody, 1),
	}
}
