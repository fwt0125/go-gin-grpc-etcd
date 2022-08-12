package middlerware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func InitMiddleware(service []interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Keys = make(map[string]interface{})
		fmt.Println("middleware:")
		fmt.Println(service[0])
		c.Keys["user"] = service[0]
		c.Next()
	}
}
