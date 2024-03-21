package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/database"
	"github.com/adriangarcia1984/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCard() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "succesfully added to the cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		//similar addtocard, check the user and product
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))

			return
		}

		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "succesfully removed item from cart")

	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		var fillercart  models.User
		user_id := c.Query("id")

		if user_id == "" {
			log.Println("user_id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search index"})
			c.Abort()
			return
		}
		usert_id, err:= primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usert_id}}).Decode(&fillercart)
		if err != nil {
			c.IndentedJSON(500,"not found")
			return
		}
		//special search in the database aggretation query
		filter_match := bson.D{{Key: "$match",Value: bson.D{primitive.E{Key:"_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind",Value: bson.D{primitive.E{Key:"path", Value: "$usercart"}}}}
		grouping:=bson.D{{Key: "$group", Value: bson.D{primitive.E{Key:"_id", Value: "$_id"}, {Key: "total", Value:bson.D{primitive.E{Key:"$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err:=UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind,grouping})
		if err != nil {
		log.Println(err)
		}
		var listing []bson.M

		if err = pointcursor.All(ctx,&listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing{
			c.IndentedJSON(200,json["total"])
			c.IndentedJSON(200, fillercart.UserCart)
		}
		ctx.Done()

	}
}

func (app Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")
		if userQueryId == "" {
			log.Panicln("user id is empty")
			_=c.AbortWithError(http.StatusBadRequest, errors.New("userId is empty"))
		}
		var ctx, cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		err:=database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err !=nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200,"successfully placed the order")
	}
}

func (app Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId, userQueryId := checkUserProductId(c)

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productId, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(http.StatusOK, "succesfully placed order")
	}
}

// func for check user replace on every func
func checkUserProductId(ctx *gin.Context) (primitive.ObjectID, string) {
	productQueryId := ctx.Query("id")
	if productQueryId == "" {
		log.Println("product id is empty")
		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
	}

	userQueryId := ctx.Query("userID")
	if userQueryId == "" {
		log.Println("user id is empty")

		_ = ctx.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
	}
	productId, err := primitive.ObjectIDFromHex(productQueryId)

	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	return productId, userQueryId
}
