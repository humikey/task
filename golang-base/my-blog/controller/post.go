package controller

import (
	"github.com/gin-gonic/gin"
	"my-blog/config"
	"my-blog/model"
	"net/http"
	"strconv"
)

type CreatePostInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// CreatePost 发布博客（需要认证）
func CreatePost(c *gin.Context) {
	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前登录用户ID（从 JWT 中间件注入）
	userID := c.MustGet("userID").(uint)

	post := model.Post{
		Title:   input.Title,
		Content: input.Content,
		UserID:  userID,
	}

	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发布失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功", "post": post})
}

// UpdatePost 实现文章的更新功能，只有文章的作者才能更新自己的文章。
func UpdatePost(c *gin.Context) {
	postID, err := c.Params.Get("id")
	if !err {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未获取到博客ID"})
		return
	}
	// 获取当前登录用户ID（从 JWT 中间件注入）
	userID := c.MustGet("userID").(uint)
	var post model.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查是否是作者
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文章"})
		return
	}

	var input CreatePostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据不合法"})
		return
	}
	post.Title = input.Title
	post.Content = input.Content
	if err := config.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "post": post})
}

// DeletePost 实现文章的删除功能，只有文章的作者才能删除自己的文章。
func DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
	}
	var post model.Post
	err = config.DB.First(&post, postID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}
	userID := c.MustGet("userID").(uint)
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
	}
	err = config.DB.Delete(&post).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除文章成功"})
}

func GetMostCommentedPost(c *gin.Context) {
	type PostWithCommentCount struct {
		ID           uint
		Title        string
		Content      string
		CommentCount int
	}

	var result PostWithCommentCount
	err := config.DB.
		Table("posts").
		Select("posts.id, posts.title, posts.content, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}
