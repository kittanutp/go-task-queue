package service

import "github.com/kittanut/go-task-queue/request"

type QueueServiceInterface interface {
	ProcessNewQueue(request request.RequestSchema) error
	ManageQueue()
	Stop()
}
