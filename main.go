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
	"github.com/tealeg/xlsx"
)

const xlsxFile string = "//33.66.96.14/public/2018台账.xlsx"
const 临时Ip string = "33.66.100.255"

// 定义接口
type 系统 interface {
	设置网络地址()
	获取设备信息()
	获取数据库信息()
	上传设备信息(信息 设备)
}

// 定义结构
type 设备 struct {
	用户名     string
	部门      string
	硬盘序列号列表 []string
	网卡MAC列表 []map[string]string
	IP地址    string
	系统类型    string
}

// 获取设备信息并初始化窗口中的列表元素
func (p *设备) 获取设备信息() (硬盘 []string, 网卡 []map[string]string) {
	// 硬盘序列号列表
	ids := 运行命令("wmic diskdrive get serialnumber")
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		硬盘 = append(硬盘, j)
	}

	// 网卡信息列表
	intf, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, v := range intf {
		tmp := make(map[string]string)
		tmp["Name"] = v.Name
		tmp["Mac"] = v.HardwareAddr.String()
		if tmp["Mac"] != "" {
			网卡 = append(网卡, tmp)
		}
	}

	// 系统版本号
	s := 运行命令("ver")
	if string.Contains(s, "10") {
		p.系统类型 = "win10"
	} else if string.Contains(s, "6.1") {
		p.系统类型 = "win7"
	} else if string.Contains(s, "5.1") {
		p.系统类型 = "winxp"
	} else {
		p.系统类型 = "null"
	}

	return
}

func 查号(目标 string, 表格 *xlsx.File, 查号信道 chan []string) {
	var tlist []string
	for k, v := range 表格.Sheets[0].Rows {
		for _, l := range v.Cells {
			if l.Value == 目标 {
				for _, ce := range 表格.Sheets[0].Rows[k].Cells {
					//println(ce.Value)
					tlist = append(tlist, ce.Value)
				}
			}
		}
	}

	查号信道 <- tlist

}

// 用户点击查询按钮，连接数据库获取相关信息
// 这里获取窗口中被选中的硬盘序列号和网卡mac地址。
// 这里打算使用协程，同时在xlsx中查询硬盘和网卡
// 尝试通过三维数组加快搜索速度
func (p *设备) 获取数据库信息() {

	// 这两个变量内容从窗口中的两个选择框获取
	var 硬盘序列号 string
	var 网卡MAC地址 string

	表格, _ := xlsx.OpenFile(xlsxFile)
	info := make(chan []string)

	go 查号(硬盘序列号, 表格, info)
	go 查号(网卡MAC地址, 表格, info)
	通过硬盘找到的目标, 通过网卡找到的目标 := <-info, <-info
	switch {
	case len(通过硬盘找到的目标) > 0:
		p.用户名 = 通过硬盘找到的目标[3]
		p.部门 = 通过硬盘找到的目标[2]
		p.IP地址 = 通过硬盘找到的目标[10]
	case len(通过网卡找到的目标) > 0:
		p.用户名 = 通过网卡找到的目标[3]
		p.部门 = 通过网卡找到的目标[2]
		p.IP地址 = 通过网卡找到的目标[10]
	default:
		// 两种数据都没有查询到设备在服务器中的记录信息
	}

	return // 无需返回值，因为是使用 *设备 直接操作结构体本身中的元素
}

// 调用系统CMD命令执行外部程序
func 运行命令(s string) (echo string) {
	t, _ := exec.Command("cmd", "/C", s).Output()
	// echo = ConvertToString(string(t), "gbk", "utf-8")
	echo = t
	return
}

// 查询成功自动设置网络失败提示用户填写新设备相关信息并提交数据库
func (p *设备) 设置网络地址(name string, ip string) {
	var ipstr string
	var dnsstr string
	// 判断当前操作系统版本 避免出现兼容性问题

	switch p.系统类型 {
	case "winxp":
		ipstr = strings.Join([]string{"netsh interface ip set address",
			"name=\"" + name + "\"",
			"source=static",
			"addr=" + ip,
			"mask=255.255.224.0 gateway=33.66.99.169 gwmetric=0"}, " ")

		dnsstr = strings.Join([]string{"netsh interface ip set dns",
			"name=\"" + name + "\"",
			"source=static addr=1.1.1.1 register=primary"}, " ")

	case "win10":
		ipstr = strings.Join([]string{"netsh interface ip set address",
			"name=\"" + name + "\"",
			"source=static",
			"address=" + ip,
			"mask=255.255.225.0 gateway=33.66.99.169"}, " ")

		dnsstr = strings.Join([]string{"netsh interface ip add dnsservers",
			name,
			"address=1.1.1.1"}, " ")

	case "win7":
		ipstr = strings.Join([]string{"netsh interface ip set address",
			"name=\"" + name + "\"",
			"source=static",
			"address=" + ip,
			"mask=255.255.225.0 gateway=33.66.99.169"}, " ")

		dnsstr = strings.Join([]string{"netsh interface ip add dnsservers",
			name,
			"address=1.1.1.1"}, " ")
	default:
		fmt.Println("未知的操作系统.")
		os.Exit()
	}
	// 调用系统命令行进行网络设置
	运行命令(ipstr)
	运行命令(dnsstr)

	fmt.Println("网络设置完毕")
	return

}

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

	/*
		三个按钮
		1. 设置临时IP地址,根据网卡选择框中的选项来设置
		2. 对比服务器数据,根据获取到的本机的网卡和硬盘序列号来在服务器数据中查找其他数据. 显示在窗口中
		3. 上传数据按钮, 用于将当前窗口中所有编辑框中和选择框中的所有数据全部上传到服务器数据库中
	*/

	//启动前的准备工作,获取设备信息并修改页面内容
	//设备.硬盘序列号列表 = 获取硬盘序列号()
	//设备.网卡MAC列表 = 获取网卡MAC列表()
	//窗口重绘()
	//加载文件
	w.LoadFile("demo1.html")
	//设置标题
	w.SetTitle("三洲特管信息化台账录入系统 v1.0")
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
