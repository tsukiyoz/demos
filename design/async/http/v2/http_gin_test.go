package v2

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestGinHTTP(t *testing.T) {
	engine := gin.New()
	engine.GET("/", func(c *gin.Context) {
		// time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)+10))
		time.Sleep(time.Millisecond * time.Duration(30))
		c.String(http.StatusOK, "hello, world")
	})
	engine.Run(":8081")
}
