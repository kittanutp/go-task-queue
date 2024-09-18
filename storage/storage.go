package storage

import (
	"os"
	"sync"

	"github.com/kittanut/go-task-queue/request"
)

type DefaultStorage struct {
	queue    []request.RequestSchema
	mux      sync.Mutex
	occupied chan os.Signal
}

// NewStorage creates a new instance of Storage
func NewStorage() StorageInterface {
	return &DefaultStorage{
		queue:    make([]request.RequestSchema, 0), // initialize an empty slice for the queue
		mux:      sync.Mutex{},
		occupied: make(chan os.Signal), // initialize a mutex for thread safety
	}
}

// Push adds a request to the queue in a thread-safe manner
func (s *DefaultStorage) Push(req request.RequestSchema) error {
	s.mux.Lock()         // acquire the lock
	defer s.mux.Unlock() // ensure the lock is released after the operation

	// Append the request to the queue
	s.queue = append(s.queue, req)
	return nil
}

// Fetch retrieves and removes the first request from the queue in a thread-safe manner
func (s *DefaultStorage) Fetch() error {
	s.mux.Lock()         // acquire the lock
	defer s.mux.Unlock() // ensure the lock is released after the operation

	// Get the first request in the queue and perform request
	req := s.queue[0]
	err := req.MakeRequest()

	// Remove the first request from the queue
	s.queue = s.queue[1:]
	return err
}

func (s *DefaultStorage) IsEmpty() bool {
	return len(s.queue) == 0
}

func (s *DefaultStorage) QueueOccupied() chan os.Signal {
	return s.occupied
}
