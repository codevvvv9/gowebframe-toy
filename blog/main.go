package main

import (
	"fmt"
	"github.com/codevvvv9/msgo"
)

func main() {
	//1. 使用远程实现web接口
	//http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Fprintf(writer, "%s 1欢迎来到手写web框架", "wushao")
	//})
	//http.ListenAndServe(":8111", nil)

	//2. 使用手写框架实现
	engine := msgo.New()
	groupUser := engine.Group("user")
	groupUser.Get("/hi", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s 欢迎来到手写web框架", "wushao")
	})
	groupUser.Post("/up", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s post 欢迎来到手写web框架", "wushao")
	})
	groupUser.Get("/up", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s get 欢迎来到手写web框架", "wushao")
	})
	groupGoods := engine.Group("goods")
	groupGoods.Any("goodList", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s 货物清单", "wushao")
	})
	groupGoods.Get("goodList2", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s 2货物清单", "wushao")
	})
	engine.Run()
}
