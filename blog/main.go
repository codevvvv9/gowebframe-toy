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
	engine.Add("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s 欢迎来到手写web框架", "wushao")
	})

	engine.Run()
}
