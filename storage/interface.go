package storage

import "github.com/kittanut/go-task-queue/request"

type StorageInterface interface {
	Push(request request.RequestSchema) error
	Fetch() error
}
