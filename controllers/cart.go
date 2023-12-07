package controllers

import "github.com/gin-gonic/gin"

func AddToCard() gin.HandlerFunc

func RemoveItem() gin.HandlerFunc

func GetItemFromCart() gin.HandlerFunc

func BuyItemFromCart() gin.HandlerFunc

func InstantBuy() gin.HandlerFunc
