package main

import (
	"fmt"
	"sync"
	"time"
)

// 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 考察点 ： go 关键字的使用、协程的并发执行。
func main() {
	odd, even := make(chan bool), make(chan bool)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10; i = i + 2 {
			<-odd
			fmt.Println(i)
			even <- true
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= 10; i = i + 2 {
			<-even
			fmt.Println(i)
			if i == 10 {
				break
			}
			odd <- true
		}
	}()

	odd <- true
	wg.Wait()
	//
	//run2()

	runPrint()
}

//题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
//考察点 ：协程原理、并发任务调度。

type Task func()

type TaskWithTime struct {
	Index int
	Time  time.Duration
}

func scheduleTask(tasks []Task) []TaskWithTime {
	var wgInner sync.WaitGroup
	result := make([]TaskWithTime, len(tasks))
	for index, task := range tasks {
		wgInner.Add(1)
		go func(i int, run Task) {
			defer wgInner.Done()
			start := time.Now()
			run()
			elapsed := time.Since(start)
			result[i] = TaskWithTime{Index: i, Time: elapsed}
		}(index, task)
	}
	wgInner.Wait()
	return result
}

func run2() {
	tasks := []Task{
		func() {
			time.Sleep(time.Millisecond * 10)
		},
		func() {
			time.Sleep(time.Millisecond * 500)
		},
		func() {
			time.Sleep(time.Millisecond * 50)
		},
	}

	result := scheduleTask(tasks)

	for i := range result {
		fmt.Printf("任务 %d 执行耗时: %v \n", i, result[i].Time)
	}
}

func runPrint() {
	odd, even := make(chan bool), make(chan bool)

	wg := sync.WaitGroup{}

	go func() {
		num := 1
		for {
			select {
			case <-odd:
				fmt.Println(num)
				num = num + 2
				even <- true
			}
		}
	}()

	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		num := 2
		for {
			if num > 100 {
				wg.Done()
				return
			}
			select {
			case <-even:
				fmt.Println(num)
				num = num + 2
				odd <- true
			}
		}

	}(&wg)

	odd <- true
	wg.Wait()
}
