package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {

	lock1()
	lock2()
}

//题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//考察点 ： sync.Mutex 的使用、并发数据安全。

var lock sync.Mutex
var wg sync.WaitGroup

func lock1() {

	counter := 0
	wg.Add(10)
	for i := 1; i <= 10; i++ {
		go addNum(&counter)
	}
	wg.Wait()

	fmt.Println(counter)
}

func addNum(counter *int) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		lock.Lock()
		*counter++
		lock.Unlock()
	}
}

//题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
//考察点 ：原子操作、并发数据安全。

func lock2() {
	var counter int64 // 使用 int64 类型，确保和 atomic 操作类型一致
	var wg sync.WaitGroup

	// 启动 10 个协程
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1) // 原子加 1
			}
		}()
	}

	wg.Wait()
	fmt.Println("最终计数值：", counter)
}
