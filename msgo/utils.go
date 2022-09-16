package msgo

import "strings"

func SubStringLast(str string, subStr string) string {
	//先查找有没有
	index := strings.Index(str, subStr)
	if index == -1 {
		return ""
	}
	len := len(subStr)
	return str[(index + len):]
}
