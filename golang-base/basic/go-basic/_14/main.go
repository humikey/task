package main

import "fmt"

// 最长公共前缀
func main() {

	fmt.Println(longestCommonPrefix([]string{"flower", "flow", "flight"}))
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	for i := 0; i < len(strs[0]); i++ {
		ch := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || ch != strs[j][i] {
				return strs[0][:i]
			}
		}
	}

	return strs[0]
}
