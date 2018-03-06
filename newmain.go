// 通过查询当前设备部分信息到服务器中获取更多的信息
// 设置当前设备的网络IP地址

package main

import (
	"fmt"
)

const xlsxFile string = "//33.66.96.14/public/2018Taizhang.xlsx"
const tempIp string = "33.66.100.255"

type ultinvsupgader interface {
	Get() ([]string, []map[string]string)
	Search(string) []string
	Add(map[string]string)
	Set(string)
}
type uisg struct{}

// 返回: 一组硬盘序列号和一组网卡信息
func (p *uisg) Get() (s []string, m []map[string]string) {
	//TODO
}

// 参数：序列号或MAC号 返回：记录中符合条件的所有信息
func (p *uisg) Search(s string) (ss map[string]string) {
	//TODO
}

// 参数：当前设备所有信息
func (p *uisg) Add(ss map[string]string) {
	//TODO
}

// 参数：ip
func (p *uisg) Set(ip string) {
	//TODO
}

func main() {
	us := new(ultinvsupgader)
	hdds, nets := us.Get()
	if len(nets) > 0 {
		//设置临时地址
		us.Set(tempIp)
	} else {
		//没有发现网卡信息，程序结束
		return
	}
	if len(hdds) < 1 {
		// 读取到硬盘序列号
		for _, v := range hdds {
			infos := us.Search(v)
			if infos["Ip"] == "" {
				// 查询结果匹配则直接设置ip
				return
			} else {
				// 否则继续查询
				continue
			}
			// 如果所有匹配都没有成功，使用网卡mac进行匹配

		}

	} else {
		// 没有读取到硬盘序列号

	}
}
