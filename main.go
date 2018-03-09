// 通过查询当前设备部分信息到服务器中获取更多的信息
// 设置当前设备的网络IP地址

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
	"github.com/tealeg/xlsx"
)

const xlsxFile string = "//33.66.96.14/public/2018Taizhang.xlsx"
const tempIp string = "33.66.100.255"

var stdin = bufio.NewReader(os.Stdin)

type ultinvsupgader interface {
	Get() ([]string, []map[string]string)
	Search(s string) (ss map[string]string)
	Add(aOrm string, ss map[string]string)
	Set(name string, ip string)
}
type uisg struct{}

// *uisg.Get()(硬盘序列号切片，网卡键值对切片)
func (p *uisg) Get() (s []string, m []map[string]string) {

	ids_byte, _ := exec.Command("cmd", "/C", "wmic diskdrive get serialnumber").Output()
	ids := string(ids_byte)
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		s = append(s, j)
	}

	intf, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, v := range intf {
		tmp := make(map[string]string)
		tmp["Name"] = v.Name
		tmp["Mac"] = v.HardwareAddr.String()
		if tmp["Mac"] != "" {
			m = append(m, tmp)
		}
	}

	return
}

// *uisg.Search(序列号或mac号)(服务器记录哈希表)
func (p *uisg) Search(s string) (ss map[string]string) {

	ss = map[string]string{"部门": "", "人员": "", "序列号": "", "Ip": "", "Mac": "", "行号": ""}
	xlFile, err := xlsx.OpenFile(xlsxFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "计算机" {
			for column, row := range sheet.Rows {
				for _, cell := range row.Cells {
					if cell.String() == s {
						for i, cell := range row.Cells {
							switch i {
							case 2:
								ss["部门"] = cell.String()
							case 3:
								ss["人员"] = cell.String()
							case 9:
								ss["序列号"] = cell.String()
							case 10:
								ss["Ip"] = cell.String()
							case 11:
								ss["Mac"] = cell.String()
							}
						}
						// 返回数据中包含查到数据的行号，方便修改信息时使用。
						ss["行号"] = strconv.Itoa(column)
					}
				}
			}
		}
	}

	return
}

// *uisg.Add(要修改的硬盘序列号, 硬件设备信息哈希表)
func (p *uisg) Add(aOrm string, ss map[string]string) {

	xlFile, err := xlsx.OpenFile(xlsxFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	if aOrm == "" {
		// 存放硬盘序列号的参数为空，代表需要增加数据
		for _, sheet := range xlFile.Sheets {
			if sheet.Name == "计算机" {
				row := sheet.AddRow()
				s := []string{"-", "-", ss["部门"], ss["人员"], "-", "-", "-", "-",
					"-", ss["序列号"], ss["Ip"], ss["Mac"], "-"}
				for _, v := range s {
					row.AddCell().Value = v
				}
				xlFile.Save(xlsxFile)
				return
			}
		}
	} else {
		// 修改数据
		for _, sheet := range xlFile.Sheets {
			if sheet.Name == "计算机" {
				u, _ := strconv.Atoi(ss["行号"])
				sheet.Cell(u, 10).Value = aOrm
				xlFile.Save(xlsxFile)
				return
			}
		}
	}

}

// 转换GBK到UTF8,针对在windows系统下遇到的exec中文乱码问题。
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

// *uisg.Set(要设置的IP地址)
func (p *uisg) Set(name string, ip string) {

	var nii string = "netsh interface ip "
	var san string = "set address name="
	var ad string = "add dns "

	var setIpstr string = nii + san + name + " static " + ip + " 255.255.224.0 33.66.99.169"
	var setDns1str string = nii + ad + name + " 202.98.96.68 index=1"
	var setDns2str string = nii + ad + name + " 61.139.2.69 index=2"

	rinfo, _ := exec.Command("cmd", "/C", setIpstr).Output()
	fmt.Println(ConvertToString(string(rinfo), "gbk", "utf-8"))
	rinfo, _ = exec.Command("cmd", "/C", setDns1str).Output()
	fmt.Println(ConvertToString(string(rinfo), "gbk", "utf-8"))
	rinfo, _ = exec.Command("cmd", "/C", setDns2str).Output()
	fmt.Println(ConvertToString(string(rinfo), "gbk", "utf-8"))

	fmt.Println("网络设置完毕,请检查.")
	pingStr, _ := exec.Command("cmd", "/C", "ping -n 1 33.66.96.14").Output()
	fmt.Println(ConvertToString(string(pingStr), "gbk", "utf-8"))

	return
}

// 跟用户聊天的方式
func chat(s string) (str string) {
	fmt.Println(s)
	fmt.Fscan(stdin, &str)
	stdin.ReadString('\n')
	return
}

// 主流程中通用的用户交互函数,根据输入的参数执行不同的行为。返回信息表。
func HumanComputerInteraction(s []string, ms []map[string]string) map[string]string {

	ss := map[string]string{"部门": "", "人员": "", "序列号": "", "Ip": "", "Mac": "", "行号": ""}

	if len(s) > 0 {
		for k, v := range s {
			fmt.Println(k, ":", v)
		}
		hddNum, _ := strconv.Atoi(chat("请选择你要添加到服务器记录中的硬盘"))
		ss["序列号"] = s[hddNum]
	}
	if len(ms) > 0 {
		for k, v := range ms {
			fmt.Println(k, ":", v["Mac"])
		}
		netNum, _ := strconv.Atoi(chat("选择你要新添加到服务器记录中的网卡MAC号"))
		ss["Mac"] = ms[netNum]["Mac"]
	}
	ss["人员"] = chat("请输入用户姓名：")
	ss["部门"] = chat("请输入用户部门：")

	return ss

}

func main() {
	var us ultinvsupgader
	us = new(uisg)
	hdds, nets := us.Get()
	var devIt map[string]string
	var cadName string

	if len(nets) > 0 {
		//设置临时地址
		fmt.Println("请选择需要设置IP的网卡名称：")
		for k, v := range nets {
			fmt.Println(k, v["Name"])
		}
		l, _ := strconv.Atoi(chat("请选择需要设置IP的网卡名称："))
		cadName = nets[l]["Name"]
		us.Set(cadName, tempIp)
	}

	if len(hdds) != 0 {
		// 读取到硬盘序列号
		for _, v := range hdds {
			infos := us.Search(v)
			if infos["序列号"] != "" {
				// 查询结果匹配则直接设置ip
				fmt.Println("该设备在服务器中记录的信息如下", infos)
				us.Set(cadName, infos["Ip"])

				return
			}
		}
		// 本地有硬盘序列号,但是数据库中没有硬盘信息,使用网卡地址搜索
		for _, v := range nets {
			infos := us.Search(v["Mac"])
			if infos["Mac"] != "" {
				// 找到信息,直接设置IP，提示增加硬盘信息
				if infos["Ip"] != "" {
					us.Set(cadName, infos["Ip"])
				} else {
					fmt.Println("数据库中发现Mac,但是没有发现Ip地址,请联系管理员尽快更新.")
					return
				}

				fmt.Println("服务器中发现网卡信息但发现硬盘序列号不匹配,请修改更新: ")
				devIt = HumanComputerInteraction(hdds, nets)
				us.Add(devIt["序列号"], infos)
				fmt.Println("已更新硬盘序列号。", devIt["序列号"])
				return
			}
		}
		// 服务器记录中没有设备中的硬盘信息和网卡信息，提示添加。并返回。
		fmt.Println("服务器记录中没有发现该设备的硬盘和网卡信息,请添加到服务器中.")
		devIt = HumanComputerInteraction(hdds, nets)
		us.Add("", devIt)
		return
	} else {
		// 没有从本设备中读取到硬盘序列号，使用网卡搜索信息
		for _, v := range nets {
			infos := us.Search(v["Mac"])
			if infos["Ip"] != "" {
				// 找到信息,返回信息,设置IP
				fmt.Println("无法读取该设备的硬盘序列号,不过已找到服务器中记录的其他信息.",
					infos)
				us.Set(cadName, infos["Ip"])
				return
			}
		}
		// 没有在服务器中找到该设备的网卡信息，提示添加，并返回。
		fmt.Println("该无法读取该设备的硬盘信息,网卡信息在服务器中也没有记录.请添加: ")
		devIt = HumanComputerInteraction(hdds, nets)

		us.Add("", devIt)
		fmt.Println("已添加记录到服务器中, 内容如下: ", devIt)
		return
	}

}
