package client

// Config 配置
type Config struct {
	RpmLimit      int           // 客户发包控制
	RouteFunc     RouterFunc    // 路由
	OnBeforeClose *Evt          // 客户端关闭以前事件处理
	OnAfertClose  *Evt          // 客户端关闭以后事件处理
	PendingSize   int           // 等待通道缓冲队列大小
	Closed        chan struct{} //
}

// NewConfig 配置
func NewConfig() *Config {
	return &Config{
		RpmLimit:      300,
		OnBeforeClose: NewEvt(),
		OnAfertClose:  NewEvt(),
		PendingSize:   128,
	}
}
