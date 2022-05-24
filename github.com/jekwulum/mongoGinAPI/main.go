package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/jekwulum/mongoGinAPI/controllers"
	"github.com/jekwulum/mongoGinAPI/models"
	"github.com/jekwulum/mongoGinAPI/services"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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

// func postgresInit() {
// 	envErr := godotenv.Load(".env")

// 	if envErr != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	dbDriver := os.Getenv("DB_DRIVER")
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("PG_DB_PORT")
// 	db_name := os.Getenv("PG_DB_NAME")
// 	user := os.Getenv("PG_DB_USER")
// 	password := os.Getenv("PG_DB_PASSWORD")
// 	DB_URL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, db_name, password)
	
// 	pg_DB, dbErr := gorm.Open(dbDriver, DB_URL)
// 	if dbErr != nil {
// 		fmt.Println("cannot connect to %s database", dbDriver)
// 		log.Fatal("This is the error connecting to postgres: ", dbErr)
// 	} else {
// 		fmt.Println("successfully connected to %s database", dbDriver)
// 	}

// 	pg_DB.Debug().AutoMigrate(
// 		&models.Blog{},
// 	)
// }

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

	fmt.Println("connection to mongodb extablished")

	usercollection = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(usercollection, ctx)
	usercontroller = controllers.New(userservice)
	
	server = gin.Default()
	
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
		fmt.Println("cannot connect to %s database", dbDriver)
		log.Fatal("This is the error connecting to postgres: ", dbErr)
	} else {
		fmt.Println("successfully connected to %s database", dbDriver)
	}

	pg_DB.Debug().AutoMigrate(
		&models.Blog{},
	)
}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basepath)
	blogcontroller.RegisterBlogRoutes(basepath)
	log.Fatal(server.Run("localhost:8080"))
}