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

	/*
		mux.HandleFunc("/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("test panic")
		})
	*/
	apiAuth := middleware.NewAPIKeyAuth(cfg.APIKeyHeader, cfg.APIKeys)

	finalHandler := middleware.Logging(
		middleware.Recover(
			middleware.RequestID(
				middleware.CORS(
					apiAuth.Middleware(mux),
				),
			),
		),
	)

	srv := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      finalHandler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	/*
		make(chan os.Signal, 1): You create a channel (a pipe for sending data between goroutines) of type os.Signal.
		The 1 means it is a "buffered" channel with a capacity of 1. The OS can drop exactly one signal into this pipe without waiting for anything to read it.

		signal.Notify: This binds the channel to the operating system. You are telling the Go runtime: "If the OS sends a SIGINT (Ctrl+C) or SIGTERM (kill command), don't crash the app.
		Instead, intercept it and push that signal into the stop channel."
	*/

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

	if r.URL.Path != "/healthz" {
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
