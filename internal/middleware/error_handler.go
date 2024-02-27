package middleware

import (
	"meepShopTest/internal/apierr"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	for _, err := range c.Errors {
		apierr, ok := err.Unwrap().(apierr.ApiErr)
		if ok {
			c.JSON(-1, apierr)
			break
		}

		c.JSON(-1, gin.H{"error": err.Error()})
	}
}
