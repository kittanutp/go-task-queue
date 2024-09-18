package storage

import (
	"os"

	"github.com/kittanut/go-task-queue/request"
)

type StorageInterface interface {
	Push(request request.RequestSchema) error
	Fetch() error
	IsEmpty() bool
	QueueOccupied() chan os.Signal
}
