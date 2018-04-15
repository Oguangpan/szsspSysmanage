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
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
	"github.com/tealeg/xlsx"
)

// 存放数据的网络服务器中的表格文件
const xlsxFile string = "//33.66.96.14/public/2018台账.xlsx"

// 如果没有固定的ip,设置一个临时的ip,便于访问网络服务器
const tempIp string = "33.66.100.255"

// 定义接口
type systemer interface {
	getdeviceInfo()
	getDbinfo()
	updateDeviceInfo()
}

// 定义结构
type thisComputer struct {
	userName        string
	department      string
	hardDiskNumbers []string
	harddisk        string
	macs            []map[string]string
	mac             string
	ip              string
	osType          string
}

// 获取设备信息并初始化窗口中的列表元素
func (p *thisComputer) getdeviceInfo() (h []string, c []map[string]string) {
	// 读取硬盘列表
	ids := runCmd("wmic diskdrive get serialnumber")
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		h = append(h, j)
	}

	// 读取网卡信息列表
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

// 查询服务器上xlsx表格中的IP地址
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
func (p *thisComputer) getDbinfo(xlsxObjects *xlsx.File) {

	// 这两个变量内容从窗口中的两个选择框获取
	var hdId string = p.harddisk
	var cMac string = p.mac

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
		p.userName = "查询失败"
		p.department = "查询失败"
		p.ip = "查询失败"
	}

	return // 无需返回值，因为是使用 *thisComputer 直接操作结构体本身中的元素
}

// 上传信息

func (p *thisComputer) updateDeviceInfo(xlsxObjects *xlsx.File) {
	row := xlsxObjects.Sheets[0].AddRow()
	// 向新生成的行row 插入内容
	devInfo := []string{" ", " ",
		p.department,
		p.userName,
		" ", " ", " ", " ", " ",
		p.harddisk,
		p.ip,
		p.mac, " "}
	row.WriteSlice(&devInfo, -1)
	xlsxObjects.Save(xlsxFile)
}

// 调用系统CMD命令执行外部程序
func runCmd(s string) (echo string) {
	t, _ := exec.Command("cmd", "/C", s).Output()
	// echo = ConvertToString(string(t), "gbk", "utf-8")
	echo = string(t)
	return
}

// 根据系统类型设置ip地址
func setIp(name string, ip string, ostype string) {
	var ipstr string
	var dnsstr string
	// 判断当前操作系统版本 避免出现兼容性问题

	switch ostype {
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
			"mask=255.255.224.0 gateway=33.66.99.169"}, " ")

		dnsstr = strings.Join([]string{"netsh interface ip add dnsservers",
			name,
			"address=1.1.1.1"}, " ")

	case "win7":
		ipstr = strings.Join([]string{"netsh interface ip set address",
			"name=\"" + name + "\"",
			"source=static",
			"address=" + ip,
			"mask=255.255.224.0 gateway=33.66.99.169"}, " ")

		dnsstr = strings.Join([]string{"netsh interface ip add dnsservers",
			name,
			"address=1.1.1.1"}, " ")
	default:

		os.Exit(0)
	}
	// 设置ip 掩码 网关
	runCmd(ipstr)
	// 设置DNS
	runCmd(dnsstr)

	return

}

// 获取窗口上两个下拉列表框中当前选中的内容
func (p *thisComputer) getWindowSelectValue(root *sciter.Element) (cardName string) {
	// 选择网卡下拉列表框获取当前选中的值
	editNetcad, _ := root.SelectFirst(".right>.label>#slNet")
	v, _ := editNetcad.GetValue()
	// 分离出网络名称和网卡mac地址
	netCadinfo := strings.Split(v.String(), "|")
	// 给 结构体 中的网卡属性赋值 点击查询按钮的时候会用到
	p.mac = netCadinfo[1]

	cardName = netCadinfo[0]
	editHds, _ := root.SelectFirst(".right>.label>#slHdn")
	k, _ := editHds.GetValue()
	p.harddisk = k.String()
	return
}

// 设置IP按钮被点击事件
func (p *thisComputer) setIpButtonOnclick(root *sciter.Element) {
	// 测试,获取被选择的值
	btn1, _ := root.SelectById("btn1")
	// 按钮点击事件
	btn1.OnClick(func() {
		v := p.getWindowSelectValue(root)
		// 使用网络名 设置网络地址 首先判断当前 ip编辑框中是否有地址
		editIp, _ := root.SelectFirst(".right>.label>#eIp")
		ip, _ := editIp.GetValue()
		//根据ip编辑框中是否有ip存在, 选择设置临时ip还是编辑框中的ip
		if ip.String() != "" {
			go setIp(v, ip.String(), p.osType)

		} else {
			go setIp(v, tempIp, p.osType)
		}
		btn1.CallFunction("popmsgbox", sciter.NewValue("ip地址设置完毕，请稍等片刻进行下一步操作。"))
	})
}

// 通过两个下拉列表框中的被选中项查询其他信息
func (p *thisComputer) getInfoButtonOnclick(root *sciter.Element, xlsxObjects *xlsx.File) {
	btn2, _ := root.SelectById("btn2")
	btn2.OnClick(func() {
		// 将窗口中下拉选择框中的选中项的值赋于自我对应属性
		p.getWindowSelectValue(root)
		// 查询信息
		p.getDbinfo(xlsxObjects)
		// 把信息填入到窗口中的对应编辑框中
		editName, _ := root.SelectFirst("#eName")
		editGroup, _ := root.SelectFirst("#eGroup")
		editIp, _ := root.SelectFirst("#eIp")
		editName.SetValue(sciter.NewValue(p.userName))
		editGroup.SetValue(sciter.NewValue(p.department))
		editIp.SetValue(sciter.NewValue(p.ip))

	})
}

// 上传数据到服务器中的xlsx表格中
func (p *thisComputer) UpdateButtonOnclick(root *sciter.Element, xlsxObjects *xlsx.File) {
	btn3, _ := root.SelectById("btn3")
	btn3.OnClick(func() {
		p.updateDeviceInfo(xlsxObjects)
		btn3.CallFunction("popmsgbox", sciter.NewValue("上传数据完毕。"))
	})
}

// 窗口右上角的关闭按钮被点击事件
func (p *thisComputer) closeWindow(root *sciter.Element) {
	closeBtn, _ := root.SelectById("closebtn")
	closeBtn.OnClick(func() {
		os.Exit(0)
	})
}

func newWindowtextSet(root *sciter.Element, hds []string, mas []map[string]string, cmp *thisComputer) {
	setHd, _ := root.SelectById("slHdn")
	setNe, _ := root.SelectById("slNet")
	// 添加硬盘们的序列号到列表选择框中
	for _, j := range hds {
		setHd.CallFunction("addOp", sciter.NewValue(j))
	}
	// 添加网卡MAC地址列表到列表选择框中
	for _, j := range mas {
		setNe.CallFunction("addMac", sciter.NewValue(j["Name"]+"|"+j["Mac"]))
	}
	// 更新版本编辑框中的设备系统版本数据
	editOstype, _ := root.SelectFirst(".right>.label>#eVersion")
	editOstype.SetValue(sciter.NewValue(cmp.osType))
}

func main() {

	//cmp 是主接口
	cmp := new(thisComputer)
	hds, mas := cmp.getdeviceInfo()
	// w 是窗口对象
	w, _ := window.New(sciter.SW_TITLEBAR|sciter.SW_CONTROLS|sciter.SW_MAIN, &sciter.Rect{Left: 0, Top: 0, Right: 720, Bottom: 340})

	// xlsx 对象
	xlsxObjects, err := xlsx.OpenFile(xlsxFile)
	if err != nil {
		os.Exit(1)
	}

	w.LoadFile("newgui.html")
	w.SetTitle("三洲特管信息化台账录入系统 v1.0")
	// 所有按钮的事件响应
	root, _ := w.GetRootElement()
	cmp.setIpButtonOnclick(root)
	cmp.getInfoButtonOnclick(root, xlsxObjects)
	cmp.UpdateButtonOnclick(root, xlsxObjects)
	cmp.closeWindow(root)

	newWindowtextSet(root, hds, mas, cmp)

	w.Show()
	w.Run()
}
