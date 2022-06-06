package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthRequired(tokens map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, err := c.Cookie("Authorization")
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		login, ok := tokens[auth]
		if !ok {
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		c.Set("login", login)

		c.Next()
	}
}
