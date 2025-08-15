package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"my-blog/config"
	"my-blog/model"
)

// RegisterUser 注册用户
func RegisterUser(username, password, email, nickname string) error {
	var user model.User
	result := config.DB.Where("username = ? OR email = ?", username, email).First(&user)
	if result.RowsAffected > 0 {
		return errors.New("用户名或邮箱已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nickname: nickname,
	}

	return config.DB.Create(&newUser).Error
}

// LoginUser 登录验证
func LoginUser(username, password string) (model.User, error) {
	var user model.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return user, errors.New("用户不存在")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, errors.New("密码错误")
	}

	// 返回用户对象（可根据需要去掉密码字段）
	user.Password = ""
	return user, nil
}
