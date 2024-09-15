package storage

import (
	"errors"
	"log"
	"sync"

	"github.com/kittanut/go-task-queue/request"
)

type DefaultStorage struct {
	queue []request.RequestSchema
	mux   sync.Mutex
}

// NewStorage creates a new instance of Storage
func NewStorage() StorageInterface {
	return &DefaultStorage{
		queue: make([]request.RequestSchema, 0), // initialize an empty slice for the queue
		mux:   sync.Mutex{},                     // initialize a mutex for thread safety
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

	// Check if the queue is empty
	if len(s.queue) == 0 {
		return errors.New("queue is empty")
	}

	// Get the first request in the queue and perform request
	req := s.queue[0]
	err := req.MakeRequest()
	if err != nil {
		log.Println(req.Url)
	}

	// Remove the first request from the queue
	s.queue = s.queue[1:]
	return err
}
