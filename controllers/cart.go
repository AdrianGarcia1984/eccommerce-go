package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/database"
	"github.com/gin-gonic/gin"
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
		
	}
}

func BuyItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func InstantBuy() gin.HandlerFunc {
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


//func for check user replace on every func
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
