package main

import (
	"log"
	"os"

	"github.com/adriangarcia1984/ecommerce-go/controllers"
	"github.com/adriangarcia1984/ecommerce-go/database"
	"github.com/adriangarcia1984/ecommerce-go/middleware"
	"github.com/adriangarcia1984/ecommerce-go/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "User"))

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocard", app.AddToCard())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))

}
