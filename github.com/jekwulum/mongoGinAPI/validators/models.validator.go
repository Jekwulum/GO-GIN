package validators

import (
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/jekwulum/mongoGinAPI/models"
)

type CreateBlogInput struct {
	Title string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateBlogInput struct {
	Title string `json:"title"`
	Content string `json:"content"`
}

func CreateBlogModel (c *gin.Context) (*CreateBlogInput, error) {
	var input CreateBlogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

func UpdateBlogModel (c *gin.Context) (*models.Blog, *UpdateBlogInput, error) {
	var blog models.Blog
	var input UpdateBlogInput

	existsErr := models.DB.Where("id = ?", c.Param("id")).First(&blog).Error
	if existsErr != nil {
		return nil, nil, existsErr
	}

	if bindErr := c.ShouldBindJSON(&input); bindErr != nil {
		return nil, nil, bindErr
	}

	return &blog, &input, nil
}