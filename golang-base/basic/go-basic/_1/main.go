package main

import "fmt"

// 两数之和
func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum(nums, target)
	fmt.Println("结果索引:", result) // 输出: [0 1]
}

func twoSum(nums []int, target int) []int {
	var sumMap = make(map[int]int)

	for index, value := range nums {
		if i, ok := sumMap[target-value]; ok {
			return []int{i, index}
		}
		sumMap[value] = index
	}

	return nil
}
