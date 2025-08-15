package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Email    string `gorm:"unique;not null" json:"email"`
	Nickname string `json:"nickname"`
	Posts    []Post
	Comments []Comment
}

type Post struct {
	gorm.Model
	Title    string `gorm:"not null"`
	Content  string `gorm:"not null"`
	UserID   uint
	User     User
	Comments []Comment
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint
	User    User
	PostID  uint
	Post    Post
}
