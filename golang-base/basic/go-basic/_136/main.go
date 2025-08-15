package main

import "fmt"

// 136. 只出现一次的数字
func main() {
	nums := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	fmt.Println(singleNumber(nums))
}

func singleNumber(nums []int) int {
	res := 0
	for _, v := range nums {
		res ^= v
	}
	return res
}

func singleNumber2(nums []int) int {
	var mapSet = make(map[int]int)
	for _, value := range nums {
		if _, ok := mapSet[value]; ok {
			mapSet[value]++
		} else {
			mapSet[value] = 1
		}
	}

	for i, value := range mapSet {
		if value == 1 {
			return i
		}
	}
	return 0
}
