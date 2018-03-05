/*
2018.02.28 by dead pig panndora

Process:
1. Get a list of hard drive addresses.
2. Get mac address list.
3. visit xlsx file search mac address, return the relevant data.
4. According to the data set ip address.
5. Prompt and modify the contents of the xlsx file according to the operation.
*/

package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Oguangpan/szsspSysManage/ynm3000"
)

// temp ip
const tmpIp = "33.66.100.255"

func main() {
	var dn *ynm3000.Computer
	dn = new(ynm3000.Computer)
	dn.GetHardNumber()
	dn.GetMacAddress()
	cowNum, _ := dn.SearchXlsx()
	hddNum := len(dn.HddIds)
	macNum := len(dn.MacAddress)
	var yon string
	var num int = 0
	/*
		分为两种情况, 依据mac地址为基础,第一种.如果有mac地址,而当前硬盘数据不相符,那么提示修改.
		第二种情况是没有查询到相关mac地址,提示是否新增数据,添加一row,输入 部门 使用者, 上传 部门
		使用者, 硬盘序列号,mac地址到服务器中, 后续ip分配工作直接编辑xlsx文件完成.
	*/
	var update bool
	if cowNum == 0 {
		// 没有查询到mac地址

		stdin := bufio.NewReader(os.Stdin)
		fmt.Println("您的电脑信息并没有在服务器记录中，是否添加新的设备信息。 yes or no ?>")
		fmt.Scanln(&yon)
		if yon != "yes" {
			return
		}
		fmt.Println("请输入该设备使用者姓名>")
		fmt.Scanln(&dn.User)
		fmt.Println("请输入该设备的所属部门>")
		fmt.Scanln(&dn.Department)
		if hddNum > 1 {
			fmt.Println("您的电脑拥有多个硬盘,请选择其中一个作为记录硬盘>")
			for i, j := range dn.HddIds {
				fmt.Println(i, ":", j)
			}
			fmt.Fscan(stdin, &num)
			dn.HddId = dn.HddIds[num]
		} else if hddNum > 1 {
			dn.HddId = dn.HddIds[0]
		} else {
			fmt.Println("警告:没有发现该设备上的硬盘物理序列号")
		}
		if macNum > 1 {
			fmt.Println("您的电脑中有多个网卡,请选择其中一个作为接入内网专用>")
			for i, j := range dn.MacAddress {
				fmt.Println(i, ": ", j["Name"], j["Mac"])
			}
			fmt.Fscan(stdin, &num)
			dn.MacAddres = dn.MacAddress[num]["Mac"]

		} else if macNum > 0 {
			dn.MacAddres = dn.MacAddress[0]["Mac"]
		} else {
			fmt.Println("警告:没有在该设备上发现可用的网卡")
		}
		fmt.Println("该机器信息如下: \n用户:",
			dn.User, "\n部门:",
			dn.Department, "\n硬盘物理序列号:",
			dn.HddId, "\n网卡IP地址:",
			dn.MacAddres, "\n是否直接提交到服务器中? yes or no >")
		fmt.Scanln(&yon)
		if yon == "yes" {
			dn.AmXlsx(false)
		} else {
			fmt.Println("您选择了no，结束提交，服务器记录中仍然没有该设备的记录。")
		}
		return

	} else {
		// 查询到了mac地址, 判断硬盘序列号是否与服务器
		if hddNum > 0 {
			update = true
			for _, j := range dn.HddIds {
				if j == dn.HddId {
					update = false
				}
			}
			if update {
				fmt.Println("检查到当前电脑硬盘序列号与服务器记录不相符，是否修正？ yes or no>")
				var q string
				fmt.Scanln(&q)
				if q == "yes" {
					switch hddNum {
					case hddNum == 1:
						dn.HddId = dn.HddIds[0]
					case hddNum > 1:
						fmt.Println("您的设备中有多个硬盘，请选择需要绑定的那个。>")
						for k, v := range dn.HddIds {
							fmt.Println(k, ":", v)
						}
						fmt.Fscan(stdin, &num)
						dn.HddId = dn.HddIds[num]
					}
				} else {
					fmt.Println("您选择了no，将不对硬盘序列号记录做修改。")
				}

			} else {
				fmt.Println("警告:该设备没有发现可用的硬盘序列号")
			}

			dn.AmXlsx(true)
		}
		//fmt.Println(cowNum, cardName, dn.Ip)

	}
}
