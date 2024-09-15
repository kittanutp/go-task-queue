package main

import (
	"github.com/kittanut/go-task-queue/config"
	"github.com/kittanut/go-task-queue/server"
	"github.com/kittanut/go-task-queue/storage"
)

func main() {
	config := config.GetConfig()
	storage := storage.NewStorage()
	server.NewGinServer(config, storage).Start()
}
