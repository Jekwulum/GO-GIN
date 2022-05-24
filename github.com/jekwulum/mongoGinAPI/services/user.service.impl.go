package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/jekwulum/mongoGinAPI/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx				context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl {
		usercollection: usercollection,
		ctx:			ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *models.User) error {
	result, err := u.usercollection.InsertOne(u.ctx, user)
	fmt.Println("result: ", result)
	return err
}

func (u *UserServiceImpl) GetUser(name *string) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "user_name", Value: name}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*models.User, error) {
	var users []*models.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		fmt.Println("get all error 1")
		return nil, err
	}

	for cursor.Next(u.ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			fmt.Println("get all error 2")
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("get all error 3")
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(users) == 0 {
		fmt.Println("get all error 3")
		return nil, errors.New("documents not found")
	}
	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *models.User) error {
	filter := bson.D{bson.E{Key: "user_name", Value: user.Name}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "user_name", Value: user.Name},
		bson.E{Key: "user_age", Value: user.Age},
		bson.E{Key: "user_addres", Value: user.Address},
	}}}
	result, _ := u.usercollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched document found for update")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(name *string) error {
	filter := bson.D{bson.E{Key: "user_name", Value: name}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}