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
	groupUser.Use(func(next msgo.HandleFunc) msgo.HandleFunc {
		return func(ctx *msgo.Context) {
			fmt.Println("我是前面的中间件")
			next(ctx)
			fmt.Println("我是后面的中间件")
		}
	})
	groupUser.Get("/hi/*/get", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s /hi/*/get 欢迎来到手写web框架", "wushao")
	})
	groupUser.Post("/up", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s post 欢迎来到手写web框架", "wushao")
	})
	groupUser.Put("/up", func(ctx *msgo.Context) {
		fmt.Println("我是next")
		fmt.Fprintf(ctx.W, "%s get 欢迎来到手写web框架", "wushao")
	})
	groupGoods := engine.Group("goods")
	groupGoods.Any("goodList", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s 货物清单", "wushao")
	})
	groupGoods.Get("/goodList2", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s 2货物清单", "wushao")
	})
	groupGoods.Get("/goodList2/:id", func(ctx *msgo.Context) {
		fmt.Fprintf(ctx.W, "%s :id is ok", "wushao")
	})
	engine.Run()
}
