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
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	//"strconv"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/tealeg/xlsx"
)

const xlsxFile string = "Z:\a\2018台账.xlsx"
const tempIp string = "33.66.100.255"

// 定义接口
type systemer interface {
	setIp()
	getdeviceInfo()
	getDbinfo()
	updateDeviceInfo()
}

// 定义结构
type thisComputer struct {
	userName       string
	department     string
	hardDiskNumber []string
	macs           []map[string]string
	ip             string
	osType         string
}

// 获取设备信息并初始化窗口中的列表元素
func (p *thisComputer) getdeviceInfo() (h []string, c []map[string]string) {
	// hardDiskNumber
	ids := runCmd("wmic diskdrive get serialnumber")
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		h = append(h, j)
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
			c = append(c, tmp)
		}
	}

	// 系统版本号
	s := runCmd("ver")
	if strings.Contains(s, "10") {
		p.osType = "win10"
	} else if strings.Contains(s, "6.1") {
		p.osType = "win7"
	} else if strings.Contains(s, "5.1") {
		p.osType = "winxp"
	} else {
		p.osType = "null"
	}

	return
}

func searchIp(t string, xlsxObjects *xlsx.File, searchIpChan chan []string) {
	var tlist []string
	for k, v := range xlsxObjects.Sheets[0].Rows {
		for _, l := range v.Cells {
			if l.Value == t {
				for _, ce := range xlsxObjects.Sheets[0].Rows[k].Cells {
					//println(ce.Value)
					tlist = append(tlist, ce.Value)
				}
			}
		}
	}

	searchIpChan <- tlist

}

// 用户点击查询按钮，连接数据库获取相关信息
// 这里获取窗口中被选中的硬盘序列号和网卡mac地址。
// 这里打算使用协程，同时在xlsx中查询硬盘和网卡
// 尝试通过三维数组加快搜索速度
func (p *thisComputer) getDbinfo() {

	// 这两个变量内容从窗口中的两个选择框获取
	var hdId string
	var cMac string

	xlsxObjects, _ := xlsx.OpenFile(xlsxFile)
	info := make(chan []string)

	go searchIp(hdId, xlsxObjects, info)
	go searchIp(cMac, xlsxObjects, info)
	hdTarget, cmacTarget := <-info, <-info
	switch {
	case len(hdTarget) > 0:
		p.userName = hdTarget[3]
		p.department = hdTarget[2]
		p.ip = hdTarget[10]
	case len(cmacTarget) > 0:
		p.userName = cmacTarget[3]
		p.department = cmacTarget[2]
		p.ip = cmacTarget[10]
	default:
		// 两种数据都没有查询到设备在服务器中的记录信息
	}

	return // 无需返回值，因为是使用 *thisComputer 直接操作结构体本身中的元素
}

// 调用系统CMD命令执行外部程序
func runCmd(s string) (echo string) {
	t, _ := exec.Command("cmd", "/C", s).Output()
	// echo = ConvertToString(string(t), "gbk", "utf-8")
	echo = string(t)
	return
}

// 查询成功自动设置网络失败提示用户填写新设备相关信息并提交数据库
func (p *thisComputer) setIp(name string, ip string) {
	var ipstr string
	var dnsstr string
	// 判断当前操作系统版本 避免出现兼容性问题

	switch p.osType {
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
		os.Exit(0)
	}
	// 调用系统命令行进行网络设置
	runCmd(ipstr)
	runCmd(dnsstr)

	fmt.Println("网络设置完毕")
	return

}

func setIpButtonOnclick(root *sciter.Element) {
	btn1, _ := root.SelectById("btn1")
	btn1.OnClick(func() {
		fmt.Println("btn1被点击 了")
	})
}
func getInfoButtonOnclick(root *sciter.Element) {
	btn2, _ := root.SelectById("btn2")
	btn2.OnClick(func() {
		fmt.Println("btn2被点击 了")
	})
}
func UpdateButtonOnclick(root *sciter.Element) {
	btn3, _ := root.SelectById("btn3")
	btn3.OnClick(func() {
		fmt.Println("btn3被点击 了")
	})
}

func closeWindow(root *sciter.Element) {
	closeBtn, _ := root.SelectById("closebtn")
	closeBtn.OnClick(func() {
		os.Exit(0)
	})
}

func main() {
	w, err := window.New(sciter.SW_TITLEBAR|sciter.SW_CONTROLS|sciter.SW_MAIN, &sciter.Rect{Left: 0, Top: 0, Right: 720, Bottom: 340})
	if err != nil {
		log.Fatal(err)
	}
	w.LoadFile("newgui.html")
	w.SetTitle("三洲特管信息化台账录入系统 v1.0")

	root, _ := w.GetRootElement()
	setIpButtonOnclick(root)
	getInfoButtonOnclick(root)
	UpdateButtonOnclick(root)
	closeWindow(root)

	//var cmp systemer
	cmp := new(thisComputer)
	hds, mas := cmp.getdeviceInfo()
	// 添加硬盘们的序列号到列表选择框中
	set1, _ := root.SelectById("slHdn")
	set2, _ := root.SelectById("slNet")
	for _, j := range hds {
		set1.CallFunction("addOp", sciter.NewValue(j))
	}
	for _, j := range mas {
		set2.CallFunction("addMac", sciter.NewValue(j["Name"]+":"+j["Mac"]))
	}
	// 添加网卡MAC地址列表到列表选择框中
	w.Show()
	w.Run()
}
