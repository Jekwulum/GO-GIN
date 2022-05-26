package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	// "time"

	_ "github.com/jekwulum/mongoGinAPI/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/jekwulum/mongoGinAPI/controllers"
	"github.com/jekwulum/mongoGinAPI/middlewares"
	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/services"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server 			*gin.Engine
	blogcontroller	controllers.BlogController
	userservice		services.UserService
	usercontroller 	controllers.UserController
	ctx				context.Context
	usercollection	*mongo.Collection
	mongoclient		*mongo.Client
	err				error
	pg_DB			*gorm.DB
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}


func init() {
	// mongo
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
	
	server = gin.New()
	server.Use(gin.Recovery(), middlewares.Logger())
	

	// postgres
	envErr := godotenv.Load(".env")

	if envErr != nil {
		log.Fatal("Error loading .env file")
	}
	dbDriver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("PG_DB_PORT")
	db_name := os.Getenv("PG_DB_NAME")
	user := os.Getenv("PG_DB_USER")
	password := os.Getenv("PG_DB_PASSWORD")
	DB_URL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, db_name, password)
	
	pg_DB, dbErr := gorm.Open(dbDriver, DB_URL)
	if dbErr != nil {
		fmt.Println("cannot connect to postgress database")
		log.Fatal("This is the error connecting to postgres: ", dbErr)
	} else {
		fmt.Println("successfully connected to postgres database")
		}
	// if !pg_DB.HasTable(&models.Blog{}) {
	// 	fmt.Println("checkkkkkkkkkkkkkkkkkkkkkkkk")
	// 	err := pg_DB.CreateTable(&models.Blog{})
	// 	if err != nil {
	// 		log.Fatal("Table already exists")
	// 	} else {
	// 		fmt.Println("Table created")
	// 	}
	// }

	pg_DB.AutoMigrate(&models.Blog{})
	pg_DB.AutoMigrate(&models.Cust_User{})

	// newBlog := models.Blog{Title: "blog 1",  Content: "Content 2", CreatedAt: time.Now()}
	// result := pg_DB.Create(&newBlog)
	// if result.Error != nil {
	// 	fmt.Println("create error ----------------------->", result.Error)
	// } else {
	// 	fmt.Println("blog id---------------------> ", newBlog.ID)
	// 	fmt.Println("rows affected  ------------------->", result.RowsAffected)
	// }

	user_result := pg_DB.Create(&models.Cust_User{Name: "Bambi"})
	if user_result != nil {
		fmt.Println("error---------------> ", user_result.Error)
	} else {
		fmt.Println("user--------------->", user_result.RowsAffected)
	}
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
	setupLogOutput()
	defer mongoclient.Disconnect(ctx)

	// url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	basepath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basepath)
	blogcontroller.RegisterBlogRoutes(basepath)
	log.Fatal(server.Run("localhost:8080"))
}