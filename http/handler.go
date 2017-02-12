package http

import (
	"net/http"

	"fmt"

	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/davepgreene/tokend/utils"
	"github.com/gorilla/mux"
	"github.com/meatballhat/negroni-logrus"
	"github.com/spf13/viper"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
	"gopkg.in/tylerb/graceful.v1"
)

const (
	// MaxRequestSize is the maximum accepted request size. This is to prevent
	// a denial of service attack where no Content-Length is provided and the server
	// is fed ever more data until it exhausts memory.
	MaxRequestSize = 32 * 1024 * 1024
)

// Handler returns an http.Handler for the API.
func Handler() error {
	r := mux.NewRouter()
	statsMiddleware := stats.New()
	r.HandleFunc("/stats", newAdminHandler(statsMiddleware).ServeHTTP)

	v1 := r.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/health", newHealthHandler().ServeHTTP)
	// 	app.get('/v1/health', require('./health')(storage));

	// Auth handlers
	v1.HandleFunc("/token/{name}", newTokenHandler().ServeHTTP)
	//   app.get('/v1/token/default', require('./token')(storage));

	// Secret handlers
	v1.HandleFunc(`/secret/{token}/{path:[a-zA-Z0-9=\-\/]+}`, newSecretHandler().ServeHTTP)
	//   app.get('/v1/secret/:token/:path(*)', require('./secret')(storage));

	v1.HandleFunc(`/cubbyhole/{token}/{path:[a-zA-Z0-9=\-\/]+}`, newCubbyHoleHandler().ServeHTTP)
	//   app.get('/v1/cubbyhole/:token/:path(*)', require('./cubbyhole')(storage));

	v1.HandleFunc("/credential/{token}/{mount}/{role}", newCredentialHandler().ServeHTTP)
	//   app.get('/v1/credential/:token/:mount/:role', require('./credential')(storage));

	v1.HandleFunc("/transit/{token}/decrypt", newTransitHandler().ServeHTTP)
	//   app.post('/v1/transit/:token/decrypt', require('./transit')(storage));

	// Define our 404 handler
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

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
	return server(conn, n).ListenAndServe()
}

func server(conn string, handler http.Handler) *graceful.Server {
	return &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:    conn,
			Handler: handler,
		},
	}
}

// notFoundHandler provides a standard response for unhandled paths
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	var b []byte

	w.Write(b)
}
