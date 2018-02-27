/*
运行流程:
1. 启动获取硬件设备
2. 设置临时ip地址,访问局域网服务器中的xlsx文件
3. 搜索mac地址,获取正确ip并设置,同时对比其他信息是否正确
4. 如果有信息出现差异,提示并修改xlsx文件.包含用户和部门.
5. 程序退出.
*/
package main

import (
	"fmt"
	"os"
)

// 临时Ip地址常量
const tmpIp = "33.66.100.255"
