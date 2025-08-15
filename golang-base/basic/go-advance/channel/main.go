package main

import (
	"fmt"
	"sync"
)

func main() {
	channel1()
}

// 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 考察点 ：通道的基本使用、协程间通信。
var wg sync.WaitGroup

func channel1() {
	chanInNum := make(chan int)

	go handleChanIn(chanInNum)

	wg.Add(1)
	go handleChanOut(chanInNum)
	wg.Wait()
}

func handleChanOut(num chan int) {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		val := <-num
		fmt.Println(val)
	}
}

func handleChanIn(num chan int) {
	for i := 1; i <= 10; i++ {
		num <- i
	}
	close(num)
}

//题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
//考察点 ：通道的缓冲机制。

func channel2() {
	ch := make(chan int, 10) // 创建一个缓冲区大小为10的通道
	var wg sync.WaitGroup
	wg.Add(2)

	// 生产者协程：发送 1~100 到通道
	go func() {
		defer wg.Done()
		for i := 1; i <= 100; i++ {
			ch <- i
			// fmt.Println("发送：", i)
		}
		close(ch) // 发送完成后关闭通道
	}()

	// 消费者协程：从通道接收并打印
	go func() {
		defer wg.Done()
		for num := range ch {
			fmt.Println("接收：", num)
		}
	}()

	wg.Wait()
}
