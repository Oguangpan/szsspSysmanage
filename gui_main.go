package main

import (
	"log"

	"github.com/sciter-sdk/go-sciter"
	"github.com/sciter-sdk/go-sciter/window"
)

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
		&sciter.Rect{Left: 0, Top: 0, Right: 335, Bottom: 405})
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
