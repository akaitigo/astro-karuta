package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akaitigo/astro-karuta/backend/internal/handler"
	"github.com/akaitigo/astro-karuta/backend/internal/repository"
	"github.com/akaitigo/astro-karuta/backend/internal/seed"
	"github.com/akaitigo/astro-karuta/backend/internal/service"
	"github.com/akaitigo/astro-karuta/backend/internal/ws"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cardRepo := repository.NewInMemoryCardRepository()
	deckRepo := repository.NewInMemoryDeckRepository()
	collectionRepo := repository.NewInMemoryCollectionRepository(cardRepo)

	if err := seed.LoadCards(context.Background(), cardRepo); err != nil {
		log.Fatalf("failed to seed cards: %v", err)
	}
	log.Println("seeded card data")

	missionRepo := repository.NewInMemoryMissionRepository()

	seasonalSvc := service.NewSeasonalService(cardRepo, deckRepo)
	missionSvc := service.NewMissionService(missionRepo, cardRepo)

	cardSvc := service.NewCardService(cardRepo, deckRepo)
	cardSvc.SetSeasonalService(seasonalSvc)
	cardHandler := handler.NewCardHandler(cardSvc)

	collectionSvc := service.NewCollectionService(collectionRepo)
	collectionHandler := handler.NewCollectionHandler(collectionSvc)

	seasonalHandler := handler.NewSeasonalHandler(missionSvc, seasonalSvc)

	hub := ws.NewHub()
	gm := ws.NewGameManager(hub, cardRepo)
	wsHandler := handler.NewWSHandler(hub, gm)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})
	cardHandler.RegisterRoutes(mux)
	collectionHandler.RegisterRoutes(mux)
	seasonalHandler.RegisterRoutes(mux)
	wsHandler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
