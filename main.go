// 通过查询当前设备部分信息到服务器中获取更多的信息
// 设置当前设备的网络IP地址

package main

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

const xlsxFile string = "//33.66.96.14/public/2018Taizhang.xlsx"
const tempIp string = "33.66.100.255"

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

// *uisg.Set(要设置的IP地址)
func (p *uisg) Set(name string, ip string) {
	//TODO
}

// 主流程中通用的用户交互函数,根据输入的参数执行不同的行为。返回信息表。
func HumanComputerInteraction(hdds []string, nets []map[string]string) (ss map[string]string) {

	var w int

	if len(hdds) > 0 {
		fmt.Println("选择你要新添加的硬盘序列号。")
		for k, v := range hdds {
			fmt.Println(k, ":", v)
		}

		fmt.Scanf("%d\n", w)
		ss["序列号"] = hdds[w]
	}
	if len(nets) > 0 {
		fmt.Println("选择你要新添加的网卡MAC号。")
		for k, v := range nets {
			fmt.Println(k, ":", v["Mac"])
		}
		fmt.Scanf("%d\n", w)
		ss["Mac"] = nets[w]["Mac"]
	}

	fmt.Scanf("%d\n", w)
	ss["Mac"] = nets[w]["Mac"]
	fmt.Println("请输入用户姓名：")
	fmt.Scanln(ss["人员"])
	fmt.Println("请输入用户部门：")
	fmt.Scanln(ss["部门"])

	return

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
			fmt.Println(k, v)
		}
		var l int
		fmt.Scanf("%d\n", l)
		cadName = nets[l]["Name"]
		us.Set(cadName, tempIp)
	} else {
		//没有发现网卡信息，程序结束
		return
	}
	if len(hdds) > 0 {
		// 读取到硬盘序列号
		for _, v := range hdds {
			infos := us.Search(v)
			if infos["Ip"] == "" {
				// 查询结果匹配则直接设置ip
				us.Set(cadName, infos["Ip"])
				return
			}
		}
		// 数据库中有硬盘信息没有网卡信息,使用网卡地址搜索
		for _, v := range nets {
			infos := us.Search(v["Mac"])
			if infos["Ip"] != "" {
				// 找到信息,直接设置IP，提示增加硬盘信息
				us.Set(cadName, infos["Ip"])

				devIt = HumanComputerInteraction(hdds, nets)
				us.Add(devIt["序列号"], infos)
				fmt.Println("已更新硬盘序列号。", devIt["序列号"])
				return
			}
		}
		// 服务器记录中没有所有设备中的硬盘信息和网卡信息，提示添加。并返回。
		devIt = HumanComputerInteraction(hdds, nets)
		us.Add("", devIt)
		return
	} else {
		// 没有从本设备中读取到硬盘序列号，使用网卡搜索信息
		for _, v := range nets {
			infos := us.Search(v["Mac"])
			if infos["Ip"] != "" {
				// 找到信息,返回信息,设置IP
				us.Set(cadName, infos["Ip"])
				return
			}
		}
		// 没有在服务器中找到该设备的网卡信息，提示添加，并返回。
		devIt = HumanComputerInteraction(hdds, nets)
		us.Add("", devIt)
		return
	}
}
