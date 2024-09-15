package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kittanut/go-task-queue/config"
	"github.com/kittanut/go-task-queue/handler"
	"github.com/kittanut/go-task-queue/service"
	"github.com/kittanut/go-task-queue/storage"
)

type ginServer struct {
	app     *gin.Engine
	storage storage.StorageInterface
	config  *config.Config
}

func NewGinServer(config *config.Config, storage storage.StorageInterface) Server {
	app := gin.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"Content-Type,access-control-allow-origin, access-control-allow-headers"},
	}))

	app.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	return &ginServer{
		app:     app,
		storage: storage,
		config:  config,
	}
}

func (g *ginServer) Start() {
	g.app.Use(gin.Recovery())
	g.app.Use(gin.Logger())
	g.initializeQueueHttpHandler()
	serverUrl := fmt.Sprintf(":%d", g.config.Server.Port)

	server := &http.Server{
		Addr:    serverUrl,
		Handler: g.app,
	}

	go func() {
		log.Printf("Starting server on %s", serverUrl)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	g.handleShutdown(server)
}

func (g *ginServer) initializeQueueHttpHandler() {
	queueService := service.NewQueueService(g.storage)
	queuehandler := handler.NewQueueHTTPHandler(queueService)
	queueRoutes := g.app.Group("/queue")
	{
		queueRoutes.POST("new", queuehandler.AddQueue)

	}

}

func (g *ginServer) handleShutdown(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit // Block until we receive a signal
	log.Println("Shutting down server...")

	// Create a context with a timeout for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop the queue service gracefully

	// Shutdown the HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
