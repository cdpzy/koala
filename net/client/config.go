package client

// Config 配置
type Config struct {
	RpmLimit      int                   // 客户发包控制
	RouteFunc     RouterFunc            // 路由
	OnBeforeClose map[string]BeforeFunc // 客户端关闭以前事件处理
	OnAfertClose  map[string]AfertFunc  // 客户端关闭以后事件处理
	PendingSize   int                   // 等待通道缓冲队列大小
}

// NewConfig 配置
func NewConfig() *Config {
	return &Config{
		RpmLimit:      300,
		OnBeforeClose: make(map[string]BeforeFunc),
		OnAfertClose:  make(map[string]AfertFunc),
		PendingSize:   128,
	}
}
