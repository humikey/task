package main

import (
	"fmt"
	"sort"
)

// 合并区间
func main() {
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	result := merge(intervals)
	fmt.Println(result) // [[1 6] [8 10] [15 18]]
}

func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return intervals
	}

	// 按区间起点排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	merged := [][]int{intervals[0]}

	for i := 1; i < len(intervals); i++ {
		last := merged[len(merged)-1]
		current := intervals[i]

		if current[0] > last[1] {
			// 无重叠，直接添加
			merged = append(merged, current)
		} else {
			// 有重叠，更新终点为最大值
			if current[1] > last[1] {
				last[1] = current[1]
				merged[len(merged)-1] = last
			}
		}
	}

	return merged
}
