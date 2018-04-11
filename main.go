/*
> 2018年4月10日
> 开始抽时间来完成GUI窗口化的szssp信息化设备管理工具
## 程序流程修改:
- 运行程序
- 自动开始获取设备相关信息并设置临时地址
- 由使用者选择硬盘序列号和网卡MAC地址,查询远程数据库中的信息
- 如果查询到信息就显示出来,如果没有信息就提示使用者提交相关信息
- 然后最终设置网络IP地址.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

const xlsx文件 string = "//33.66.96.14/public/2018Taizhang.xlsx"
const 临时Ip string = "33.66.100.255"

// 定义接口
type 系统 interface {
	设置网络地址(地址 string) error
	获取设备信息() 设备
	获取数据库信息(string, string)
	上传设备信息(信息 设备)
}

// 定义结构
type 设备 struct {
	用户名     string
	部门      string
	硬盘序列号列表 []string
	网卡MAC列表 []string
	IP地址    string
}

// 获取设备信息并初始化窗口中的列表元素
func (p *设备) 获取设备信息() (m 设备) {
	设备.硬盘序列号列表 = 获取硬盘序列号()
	设备.网卡MAC列表 = 获取网卡MAC列表()
	窗口重绘()
}

// 用户点击查询按钮，连接数据库获取相关信息
func (p *设备) 获取数据库信息(硬盘序列号, 网卡MAC地址) {
	// 这里获取窗口中被选中的硬盘序列号和网卡mac地址。
	// 这里打算使用协程，同时在xlsx中查询硬盘和网卡
	// 尝试通过三维数组加快搜索速度
	表格, err := xlsx.FileToSlice(xlsxFile)
	// xlxs对象, err := xlsx.OpenFile(xlsxFile)
	if err != nil {
		// fmt.Println(err)
		// 这里应该在窗口中显示信息，暂时不知道怎么实现
		return
	}
	// 20180411进度留存
	其他信息 := make(chan []string)
	go func(string) {
		var 数据 []string
		for k, v := range 表格 {
			for i, l := range v {
				for o, p := range l {
					println(p)
				}
			}

		}
		其他信息 <- 数据
	}(硬盘序列号, xlxs对象)

	go 查询(网卡MAC地址)
	信息1, 信息2 := <-其他信息, <-其他信息

	return // 无需返回值，因为是使用 *设备 直接操作结构体本身中的元素
}

// 查询成功自动设置网络失败提示用户填写新设备相关信息并提交数据库

func main() {
	//创建window窗口
	//参数一表示创建窗口的样式
	//SW_TITLEBAR 顶层窗口，有标题栏
	//SW_RESIZEABLE 可调整大小
	//SW_CONTROLS 有最小/最大按钮
	//SW_MAIN 应用程序主窗口，关闭后其他所有窗口也会关闭
	//SW_ENABLE_DEBUG 可以调试
	//参数二表示创建窗口的矩形
	w, err := window.New(sciter.SW_TITLEBAR|
		//sciter.SW_RESIZEABLE|
		sciter.SW_CONTROLS|
		sciter.SW_MAIN,
		//sciter.SW_ENABLE_DEBUG,
		// 设置窗口大小
		&sciter.Rect{Left: 0, Top: 0, Right: 300, Bottom: 340})
	if err != nil {
		log.Fatal(err)
	}
	//加载文件
	w.LoadFile("demo1.html")
	//设置标题
	w.SetTitle("三洲特管信息化台账录入系统 v1.0")
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
