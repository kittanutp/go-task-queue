package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kittanut/go-task-queue/request"
	"github.com/kittanut/go-task-queue/service"
)

type QueueHandler struct {
	queueService service.QueueServiceInterface
}

func NewQueueHTTPHandler(queueService service.QueueServiceInterface) *QueueHandler {
	return &QueueHandler{
		queueService: queueService,
	}
}

func (qh *QueueHandler) AddQueue(c *gin.Context) {
	var request request.RequestSchema
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatusJSON(400, fmt.Sprintf("invalid json request as %v", err))
		return
	}
	if err := qh.queueService.ProcessNewQueue(request); err != nil {
		c.AbortWithStatusJSON(400, fmt.Sprintf("unable to add request as %v", err))
		return
	}
	c.SecureJSON(200, gin.H{"data": nil})
}

func (qh *QueueHandler) CheckQueue(c *gin.Context) {
	c.SecureJSON(200, gin.H{"is_queue_empty": qh.queueService.CheckEmptyQueue()})
}
