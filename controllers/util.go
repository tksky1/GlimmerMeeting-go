package controllers

import "github.com/gin-gonic/gin"

func set400(c *gin.Context, msg string) {
	c.JSON(400, gin.H{
		"message": "请求有误 " + msg,
	})
}

func set200(c *gin.Context, msg string) {
	c.JSON(200, gin.H{
		"message": msg,
	})
}

func set500(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"message": "内部错误：" + err.Error(),
	})
}
