package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/dipmed/backend-chat/internal"
	"github.com/dipmed/backend-chat/internal/config"
	"github.com/dipmed/backend-chat/internal/llm"
	"github.com/dipmed/backend-chat/internal/sessions"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configPath, err := config.ConfigPath()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	provider, err := llm.NewGeminiProvider(ctx, cfg.LLM.APIKey, cfg.LLM.Model)
	if err != nil {
		log.Fatalf("failed to init llm provider: %v", err)
	}

	store := sessions.NewMemoryStore()
	srv := internal.NewServer(provider, store)

	httpServer := &http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: srv.Handler(),
	}

	go func() {
		log.Printf("listening on %s", cfg.HTTPServer.Address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")
	httpServer.Close()
}
