package tokens

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/adriangarcia1984/ecommerce-go/database"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string ` json:"email"`
	First_Name string
	Last_Name  string
	Uid        string
	jwt.RegisteredClaims
}

var UserData *mongo.Collection = database.UserData(database.Client, "Users")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedtoken string, signedrefreshtoken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		RegisteredClaims: jwt.RegisteredClaims{
			//ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),// update for v4 at v5 jwt-golang
			ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		},
	}
	refreshclaims := SignedDetails{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(1516239022, 0)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodES256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshtoken, err

}



func ValidateToken(signedtoken string) (claims *SignedDetails, msg string) {

	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token not valid"
		return
	}

	if claims.ExpiresAt.Time.Unix() < time.Now().Local().Unix() {
		msg = "the token not valid"
		return
	}
	return claims, msg
}

func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	update_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: update_at})
	upsert := true

	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_,err := UserData.UpdateOne(ctx, filter, bson.D{
			{Key: "$set", Value: updateobj},
		}, &opt)
		defer cancel()
		if err != nil{
			log.Panic(err)
		}
}
