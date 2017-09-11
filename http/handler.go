package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/davepgreene/tokend/api"
	"github.com/davepgreene/tokend/utils"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

var requiredHeaders map[string]string = map[string]string{
	"Content-Type": "application/json; charset=utf-8",
}

// Handler returns an http.Handler for the API.
func Handler(storage *api.Storage) {
	r := mux.NewRouter()
	statsMiddleware := stats.New()
	r.HandleFunc("/stats", newAdminHandler(statsMiddleware).ServeHTTP)

	v1 := r.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/health", Adapt(newHealthHandler(storage), setHeaders()).ServeHTTP)

	// Auth handlers
	v1.HandleFunc("/token/{name}", Adapt(newTokenHandler(storage), setHeaders()).ServeHTTP)

	// Secret handlers
	//v1.HandleFunc(`/secret/{token}/{path:[a-zA-Z0-9=\-\/]+}`, newSecretHandler(storage).ServeHTTP)
	//v1.HandleFunc(`/cubbyhole/{token}/{path:[a-zA-Z0-9=\-\/]+}`, newCubbyHoleHandler(storage).ServeHTTP)
	//v1.HandleFunc("/credential/{token}/{mount}/{role}", newCredentialHandler(storage).ServeHTTP)
	v1.HandleFunc("/transit/{token}/decrypt", Adapt(newTransitHandler(storage), checkHeaders(requiredHeaders), setHeaders()).ServeHTTP)

	v1.HandleFunc("/kms/decrypt", Adapt(newKMSHandler(storage), checkHeaders(requiredHeaders), setHeaders()).ServeHTTP)

	// Define our 404 handler
	r.NotFoundHandler = Adapt(http.HandlerFunc(notFoundHandler), setHeaders())

	// Add middleware handlers
	n := negroni.New()
	n.Use(negroni.NewRecovery())

	if viper.GetBool("log.requests") {
		n.Use(negronilogrus.NewCustomMiddleware(utils.GetLogLevel(), utils.GetLogFormatter(), "requests"))
	}

	n.Use(statsMiddleware)
	n.UseHandler(r)

	// Set up connection
	conn := fmt.Sprintf("%s:%d", viper.GetString("service.host"), viper.GetInt("service.port"))
	log.Info(fmt.Sprintf("Listening on %s", conn))

	// Bombs away!
	server := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Addr:              conn,
		Handler:           n,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown(server)
}

func shutdown(s *http.Server) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	var timeout time.Duration

	if s.ReadHeaderTimeout != 0 {
		timeout = s.ReadHeaderTimeout
	} else if s.IdleTimeout != 0 {
		timeout = s.IdleTimeout
	} else {
		timeout = s.ReadTimeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Infof("Shutdown with timeout: %s", timeout)

	if err := s.Shutdown(ctx); err != nil {
		log.Error(err)
	} else {
		log.Info("Server stopped")
	}
}

// notFoundHandler provides a standard response for unhandled paths
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	w.Write([]byte(""))
}
