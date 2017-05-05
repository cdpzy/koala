package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/doublemo/koala/cluster"
	"github.com/doublemo/koala/svc"
	cli "gopkg.in/urfave/cli.v2"
	"gopkg.in/urfave/cli.v2/altsrc"
)

func main() {
	app := &cli.App{
		Name:    "Agent",
		Usage:   "a agent server",
		Version: "4.0.2",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "data/conf/agent.conf",
				Usage:   "Load configuration from `FILE`",
			},

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.name",
				Value: "agent",
				Usage: "网关名称",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.tcp.addr",
				Value: ":19020",
				Usage: "Listen TCP `ADDR`",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.tcp.readbuffersize",
				Value: 32767,
				Usage: "TCP读取缓存大小",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.tcp.writebuffersize",
				Value: 32767,
				Usage: "TCP写入缓存大小",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.tcp.readdeadline",
				Value: 30,
				Usage: "TCP读取超时",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.tcp.writedeadline",
				Value: 30,
				Usage: "TCP写入超时",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.http.addr",
				Value: "",
				Usage: "Listen HTTP `ADDR`",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.http.readtimeout",
				Value: 0,
				Usage: "http读取超时",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.http.writetimeout",
				Value: 0,
				Usage: "http写入超时",
			}),

			altsrc.NewBoolFlag(&cli.BoolFlag{
				Name:  "server.http.ssl",
				Value: false,
				Usage: "SSL",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.http.sslkey",
				Value: "",
				Usage: "ETCD集群地址",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.http.sslcert",
				Value: "",
				Usage: "ETCD集群地址",
			}),

			altsrc.NewInt64Flag(&cli.Int64Flag{
				Name:  "server.http.maxrequestsize",
				Value: 0,
				Usage: "最大http内容大小限制",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "server.kcp.addr",
				Value: "",
				Usage: "Listen KCP `ADDR`",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.writebuffersize",
				Value: 32767,
				Usage: "KCP写入缓存大小",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.readbuffersize",
				Value: 32767,
				Usage: "KCP读取缓存大小",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.readdeadline",
				Value: 30,
				Usage: "KCP读取超时",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.writedeadline",
				Value: 30,
				Usage: "KCP写入超时",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.dscp",
				Value: 46,
				Usage: "set DSCP(6bit)",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.sndwnd",
				Value: 32,
				Usage: "per connection UDP send window",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.rcvwnd",
				Value: 32,
				Usage: "per connection UDP recv window",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.nodelay",
				Value: 1,
				Usage: "ikcp_nodelay()",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.interval",
				Value: 20,
				Usage: "ikcp_nodelay()",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.resend",
				Value: 1,
				Usage: "ikcp_nodelay()",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.nc",
				Value: 1,
				Usage: "ikcp_nodelay()",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.kcp.mtu",
				Value: 1280,
				Usage: "MTU of UDP packets, without IP(20) + UDP(8)",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "server.errorlv",
				Value: 5,
				Usage: "错误日志等级",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "cluster.endpoints",
				Value: "127.0.0.1:2379",
				Usage: "ETCD集群地址",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "cluster.dialTimeout",
				Value: 30,
				Usage: "连接ETCD超时时间",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "cluster.etcdurl",
				Value: "/backends",
				Usage: "ETCD存储集群信息前缀",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "cluster.priority",
				Value: 0,
				Usage: "集群优先级",
			}),

			altsrc.NewStringFlag(&cli.StringFlag{
				Name:  "cluster.addr",
				Value: "",
				Usage: "当前节点服务IP",
			}),

			altsrc.NewIntFlag(&cli.IntFlag{
				Name:  "cluster.port",
				Value: 6061,
				Usage: "当前节点服务端口",
			}),
		},
		Action: Action,
	}

	// Before
	app.Before = altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewTomlSourceFromFlagFunc("config"))
	app.Run(os.Args)
}

// Action 程序启动
func Action(c *cli.Context) error {
	cli.ShowVersion(c)
	log.Println("server.name", c.String("server.name"))
	log.Println("server.tcp.addr", c.String("server.tcp.addr"))
	log.Println("server.tcp.readbuffersize", c.Int("server.tcp.readbuffersize"))
	log.Println("server.tcp.writebuffersize", c.Int("server.tcp.writebuffersize"))
	log.Println("server.tcp.readdeadline", c.Int("server.tcp.readdeadline"))
	log.Println("server.tcp.writedeadline", c.Int("server.tcp.writedeadline"))
	log.Println("server.http.readtimeout", c.Int("server.http.readtimeout"))
	log.Println("server.http.writetimeout", c.Int("server.http.writetimeout"))
	log.Println("server.http.ssl", c.Bool("server.http.ssl"))
	log.Println("server.http.sslkey", c.String("server.http.sslkey"))
	log.Println("server.http.sslcert", c.String("server.http.sslcert"))
	log.Println("server.http.maxrequestsize", c.Int64("server.http.maxrequestsize"))
	log.Println("server.kcp.addr", c.String("server.kcp.addr"))
	log.Println("server.kcp.readbuffersize", c.Int("server.kcp.readbuffersize"))
	log.Println("server.kcp.writebuffersize", c.Int("server.kcp.writebuffersize"))
	log.Println("server.kcp.readdeadline", c.Int("server.kcp.readdeadline"))
	log.Println("server.kcp.writedeadline", c.Int("server.kcp.writedeadline"))
	log.Println("server.kcp.dscp", c.String("server.kcp.dscp"))
	log.Println("server.kcp.sndwnd", c.String("server.kcp.sndwnd"))
	log.Println("server.kcp.rcvwnd", c.String("server.kcp.rcvwnd"))
	log.Println("server.kcp.nodelay", c.String("server.kcp.nodelay"))
	log.Println("server.kcp.interval", c.String("server.kcp.interval"))
	log.Println("server.kcp.resend", c.String("server.kcp.resend"))
	log.Println("server.kcp.nc", c.String("server.kcp.nc"))
	log.Println("server.kcp.mtu", c.String("server.kcp.mtu"))
	log.Println("server.errorlv", c.String("server.errorlv"))
	log.Println("cluster.endpoints:", c.String("cluster.endpoints"))
	log.Println("cluster.dialTimeout", c.Int("cluster.dialTimeout"))
	log.Println("cluster.etcdurl", c.String("cluster.etcdurl"))
	log.Println("cluster.priority", c.Int("cluster.priority"))
	log.Println("cluster.addr", c.String("cluster.addr"))
	log.Println("cluster.port", c.Int("cluster.port"))

	// 日志等级设置
	log.SetLevel(log.Level(c.Int("server.errorlv")))

	// 服务支持
	return svc.WaitExit(start, stop, c)
}

func start(c *cli.Context) error {
	// 集群节点
	cl, err := cluster.New(&cluster.Options{
		Endpoints:   strings.Split(c.String("cluster.endpoints"), ","),
		DialTimeout: c.Int("cluster.dialTimeout"),
		ETCDUrl:     c.String("cluster.etcdurl"),
		Services:    []string{},
	})

	if err != nil {
		log.Fatalf("Cluster init:%v", err)
	}

	// 本地节点信息
	cl.Local.Name = c.String("server.name")
	cl.Local.Type = "agent"
	cl.Local.Priority = c.Int("cluster.priority")
	cl.Local.Status = cluster.NodeStatusOK
	cl.Local.Params = map[string]string{}
	cl.Local.Port = c.Int("cluster.port")

	if len(c.String("cluster.addr")) > 0 {
		cl.Local.Addr = net.ParseIP(c.String("cluster.addr"))
	}

	// 节点上线
	cluster.AddEvent(&cluster.Event{
		Name: cluster.EventNodeOnline,
		Data: make(map[string]interface{}),
		CallBack: func(e cluster.Event) {
			fmt.Println("ONLINE", e.Node)
		},
	})

	// HTTP 服务启动
	go httpServe(c)

	// kcp
	go kcpServe(c)

	// TCP 服务启动
	go tcpServe(c)

	// 加入节点
	cl.Join(cl.Local)

	fmt.Println("started")
	return nil
}

func stop(c *cli.Context) error {
	fmt.Println("Stop;.....")
	return nil
}
