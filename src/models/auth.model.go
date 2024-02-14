package models

import (
	"context"
	"fmt"
	"log"
	"root/src/common"
	"root/src/config"
	"root/src/core/db"
	"root/src/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func AuthModel() *BaseModel {
	mod := &BaseModel{
		ModelConstructor: &common.ModelConstructor{
			Collection: db.GetMongoDb().Collection("users"),
		},
	}

	return mod
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (mod *BaseModel) CurrentUser(ctx *gin.Context) (User, error) {
	var user User

	_, err := utils.ExtractTokenID(ctx)
	fmt.Println(user)

	if err != nil {
		return user, err
	}

	//user = UsersModel().GetOneUser(user_id)

	if err != nil {
		return user, err
	}

	return user, nil
}

func GenerateToken(user User) (string, error) {

	tokenLifespan, err := strconv.Atoi(config.LoadConfig("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.LoadConfig("API_SECRET")))

}

func (mod *BaseModel) LoginCheck(username, password string) (result string, err error) {
	var user User
	err = mod.Collection.FindOne(context.TODO(), bson.M{"userName": username}).Decode(&user)
	if err != nil || user.ID == primitive.NilObjectID {
		return "", err
	}
	log.Println("user => ", user.ID, password, user.Password)
	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	log.Println("user => ", user.ID, password, user.Password)
	token, err := GenerateToken(user)
	log.Println("token => ", token)
	if err != nil {
		return "", err
	}

	return token, nil
}
