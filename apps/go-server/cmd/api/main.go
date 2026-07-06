// Command api is the BitwiseLearn Go API server entrypoint.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/config"
	"github.com/bitwiselearn/go-server/internal/db"
	appmw "github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/server"
	"github.com/bitwiselearn/go-server/internal/services/blob"
	"github.com/bitwiselearn/go-server/internal/services/piston"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

func main() {
	cfg := config.Load()

	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbClient, err := db.Connect(startCtx, cfg.DatabaseURL, cfg.DBName(), cfg.MongoServerSelectionTimeout)
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	log.Printf("Database connected: %s", cfg.DBName())

	if err := dbClient.EnsureIndexes(startCtx); err != nil {
		log.Fatalf("index sync failed: %v", err)
	}
	log.Printf("Indexes ensured")

	jwtMgr := auth.New(cfg.JWTAccessSecret, cfg.JWTRefreshSecret, cfg.ResetTokenSecret)
	authMW := appmw.NewAuth(jwtMgr, dbClient)

	store, err := blob.New(cfg.AzureStorageConnectionString, cfg.AzureStorageContainer)
	if err != nil {
		log.Fatalf("blob store init failed: %v", err)
	}

	publisher := queue.NewPublisher(cfg.MQClient)
	defer publisher.Close()

	e := server.New(server.Deps{
		Config:    cfg,
		DB:        dbClient,
		JWT:       jwtMgr,
		Auth:      authMW,
		Store:     store,
		Publisher: publisher,
		OTP:       auth.NewMemoryOTPStore(),
		Blocklist: auth.NewMemoryBlocklist(),
		Piston:    piston.New(cfg.CodeExecutionServer),
	})

	// Start server in a goroutine so we can handle graceful shutdown.
	go func() {
		addr := ":" + cfg.Port
		log.Printf("API listening on %s", addr)
		if err := e.Start(addr); err != nil && err.Error() != "http: Server closed" {
			log.Printf("server stopped: %v", err)
		}
	}()

	// Wait for SIGINT/SIGTERM (ACA sends SIGTERM on scale-in).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
	if err := dbClient.Disconnect(shutdownCtx); err != nil {
		log.Printf("db disconnect error: %v", err)
	}
	log.Printf("shutdown complete")
}
