package msgo

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

// 路由组
type routerGroup struct {
	name             string
	handlerFuncMap   map[string]HandleFunc
	handlerMethodMap map[string][]string
}

// 其实已经用不到这个方法了
//func (r *routerGroup) Add(name string, handleFunc HandleFunc) {
//	r.handlerFuncMap[name] = handleFunc
//}

func (r *routerGroup) Any(name string, handlerFunc HandleFunc) {
	r.handlerFuncMap[name] = handlerFunc
	r.handlerMethodMap["Any"] = append(r.handlerMethodMap["any"], name)
}
func (r *routerGroup) Get(name string, handlerFunc HandleFunc) {
	r.handlerFuncMap[name] = handlerFunc
	r.handlerMethodMap[http.MethodGet] = append(r.handlerMethodMap[http.MethodGet], name)
}
func (r *routerGroup) Post(name string, handlerFunc HandleFunc) {
	r.handlerFuncMap[name] = handlerFunc
	r.handlerMethodMap[http.MethodPost] = append(r.handlerMethodMap[http.MethodPost], name)
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
		handlerFuncMap:   make(map[string]HandleFunc),
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
				// 先看看属不属于any类型的
				routers, ok := group.handlerMethodMap["Any"]
				if ok {
					for _, routerName := range routers {
						// 比较map中存储的切片中的最后一级url和name
						if routerName == name {
							methodHandle(writer, request)
							return
						}
					}
				}
				// any中不ok 就去methodMap中遍历
				routers, ok = group.handlerMethodMap[method]
				if ok {
					for _, routerName := range routers {
						// 比较map中存储的切片中的最后一级url和name
						if routerName == name {
							methodHandle(writer, request)
							return
						}
					}
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
