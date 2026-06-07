package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vaForge/TaskManagerAPI/config"
	"github.com/vaForge/TaskManagerAPI/handlers"
	"github.com/vaForge/TaskManagerAPI/middleware"
	"github.com/vaForge/TaskManagerAPI/store"
)

// Now We change the flow from logger->recover->mux | mux -> handler ->response
// To : request -> CORS -> requestID
//
//	-> panic recovery -> logging -> mux -> handlers
func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// Create taskStore & taskHandler
	taskStore := store.NewStore()
	taskHandler := handlers.NewTaskHandler(taskStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)
	//collection routes
	mux.HandleFunc("/tasks", taskHandler.TaskHandler)
	// Item routes
	mux.HandleFunc("/tasks/", taskHandler.TaskByIDHandler)
	// home route
	mux.HandleFunc("/", homeHandler)

	finalHandler := middleware.Logging(
		middleware.Recover(
			middleware.RequestID(
				middleware.CORS(mux),
			),
		),
	)

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      finalHandler,
		ReadTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-stop
		logger.Info("shutdown signal received", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("graceful shutdown failed", "error", err)
		}
	}()

	fmt.Println("Server Listening on http://localhost:8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	logger.Info("server stopped cleanly")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, `{"status":"ok"}`)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "Welcome to Task Manager API bros")
}
