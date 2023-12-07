package controllers

import (
	"github.com/gin-gonic/gin"
)

// HashPassword
func HashPassword(password string) string

// verfyPassword
func VerifyPassword(userPassword string, givvenPassword string) (bool, string)

//login

func Login() gin.HandlerFunc

// signup
func Singup() gin.HandlerFunc

// ProductViewerAdmin
func ProductViewerAdmin() gin.HandlerFunc

// searchProduct
func SearchProduct() gin.HandlerFunc

// searchProductByQuery
func SearchProductByQuery() gin.HandlerFunc
