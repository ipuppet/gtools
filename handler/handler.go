package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonStatus(c *gin.Context, err error) {
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]bool{"status": true})
}

func JsonStatusWithData(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": true,
		"data":   data,
	})
}
