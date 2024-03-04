package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/database"
	"github.com/adriangarcia1984/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

// HashPassword
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

// verfyPassword
func VerifyPassword(userPassword string, givvenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givvenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		valid = false
		msg = "login or password is incorrect"
	}

	return valid, msg
}

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
		generate.UpdateAllTokens(token, refresToken, founduser.User_Id)

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
func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "someting went wrong, try after some time")
			return
		}
		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)

	}
}

// searchProductByQuery
func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProduct []models.Product
		queryParam := c.Query("name")
		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search index"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		searchquerydb, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(400, "someting went wrongwhile fetching the data")
			return
		}
		err = searchquerydb.All(ctx, &searchProduct)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)
		if err := searchquerydb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(200, searchProduct)
	}
}
