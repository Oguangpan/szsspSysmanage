/*
表格数据读写包
用途: 获取远程服务器上的表格文件中的数据.并支持修改数据
使用方法:
  查询数据 : xlsx.inquire(mac)
  对比数据 : xlsx.compared([]string)
  修改数据 : xlsx.modify([]string)
*/
package xlsx

import (
	"errors"

	"github.com/tealeg/xlsx"
)

const xlsxFilePath = "//33.66.96.14/Public/2017盘点/2017年公司办信息设备台账(2017-12-15最后编辑).xls"

// xlsx表单接口,三个方法 查询 对比 修改
type xlsx interface {
	inquire(mac string)
	compared([]string) bool
	modify()
}

// 设备相关信息
type compData struct {
	userName string
	divSio   string
	os       string
	id       string
	mac      string
	ip       string
}
