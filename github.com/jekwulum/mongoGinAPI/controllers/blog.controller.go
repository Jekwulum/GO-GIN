package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/services"
	"github.com/jinzhu/gorm"
)

type BlogController struct {
	DB *gorm.DB
}

func (bc *BlogController) GetBlogs(ctx *gin.Context) {
	blogs, err := services.GetBlogs(bc.DB)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, blogs)
}

func (bc *BlogController) GetBlog(ctx *gin.Context) {
	str_id := ctx.Params.ByName("id")
	id, _ := strconv.ParseUint(str_id, 10, 64)
	blog, exists, err := services.GetBlog(id, bc.DB)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !exists {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "blog not found"})
		return
	}
	ctx.IndentedJSON(http.StatusOK, blog)
}

func (bc *BlogController) CreateBlog(ctx *gin.Context) {
	blog := models.Blog{}
	err := ctx.BindJSON(&blog)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := bc.DB.Create(&blog).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"data": blog})
}

func (bc *BlogController) UpdateBlog(ctx *gin.Context) {
	str_id := ctx.Params.ByName("id")
	id, _ := strconv.ParseUint(str_id, 10, 64)
	_, exists, err := services.GetBlog(id, bc.DB)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not exists"})
		return
	}

	updatedBlog := models.Blog{}
	err = ctx.BindJSON(&updatedBlog)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := services.UpdateBlog(bc.DB, &updatedBlog); err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	bc.GetBlog(ctx)
}

func (bc *BlogController) DeleteBlog(ctx *gin.Context) {
	str_id := ctx.Params.ByName("id")
	id, _ := strconv.ParseUint(str_id, 10, 64)
	_, exists, err := services.GetBlog(id, bc.DB)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "blog not exists"})
		return
	}

	err = services.DeleteBlog(id, bc.DB)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}

func GetDB() *gorm.DB {
	var db *gorm.DB
	return db
}

func myGet(c *gin.Context) {
	var blogs []models.Blog
	db := GetDB()
	db.Find(&blogs)
	c.IndentedJSON(200, gin.H{"data": blogs})
}

func (bc *BlogController) RegisterBlogRoutes(rg *gin.RouterGroup) {
	blogRoute := rg.Group("/blog")
	{
		blogRoute.POST("/create", bc.CreateBlog)
		blogRoute.GET("/get-all", bc.GetBlogs)
		blogRoute.GET("/get-alls", myGet)
		blogRoute.GET("/get/:id", bc.GetBlog)
		blogRoute.PATCH("/update", bc.UpdateBlog)
		blogRoute.DELETE("/delete/:id", bc.DeleteBlog)
	}
}