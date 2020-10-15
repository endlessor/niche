package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// StatusCheck tests health
func StatusCheck(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
