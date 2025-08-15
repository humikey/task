package controller

import (
	"github.com/gin-gonic/gin"
	"my-blog/config"
	"my-blog/model"
	"my-blog/pkg/jwt"
	"my-blog/service"
	"net/http"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 注册
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.RegisterUser(input.Username, input.Password, input.Email, input.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// Login 登录
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := service.LoginUser(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成 Token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}

// GetUserPosts 获取用户所有帖子和评论
func GetUserPosts(c *gin.Context) {
	userID := c.Param("id")

	var posts []model.Post
	if err := config.DB.
		Preload("Comments.User").
		Where("user_id = ?", userID).
		Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

//后续在控制器中获取当前登录的用户ID
//userID := c.MustGet("userID").(uint)
