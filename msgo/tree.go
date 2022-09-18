package msgo

import "strings"

// 字典树或者前缀树结构，根节点root是空字符串，三个基本性质
//1. 根节点不包含字符，除根节点外每一个节点都只包含一个字符。
//2. 从根节点到某一节点，路径上经过的字符连接起来，为该节点对应的字符串。
//3. 每个节点的所有子节点包含的字符都不相同。
type treeNode struct {
	name       string //节点名称 /user
	children   []*treeNode
	routerName string
	isEnd      bool
}

// put path: /user/get/:id

func (t *treeNode) Put(path string) {
	root := t
	//根据 /来分割 [ user get :id]
	strs := strings.Split(path, "/")
	for index, name := range strs {
		// 第一个节点是个空字符串 "" 直接跳过
		if index == 0 {
			continue
		}
		// 第二个节点 user要看看它的子节点
		children := t.children
		// 标志是否匹配到了
		isMatch := false
		for _, node := range children {
			if node.name == name {
				isMatch = true
				t = node
				break
			}
		}
		// 没匹配上就把节点就组建回来
		if !isMatch {
			// put时判断是不是尾节点
			isEnd := false
			if index == len(strs)-1 {
				isEnd = true
			}
			node := &treeNode{
				name:     name,
				children: make([]*treeNode, 0),
				isEnd:    isEnd,
			}
			children = append(children, node)
			t.children = children
			t = node
		}
	}

	t = root
}

// Get path: /user/get/1
func (t *treeNode) Get(path string) *treeNode {
	strs := strings.Split(path, "/")
	routerName := ""
	for index, name := range strs {
		if index == 0 {
			continue
		}
		children := t.children
		isMatch := false
		for _, node := range children {
			if node.name == name ||
				node.name == "*" ||
				strings.Contains(node.name, ":") {
				isMatch = true
				routerName += "/" + node.name
				node.routerName = routerName
				t = node
				if index == len(strs)-1 {
					return node
				}
				break
			}
		}

		if !isMatch {
			for _, node := range children {
				// /user/**
				// /user/get/userInfo
				// /user/aa/bb
				if node.name == "**" {
					routerName += "/" + node.name
					node.routerName = routerName
					return node
				}
			}
		}
	}

	return nil
}
