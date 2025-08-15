package main

import (
	"fmt"
	"math"
)

func main() {

	runShape()

	runPerson()
}

// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。
func runShape() {
	r := &Rectangle{Width: 3, Height: 4}
	c := &Circle{Radius: 5}

	// 使用接口类型进行统一处理
	shapes := []Shape{r, c}

	for _, shape := range shapes {
		fmt.Printf("类型：%T\n", shape)
		fmt.Printf("面积：%.2f\n", shape.Area())
		fmt.Printf("周长：%.2f\n", shape.Perimeter())
		fmt.Println("--------------")
	}

}

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (this *Rectangle) Area() float64 {
	return this.Width * this.Height
}

func (this *Rectangle) Perimeter() float64 {
	return 2 * (this.Width + this.Height)
}

type Circle struct {
	Radius float64
}

func (this *Circle) Area() float64 {
	return math.Pi * this.Radius * this.Radius
}

func (this *Circle) Perimeter() float64 {
	return 2 * math.Pi * this.Radius
}

//题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
//考察点 ：组合的使用、方法接收者。

type Person struct {
	Age  int
	Name string
}

type Employee struct {
	Person
	EmployeeID int
}

func (this *Employee) PrintInfo() {
	fmt.Printf("name:%v, age:%d, employee:%d", this.Name, this.Age, this.EmployeeID)
}

func runPerson() {
	emp := &Employee{EmployeeID: 112, Person: Person{Age: 1, Name: "jhshe"}}
	emp.PrintInfo()
}
