package main

import "fmt"

// 66. 加一
func main() {

	fmt.Println(plusOne([]int{1, 2, 3, 9}))

}

// 位数太大会溢出
func plusOneWrong(digits []int) []int {

	sum := 0
	for _, value := range digits {
		sum = sum*10 + value
	}
	sum += 1

	var res []int
	for sum > 0 {
		res = append([]int{sum % 10}, res...)
		sum /= 10
	}

	return res
}

func plusOne(digits []int) []int {
	n := len(digits)

	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	// 所有位都是9，例如[9,9,9]，结果应为[1,0,0,0]
	result := make([]int, n+1)
	result[0] = 1
	return result
}
