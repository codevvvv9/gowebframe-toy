package msgo

import (
	"fmt"
	"testing"
)

func TestTreeNode(t *testing.T) {
	root := &treeNode{
		name:       "/",
		children:   make([]*treeNode, 0),
		routerName: "",
		isEnd:      false,
	}

	root.Put("/user/get/:id")
	root.Put("/user/af/hello")
	root.Put("/user/create/aaa")

	node := root.Get("/user/get/2")
	fmt.Println(node)
	node = root.Get("/user/af/hello")
	fmt.Println(node)
	node = root.Get("/user/create/aaa")
	fmt.Println(node)

}
