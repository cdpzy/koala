package daemon

type Config struct {
	Data
	DataDir    string // 数据存储目录
	UpdateTime int    // 更新心跳时间(s)
}
