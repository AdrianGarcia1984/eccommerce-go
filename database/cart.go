package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the products")
	ErrCantDecodeProducts = errors.New("can't find the products")
	ErrUserIdNotValid     = errors.New("this user not valid")
	ErrCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrCantRemoveItemCart = errors.New("cannot remove product to the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchase")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {
	searchFormDb, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		log.Print(err)
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFormDb.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProducts
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productId}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	//fetch cart of user
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	var getcarditems models.User
	var ordercart models.Order

	ordercart.Order_Id = primitive.NewObjectID()
	ordercart.Ordered_At = time.Now()
	ordercart.Order_Card = make([]models.ProductUser, 0)
	ordercart.Payment_Method.COD = true

	unwind := bson.D{{Key: " $unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
	currentItems, err:=userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil{
		panic(err)
	}
	//find the cart total
	 var getusercart []bson.M

	 if err = currentItems.All(ctx, &getusercart);err != nil{
		panic(err)
	 }
	 var total_price int32
	 for _, user_item := range getusercart{
		price:= user_item["total"]
		total_price = price.(int32)
	 }
	 ordercart.Price = int(total_price)

	 filter := bson.D{primitive.E{Key: "_id", Value: id}}
	 update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: ordercart}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	err=userCollection.FindOne(ctx,bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getcarditems)
	if err != nil {
		log.Println(err)
	}
	//create an order with the items
	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{"orders.$[].order_list": bson.M{"$each": getcarditems.UserCart}}}
	_,err=userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	//empty up the cart
	usercart_empty:= make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3:=bson.D{{Key:"$set", Value: bson.D{primitive.E{Key: "usercart", Value: usercart_empty}}}}
	_,err=userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error{
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdNotValid
	}
	var product_details models.ProductUser
	var order_details models.Order

	order_details.Order_Id = primitive.NewObjectID()
	order_details.Ordered_At = time.Now()
	order_details.Order_Card = make([]models.ProductUser, 0)
	order_details.Payment_Method.COD = true
	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&product_details)
	if err != nil{
		log.Println(err)
	}
	order_details.Price = product_details.Price

	filter:=bson.D{primitive.E{Key: "_id", Value: id}}
	update:=bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order_details}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
	}
	filter2:=bson.D{primitive.E{Key: "_id", Value: id}}
	update2:=bson.M{"$push": bson.M{"orders.$[].order_list": product_details}}
	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
	}
	
	return nil
}
