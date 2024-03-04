package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var addresses models.Address
		user_id := c.Query("id")
		if user_id == "" {
			log.Println("user_id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search index"})
			c.Abort()
			return
		}
		address, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
			return
		}

		addresses.Address_Id = primitive.NewObjectID()

		if err = c.BindJSON(&addresses); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$address_id"},{Key:  "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter,unwind,group})
		if err != nil{
			c.IndentedJSON(500, "internal server error")
		}
		var addressInfo  []bson.M
		if err = pointcursor.All(ctx, &addressInfo); err!= nil{
			panic(err)
		}
		var size int32

		for _, address_No:= range addressInfo{
			count := address_No["count"]
			size =count.(int32)
		}
		if size < 2{
			filter := bson.D{primitive.E{Key: "_id", Value: address}}
			update := bson.D{{Key: "$push", Value:  bson.D{primitive.E{Key: "address", Value: addresses}}}}
			_,err:=UserCollection.UpdateOne(ctx, filter, update)
			if err != nil{
				log.Println(err)
			}

		}else{
			c.IndentedJSON(400, "not allowed")
		}
		defer cancel()
		ctx.Done()
	}
}

func EditHomeAddress() gin.HandlerFunc

func EditWorkAddress() gin.HandlerFunc

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		addresses := make([]models.Address, 0)
		user_id := c.Query("id")

		if user_id == "" {
			log.Println("user_id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid search index"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal server error")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: addresses}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "wrong command")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "successfully deleted")
	}
}
