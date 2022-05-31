package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	
	"github.com/gin-gonic/gin"

	_ "github.com/jekwulum/mongoGinAPI/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/jekwulum/mongoGinAPI/controllers"
	"github.com/jekwulum/mongoGinAPI/middlewares"
	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/services"

	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server 			*gin.Engine
	userservice		services.UserService
	usercontroller 	controllers.UserController
	ctx				context.Context
	usercollection	*mongo.Collection
	mongoclient		*mongo.Client
	err				error
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}


func init() {
	// mongodb setup
	ctx := context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connection to mongodb established")
	usercollection = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(usercollection, ctx)
	usercontroller = controllers.New(userservice)

	
	models.ConnectPostGRES_DB()
	setupLogOutput()
	
	server = gin.New()
	server.Use(gin.Recovery(), middlewares.Logger())
}


// @title Gin Swagger CRUD REST-API
// @version 1.0
// @description This is a go-gin server implementing a mongodb & postgres database.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	// setupLogOutput()
	// models.ConnectPostGRES_DB()
	defer mongoclient.Disconnect(ctx)

	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// server.GET("/v1/all_userss", GetUsers)
	basepath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basepath)
	controllers.RegisterBlogRoutes(basepath)
	log.Fatal(server.Run("localhost:8080"))
}