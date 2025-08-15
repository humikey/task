package controller

import (
	"github.com/gin-gonic/gin"
	"my-blog/config"
	"my-blog/model"
	"net/http"
	"strconv"
)

type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
}

type CommentResponse struct {
	ID      uint         `json:"id"`
	Content string       `json:"content"`
	User    UserResponse `json:"user"`
}
type PostResponse struct {
	ID       uint              `json:"id"`
	Title    string            `json:"title"`
	Content  string            `json:"content"`
	User     UserResponse      `json:"user"`
	Comments []CommentResponse `json:"comments"`
}

func CreateComment(c *gin.Context) {
	// 获取 Post ID
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章 ID"})
		return
	}

	// 验证文章是否存在
	var post model.Post
	if err := config.DB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 绑定请求体
	var input CreateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取当前用户ID
	userID := c.MustGet("userID").(uint)

	// 创建评论
	comment := model.Comment{
		Content: input.Content,
		PostID:  uint(postID),
		UserID:  userID,
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发表评论失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论成功", "comment": comment})
}

// GetPostDetail 获取帖子详情
func GetPostDetail(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var post model.Post
	// Preload 加载关联数据：作者、评论和评论作者
	err = config.DB.
		Preload("User").
		Preload("Comments").
		Preload("Comments.User").
		First(&post, postID).Error

	// 将 Post 转换为 PostResponse 优化冗余数据展示信息
	response := PostResponse{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		User: UserResponse{
			ID:       post.User.ID,
			Username: post.User.Username,
			Email:    post.User.Email,
			Nickname: post.User.Nickname,
		},
	}

	for _, c := range post.Comments {
		response.Comments = append(response.Comments, CommentResponse{
			ID:      c.ID,
			Content: c.Content,
			User: UserResponse{
				ID:       c.User.ID,
				Username: c.User.Username,
				Email:    c.User.Email,
				Nickname: c.User.Nickname,
			},
		})
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, response)
}
