package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// HashPassword
func HashPassword(password string) string

// verfyPassword
func VerifyPassword(userPassword string, givvenPassword string) (bool, string)

//login

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user, founduser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//check user exist

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "login or password incorrect"})
			return
		}
		PasswordValid, msg := VerifyPassword(*user.Password, *founduser.Password)
		defer cancel()
		if !PasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refresToken, _ := generate.TokenGenerator(*founduser.Email, *founduser.Firt_Name, *founduser.Last_Name, founduser.User_Id)
		defer cancel()
		generate.UpdateAllTokens(token, refresToken,founduser.User_Id)

		c.JSON(http.StatusFound, founduser)
	}
}

// signup
func Singup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		valErr := Validate.Struct(user)
		if valErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": valErr})
			return
		}

		count, err := UserCollection.CountCocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exist"})
		}
		count, err = UserCollection.CountCocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "phone number already exist"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()
		token, refresToken, _ := generate.TokenGenerator(*user.Email, *user.Firt_Name, *user.Last_Name, *&user.User_Id)
		user.Token = token
		user.Refresh_Token = refresToken
		user.UserCard = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr = UserCollection.InsertOne(ctx, user)

		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "the user didn´t created"})
			return
		}

		defer cancel()

		c.JSON(http.StatusCreated, "succesfully signed in!")

	}
}

// ProductViewerAdmin
func ProductViewerAdmin() gin.HandlerFunc

// searchProduct
func SearchProduct() gin.HandlerFunc

// searchProductByQuery
func SearchProductByQuery() gin.HandlerFunc
