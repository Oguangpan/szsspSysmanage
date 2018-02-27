package devinfo

// 设备接口 三个方法 获取系统类型 获取mac地址 获取硬盘序列号
type Systemer interface {
	getos() string
	getmac() []string
	gethardid() []string
}
