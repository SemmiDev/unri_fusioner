package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	uf "unri_fusioner"
	"unri_fusioner/sinta"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HTTPServer struct {
	port int
	host string

	router *chi.Mux
	sinta  *sinta.Sinta
}

func NewHTTPServer(config *uf.Config, opts ...func(*HTTPServer) error) (*HTTPServer, error) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/api/ping"))

	httpServer := &HTTPServer{
		host:   "127.0.0.1",
		port:   3000,
		router: r,
	}

	// if opts is not empty, run all opts
	// and apply it to httpServer config
	if len(opts) != 0 {
		for _, opt := range opts {
			if err := opt(httpServer); err != nil {
				return nil, err
			}
		}
	}

	sintaScrapper := sinta.Sinta{
		SintaDomain: config.SintaDomain,
	}

	r.Get("/api/sinta/authors/profile/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := uf.CastToInt(chi.URLParam(r, "id"))

		authorProfile, err := sintaScrapper.ScrapeAuthorProfile(id)
		if err != nil {
			if errors.Is(err, sinta.ErrAuthorNotFound) {
				errorResponse(err.Error(), http.StatusNotFound, w)
				return
			}
			errorResponse(err.Error(), http.StatusInternalServerError, w)
			return
		}

		successResponse(authorProfile, http.StatusOK, w)
	})

	return httpServer, nil
}

func (s *HTTPServer) Start() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	listen := fmt.Sprintf("%s:%d", s.host, s.port)
	server := &http.Server{
		Addr:    listen,
		Handler: s.router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %s\n", err)
		}
	}()

	// wait for SIGINT or SIGTERM
	<-stop
	log.Println("shutting down server...")

	// wait for 5 seconds for the server to shut down
	// completely running jobs
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error shutting down server: %s\n", err)
	}
	log.Println("server gracefully stopped")
}

func errorResponse(message string, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	response := map[string]string{"message": message}
	dataJSON, _ := json.MarshalIndent(response, "", "  ")
	_, _ = w.Write(dataJSON)
}

func successResponse(data any, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	response, _ := json.MarshalIndent(data, "", "  ")
	_, _ = w.Write(response)
}
