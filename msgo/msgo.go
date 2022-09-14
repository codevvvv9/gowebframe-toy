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

// 路由组
type routerGroup struct {
	name string
	// 同一路径的不同请求方式
	handlerFuncMap map[string]map[string]HandleFunc
	//支持不同的请求方式  {"post": ["/hi", "/hello"]}
	handlerMethodMap map[string][]string
}

// 其实已经用不到这个方法了
//func (r *routerGroup) Add(name string, handleFunc HandleFunc) {
//	r.handlerFuncMap[name] = handleFunc
//}

func (r *routerGroup) handle(name string, method string, handleFunc HandleFunc) {
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
	}

	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	*router //这样写的原因是啥呢
}

func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// 先拿到请求的方法类型 GET POST啥的
	method := request.Method
	groups := e.router.routerGroups
	// 根据url进行匹配
	for _, group := range groups {
		for name, methodHandle := range group.handlerFuncMap {
			groupNameHasSlash := strings.HasPrefix(group.name, "/")
			routerKeyHasSlash := strings.HasPrefix(name, "/")
			var groupName string
			var requestUrl string
			if groupNameHasSlash {
				groupName = group.name
			} else {
				groupName = "/" + group.name
			}
			if routerKeyHasSlash {
				requestUrl = groupName + name
			} else {
				requestUrl = groupName + "/" + name
			}

			//http.HandleFunc(requestUrl, methodHandle)
			// 比较请求url和拼接的url是否相同
			if request.RequestURI == requestUrl {
				//构造上下文
				ctx := &Context{
					W: writer,
					R: request,
				}
				// 先看看属不属于any类型的
				handle, ok := methodHandle[Any]
				if ok {
					handle(ctx)
					return
				}
				handle, ok = methodHandle[method]
				if ok {
					handle(ctx)
					return
				}
				// 具体method中不ok 就报错了
				writer.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(writer, "%s %s not allowed \n", request.RequestURI, method)
				return
			}

		}

	}
	// url不匹配
	writer.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(writer, "%s %s not found \n", request.RequestURI, method)
	return
	//TODO implement me
	panic("implement me")
}

func New() *Engine {
	return &Engine{
		&router{},
	}
}

func (e *Engine) Run() {
	//for _, group := range e.routerGroups {
	//	for key, value := range group.handlerFuncMap {
	//		groupNameHasSlash := strings.HasPrefix(group.name, "/")
	//		routerKeyHasSlash := strings.HasPrefix(key, "/")
	//		var groupName string
	//		var handleFuncName string
	//		if groupNameHasSlash {
	//			groupName = group.name
	//		} else {
	//			groupName = "/" + group.name
	//		}
	//		if routerKeyHasSlash {
	//			handleFuncName = groupName + key
	//		} else {
	//			handleFuncName = groupName + "/" + key
	//		}
	//
	//		http.HandleFunc(handleFuncName, value)
	//	}
	//
	//}

	// 3、支持不同的方法，全拦截
	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}

}
