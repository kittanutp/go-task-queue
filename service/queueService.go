package service

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kittanut/go-task-queue/request"
	"github.com/kittanut/go-task-queue/storage"
)

type QueueService struct {
	storage storage.StorageInterface
	stop    chan bool // Channel to signal stopping the service
}

// NewQueueService creates a new QueueService and starts queue management in the background
func NewQueueService(storage storage.StorageInterface) *QueueService {
	s := &QueueService{
		storage: storage,
		stop:    make(chan bool),
	}

	// Handle shutdown signals
	go s.listenForShutdown()

	// Start the background queue manager after a delay
	go func() {
		time.Sleep(5 * time.Second)
		s.ManageQueue()
	}()

	return s
}

// ProcessNewQueue pushes a new request into the queue
func (q *QueueService) ProcessNewQueue(req request.RequestSchema) error {
	if err := q.storage.Push(req); err != nil {
		return err
	}
	return nil
}

// ManageQueue continuously fetches and processes requests from the queue
func (q *QueueService) ManageQueue() {
	for {
		select {
		case <-q.stop:
			log.Println("Stopping queue management")
			return
		default:
			if !q.storage.IsEmpty() {
				if err := q.storage.Fetch(); err != nil {
					log.Println(err)
				}
			} else {
				time.Sleep(2 * time.Second) // Wait before checking again
			}
		}
	}
}

// Stop signals the queue management to stop
func (q *QueueService) Stop() {
	q.stop <- true
}

// listenForShutdown listens for OS shutdown signals and stops the queue manager
func (q *QueueService) listenForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Block until a signal is received
	log.Println("Received shutdown signal, stopping queue service...")

	// Stop the queue manager gracefully
	q.Stop()
	log.Println("Queue service stopped")
}

func (q *QueueService) CheckEmptyQueue() bool {
	return q.storage.IsEmpty()
}
