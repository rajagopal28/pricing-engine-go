package service

import (
	"net/http"
	"time"
	"log"
	"os"
	"os/signal"
	"context"

	"pricingengine/service/app"
	"pricingengine/service/rpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Start begins a chi-Mux'd net/http server on port 3000
type Service struct {
	Server *http.Server
}

// Start method takes care of handling the initial configs and starting the server based on the handler endpoints configured
// The service typically relies on the go-chi and chi middleware libraries in constructing a rest service
func (s * Service)Start(port string) {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(5 * time.Second))

	rpc := rpc.RPC{
		App: &app.App{},
	}
	if len(port) == 0 {
		// default port 3000
		port = "3000"
	}
	r.Post("/generate_pricing", rpc.GeneratePricing)
	r.Get("/generate_pricing", rpc.GeneratePricingConfig)
	s.ListenAndServe(":"+port, r)
}

// ListenAndServe method takes care of starting the server at the given port
// The service instance also listend to the SIGNAL that typically interprets the Interrupted
func (s *Service)ListenAndServe(port string, r http.Handler) {
		log.Println("Starting Server!")
		s.Server = &http.Server{Addr: port, Handler: r}
		if err := s.Server.ListenAndServe(); err != nil {
				// handle err
		}

    // Setting up signal capturing
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)

    // Waiting for SIGINT (pkill -2)
    <-stop
		log.Println("Received Server Stop!")
    s.Stop()
    // Wait for ListenAndServe goroutine to close.
}

// Stop method will check and stop the currently running http server
func (s * Service)Stop() {
	log.Println("Stopping Server!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if s.Server != nil {
		log.Println(" Initiating Server Shutdown!")
		if err := s.Server.Shutdown(ctx); err != nil {
			// handle err
			log.Println("Error while stopping Server!", err)
		}
	}
}
