package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/services"
)

type UserController struct {
	UserService services.UserService
}

func New(userservice services.UserService) UserController {
	return UserController{
		UserService: userservice,
	}
}

// routes

// CreateUser godoc
// @Summary Create user
// @Description create and add new user to the database
// @Tags Users
// @Accept json
// @Produce json
// @Param name formData string true "user's name"
// @Param age formData string true "user's age"
// @Param state formData string true "user's name"
// @Param city formData string true "user's age"
// @Param pincode formData string true "user's name"
// @Success 201 {object} map[string]string{}
// @Failure 400	{object} map[string]interface{}
// @Failure 502	{object} map[string]interface{}
// @Router /v1/user/create [post]
func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "success"})
}

// GetUser godoc
// @Summary Get user
// @Description fetch a user from the database
// @Tags Users
// @Accept json
// @Produce json
// @Param name path string true "user's name"
// @Success 200 {object} map[string]string{}
// @Failure 404	{object} map[string]interface{}
// @Router /v1/user/get/{name} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	username := ctx.Param("name")
	user, err := uc.UserService.GetUser(&username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, 
		gin.H{"message": "success", "data": user})
}

// GetUsers godoc
// @Summary Get All Users
// @Description get all users in the database
// @Tags Users
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/user/get-all [get]
func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	fmt.Println("users: ", users)
	if err != nil {
		fmt.Println("get all error 4")
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, users)
}


// UpdateUser godoc
// @Summary Update user
// @Description update existing user in the database
// @Tags Users
// @Accept json
// @Produce json
// @Success 201 {object} map[string]string{}
// @Failure 400	{object} map[string]interface{}
// @Failure 502	{object} map[string]interface{}
// @Router /v1/user/create [patch]
func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusCreated, gin.H{"message": "success"})
}


// DeleteUser godoc
// @Summary Delete user
// @Description delete a user from the database
// @Tags Users
// @Accept json
// @Produce json
// @Param name path string true "user's name"
// @Success 200 {object} map[string]string{}
// @Failure 404	{object} map[string]interface{}
// @Router /v1/user/get/{name} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("name")
	err := uc.UserService.DeleteUser(&username)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup){
	userRoute := rg.Group("/user")
	{
		userRoute.POST("/create", uc.CreateUser)
		userRoute.GET("/get-all", uc.GetAll)
		userRoute.GET("/get/:name", uc.GetUser)
		userRoute.PATCH("/update", uc.UpdateUser)
		userRoute.DELETE("/delete/:name", uc.DeleteUser)
	}
}