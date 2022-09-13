package main

import (
	"fmt"
	"github.com/codevvvv9/msgo"
	"net/http"
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
	groupUser.Add("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s 欢迎来到手写web框架", "wushao")
	})
	groupGoods := engine.Group("goods")
	groupGoods.Add("goodList", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s 货物清单", "wushao")
	})
	engine.Run()
}
