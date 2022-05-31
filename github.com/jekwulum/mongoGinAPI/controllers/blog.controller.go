package controllers

import (
	// "fmt"
	"net/http"
	// "strconv"

	"github.com/gin-gonic/gin"
	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/validators"
)


func GetBlogs (c *gin.Context) {
	var blogs []models.Blog
	models.DB.Find(&blogs)
	c.IndentedJSON(http.StatusOK, gin.H{"data": blogs})
}

func GetBlog (c *gin.Context) {
	var blog models.Blog

	err := models.DB.Where("id = ?", c.Param("id")).First(&blog).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"data": blog})
}

func CreateBlog (c *gin.Context) {
	input, err := validators.CreateBlogModel(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	blog := models.Blog{Title: input.Title, Content: input.Content}
	createErr := models.DB.Create(&blog).Error
	if createErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": createErr.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"data": blog})
}

func UpdateBlog (c *gin.Context) {
	blog, input, err := validators.UpdateBlogModel(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateErr := models.DB.Model(&blog).Updates(&input).Error; updateErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": updateErr.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"data": blog})
}

func DeleteBlog (c *gin.Context) {
	var blog models.Blog
	err := models.DB.Where("id = ?", c.Param("id")).First(&blog).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if deleteErr := models.DB.Delete(&blog); deleteErr != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": deleteErr.Error()})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func RegisterBlogRoutes(rg *gin.RouterGroup) {
	blogRoute := rg.Group("/blog")
	{
		blogRoute.GET("/", GetBlogs)
		blogRoute.GET("/:id", GetBlog)
		blogRoute.POST("/", CreateBlog)
		blogRoute.PATCH("/:id", UpdateBlog)
		blogRoute.DELETE("/:id", DeleteBlog)
	}
}