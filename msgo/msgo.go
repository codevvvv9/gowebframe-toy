package msgo

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

// 路由组
type routerGroup struct {
	name          string
	handleFuncMap map[string]HandleFunc
}

func (r *routerGroup) Add(name string, handleFunc HandleFunc) {
	r.handleFuncMap[name] = handleFunc
}

// 3、user /get/list user组下面才是url
//路由表 由路由组组成
type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		//handleFuncMap: map[string]HandleFunc{},
		// 和上面的写法一个效果
		handleFuncMap: make(map[string]HandleFunc),
		name:          name,
	}

	r.routerGroups = append(r.routerGroups, group)
	return group
}

type Engine struct {
	*router //这样写的原因是啥呢
}

func New() *Engine {
	return &Engine{
		&router{},
	}
}

func (e *Engine) Run() {
	for _, group := range e.routerGroups {
		for key, value := range group.handleFuncMap {
			groupNameHasSlash := strings.HasPrefix(group.name, "/")
			routerKeyHasSlash := strings.HasPrefix(key, "/")
			var groupName string
			var handleFuncName string
			if groupNameHasSlash {
				groupName = group.name
			} else {
				groupName = "/" + group.name
			}
			if routerKeyHasSlash {
				handleFuncName = groupName + key
			} else {
				handleFuncName = groupName + "/" + key
			}

			http.HandleFunc(handleFuncName, value)
		}

	}

	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}

}
