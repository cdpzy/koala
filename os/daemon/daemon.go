package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
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

			fmt.Println("OKK", frame)

		case <-timer:
			timer = time.After(time.Second * d.UpdateTime)
			d.update()

		case <-d.die:
			return
		}
	}
}

func (d *Daemon) update() error {
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
	_, err = f.Write(d.Data.Encode())
	return err
}

func (d *Daemon) name() string {
	return fmt.Sprintf("%s-%d", d.Data.NodeName, d.Data.Pid)
}
