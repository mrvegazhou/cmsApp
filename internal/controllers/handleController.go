package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type handleController struct {
}

func NewHandleController() handleController {
	return handleController{}
}

func (con handleController) Handle(c *gin.Context) {
	if c.GetHeader("Accept") == "application/json" {

		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "url not fund",
			"data":    "",
		})
	} else {
		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  404,
				"message": "url not fund",
				"data":    "",
			})
		} else {
			c.HTML(http.StatusOK, "main/error.html", gin.H{
				"title":   "出错了",
				"status":  404,
				"message": "url not fund",
			})
		}

	}
}
