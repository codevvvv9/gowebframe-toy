package msgo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const Any = "ANY"

// HandleFunc 使用上下文结构体改造
type HandleFunc func(ctx *Context)

// MiddlewareFunc 中间件函数： 接受一个函数，处理完再把这个函数返回
type MiddlewareFunc func(handlerFunc HandleFunc) HandleFunc

// 路由组 一定要注意初始化
type routerGroup struct {
	name string
	// 同一路径的不同请求方式
	handlerFuncMap map[string]map[string]HandleFunc
	//支持不同的请求方式  {"post": ["/hi", "/hello"]}
	handlerMethodMap map[string][]string
	// 前缀树
	treeNode *treeNode
	// 组级别中间件map
	middlewareMap map[string]map[string][]MiddlewareFunc
	// 组级别的多个中间件函数
	middlewares []MiddlewareFunc
}

// Use 触发中间件函数
func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(method string, handleMap map[string]HandleFunc, ctx *Context) {
	h := handleMap[method]
	middlewares := r.middlewares
	if middlewares != nil {
		for _, middleware := range middlewares {
			h = middleware(h)
		}
	}
	h(ctx)
}
func (r *routerGroup) handle(name string, method string, handleFunc HandleFunc) {
	//兼容路径名没有 / 的情况
	nameHasPrefix := strings.HasPrefix(name, "/")
	if !nameHasPrefix {
		name = "/" + name
	}
	_, ok := r.handlerFuncMap[name]
	if !ok {
		r.handlerFuncMap[name] = make(map[string]HandleFunc)
	}
	_, ok = r.handlerFuncMap[name][method]
	if ok {
		panic("有重复路由")
	}
	r.handlerFuncMap[name][method] = handleFunc
	r.handlerMethodMap[method] = append(r.handlerMethodMap[method], name)

	r.treeNode.Put(name)
}
func (r *routerGroup) Any(name string, handlerFunc HandleFunc) {
	r.handle(name, Any, handlerFunc)
}
func (r *routerGroup) Get(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodGet, handlerFunc)
}
func (r *routerGroup) Post(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodPost, handlerFunc)
}
func (r *routerGroup) Put(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodPut, handlerFunc)
}
func (r *routerGroup) Delete(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodDelete, handlerFunc)
}
func (r *routerGroup) Patch(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodPatch, handlerFunc)
}
func (r *routerGroup) Options(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodOptions, handlerFunc)
}
func (r *routerGroup) Head(name string, handlerFunc HandleFunc) {
	r.handle(name, http.MethodHead, handlerFunc)
}

// 3、user /get/list user组下面才是url
//路由表 由路由组组成
type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		//handlerFuncMap: map[string]HandleFunc{},
		// 和上面的写法一个效果
		handlerFuncMap:   make(map[string]map[string]HandleFunc),
		name:             name,
		handlerMethodMap: make(map[string][]string),
		treeNode: &treeNode{
			name:       "/",
			children:   make([]*treeNode, 0),
			routerName: "",
			isEnd:      false,
		},
	}

	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	*router //这样写的原因是啥呢
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	e.requestHTTP(writer, request)
}

func New() *Engine {
	return &Engine{
		&router{},
	}
}

func (e *Engine) Run() {
	// 3、支持不同的方法，全拦截
	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// ServerHTTP真正的执行逻辑代码
func (e *Engine) requestHTTP(writer http.ResponseWriter, request *http.Request) {
	// 先拿到请求的方法类型 GET POST啥的
	method := request.Method
	groups := e.router.routerGroups
	// 根据url进行匹配
	for _, group := range groups {
		routerName := SubStringLast(request.RequestURI, "/"+group.name)
		// routerName 类似于 get/1
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			//路由匹配上了
			ctx := &Context{
				W: writer,
				R: request,
			}
			// 先查看ANY类型
			_, ok := group.handlerFuncMap[node.routerName][Any]
			if ok {
				//handle(ctx)
				group.methodHandle(Any, group.handlerFuncMap[node.routerName], ctx)
				return
			}
			// 再查看具体方法
			_, ok = group.handlerFuncMap[node.routerName][method]
			if ok {
				//handle(ctx)
				group.methodHandle(method, group.handlerFuncMap[node.routerName], ctx)
				return
			}
			// 具体method中不ok 就报错了
			writer.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(writer, "%s %s not allowed \n", request.RequestURI, method)
			return
		}
	}
	// url不匹配
	writer.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(writer, "%s %s not found \n", request.RequestURI, method)
}
