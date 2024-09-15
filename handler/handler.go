package handler

import "github.com/gin-gonic/gin"

type QueueHandlerInterface interface {
	AddQueue(c *gin.Context)
}
