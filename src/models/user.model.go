package models

import (
	"context"
	"fmt"
	"root/src/common"
	"root/src/core/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func UsersModel() *BaseModel {
	mod := &BaseModel{
		ModelConstructor: &common.ModelConstructor{
			Collection: db.GetMongoDb().Collection("users"),
		},
	}

	return mod
}

func MigrateUsers() {
	db.GetGorm().AutoMigrate(&User{})
}

// models definitions

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	Username  string             `db:"title" json:"title"`
	Email     *string
	FirstName string `db:"firstName" json:"firstName"`
	LastName  string `db:"lastName" json:"lastName"`
	Password  string `db:"password" json:"password" validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserForm struct {
	Username  string `form:"Username" json:"Username" binding:"required"`
	Email     string `form:"Email" json:"Email" binding:"required"`
	FirstName string `form:"FirstName" json:"FirstName" binding:"required"`
	LastName  string `form:"LastName" json:"LastName" binding:"required"`
	Password  string `form:"Password" json:"Password" binding:"required"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}
type UsersFindParam struct {
	ID primitive.ObjectID `bson:"_id"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// models methods
func (mod *BaseModel) GetOneUser(userId primitive.ObjectID) User {
	var user User

	result := mod.Gorm.Limit(1).Where("id = ?", userId).Find(&user)

	if result.Error != nil {
		fmt.Println(result.Error)
		return user
	}

	return user
}

func (mod *BaseModel) GetAllUsers(limit int64, skip int64, search string) (user []User, err error) {
	var results []User
	opts := options.Find().SetLimit(limit).SetSkip(skip)

	filter := bson.D{{}}

	cur, err := mod.Collection.Find(context.TODO(), filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}
		panic(err)
	}

	if err = cur.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	count, err := mod.Collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	fmt.Println(count)

	return results, err
}

func (mod *BaseModel) UpdateUser(param UsersFindParam, body User) User {
	body.ID = param.ID

	result := mod.Gorm.Save(&body)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return body
}

func (mod *BaseModel) CreateUser(body User) User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	body.Password = string(hashedPassword)

	result := mod.Gorm.Create(&body)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return body
}

func (mod *BaseModel) DeleteUser(param UsersFindParam) bool {
	var user User
	user.ID = param.ID

	result := mod.Gorm.Delete(&user)

	fmt.Println(result)
	if result.Error != nil || result.RowsAffected == 0 {
		return false
	}

	return true
}
