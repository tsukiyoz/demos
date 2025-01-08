package main

import (
	"log"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	var id int32
	r.GET("/sleep", func(c *gin.Context) {
		atomic.AddInt32(&id, 1)
		time.Sleep(30 * time.Millisecond)
		log.Printf("ID:%d Time:%d\n", atomic.LoadInt32(&id), time.Now().UnixMilli())
		c.JSON(200, gin.H{
			"message": "success",
		})
	})
	r.Run(":8080")
}
