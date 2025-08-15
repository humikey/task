package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
//要求 ：
//编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
//编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:root@tcp(172.25.129.74:3306)/gorm_learn?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	//questionOne(err, db)
	questionTwo(err, db)
}

func questionOne(err error, db *sqlx.DB) {
	var employees []Employee
	err = db.Select(&employees, "select id,name,department,salary from employees where department = ?", "技术部")
	if err != nil {
		panic(err)
	}

	for _, u := range employees {
		fmt.Println(u.ID, u.Name, u.Department, u.Salary)
	}

	var employee []Employee
	err = db.Select(&employee, "select id,name,department,salary from employees order by salary desc limit 1")
	if err != nil {
		panic(err)
	}
	fmt.Println(employee[0].Name)
}

//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//要求 ：
//定义一个 Book 结构体，包含与 books 表对应的字段。
//编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

type Book struct {
	ID     int
	Title  string
	Author string
	Price  float64
}

func questionTwo(err error, db *sqlx.DB) {
	var books []Book
	err = db.Select(&books, "select id, title, author, price from books where price > ?", 50)
	if err != nil {
		panic(err)
	}

	for _, book := range books {
		fmt.Println(book.ID, book.Price, book.Title)
	}
}
