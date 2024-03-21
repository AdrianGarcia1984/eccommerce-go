package middleware

import (
	token "github.com/adriangarcia1984/ecommerce-go/tokens"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"not autorization header"})
			c.Abort()
			return
		}
		claims, err := token.ValidateToken(ClientToken)
		 if err != ""{
			c.JSON(http.StatusInternalServerError, gin.H{"error":err})
			c.Abort()
			return
		 }
		 c.Set("email", claims.Email)
		 c.Set("uid", claims.Uid)
		 c.Next()

	}
}
