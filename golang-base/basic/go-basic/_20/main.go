package main

import (
	"fmt"
	"strings"
)

// 20. 有效的括号
func main() {
	s := "()(()){([])}"
	fmt.Println(isValid(s))
}

func isValid(s string) bool {
	var groupMap = make(map[byte]byte)
	groupMap['('] = ')'
	groupMap['['] = ']'
	groupMap['{'] = '}'

	var stack []byte
	for _, value := range []byte(s) {
		if _, ok := groupMap[value]; ok {
			stack = append(stack, value)
		} else {
			if len(stack) == 0 || groupMap[stack[len(stack)-1]] != value {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func isValidA(s string) bool {
	for strings.Contains(s, "{}") || strings.Contains(s, "[]") || strings.Contains(s, "()") {
		s = strings.ReplaceAll(s, "{}", "")
		s = strings.ReplaceAll(s, "[]", "")
		s = strings.ReplaceAll(s, "()", "")
	}
	return s == ""
}
