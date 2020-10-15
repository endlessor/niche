package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RtbListenCheck test send request
func RtbListenCheck(c *gin.Context) {
	fmt.Println("request body", time.Now().String())
	c.JSON(http.StatusOK, "status ok")
}
