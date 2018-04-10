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
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

type xitong interface {
	Shezhiwangluodizhi(dizhi string) error
	Huoqubenjixinxi() xt
	Huoqushujukuxinxi() xt
	Shangchuanxinxi(xinxi xt)
}

type xt struct {
	yonghuming      string
	bumen           string
	yingpanxuliehao []string
	wangkamac       []string
	ipdizhi         string
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
	//加载文件
	w.LoadFile("demo1.html")
	//设置标题
	w.SetTitle("三洲特管信息化台账录入系统 v1.0")
	//显示窗口
	w.Show()
	//运行窗口，进入消息循环
	w.Run()
}
