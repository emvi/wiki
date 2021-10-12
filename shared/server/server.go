package server

import (
	"context"
	"emviwiki/shared/config"
	"github.com/emvi/logbuch"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	shutdownTimeout = time.Second * 30
)

// Start starts a new HTTP REST server.
// It allows configuring read/write timeouts and TLS through environment variables.
// The optional shutdown function is called right before shutdown and can be used to clean up.
func Start(handler http.Handler, shutdown func()) {
	c := config.Get()
	logbuch.Info("Starting server...")
	logbuch.Info("Using HTTP read/write timeouts", logbuch.Fields{"write_timeout": c.Server.HTTP.Timeout.Write, "read_timeout": c.Server.HTTP.Timeout.Read})

	server := &http.Server{
		Handler:      handler,
		Addr:         c.Server.Host,
		WriteTimeout: time.Duration(c.Server.HTTP.Timeout.Write) * time.Second,
		ReadTimeout:  time.Duration(c.Server.HTTP.Timeout.Read) * time.Second,
	}

	go func() {
		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		logbuch.Info("Shutting down server...")

		if shutdown != nil {
			shutdown()
		}

		ctx, _ := context.WithTimeout(context.Background(), shutdownTimeout)

		if err := server.Shutdown(ctx); err != nil {
			logbuch.Fatal("Error shutting down server gracefully", logbuch.Fields{"err": err})
		}
	}()

	if c.Server.HTTP.TLS {
		logbuch.Info("TLS enabled")

		if err := server.ListenAndServeTLS(c.Server.HTTP.TLSCert, c.Server.HTTP.TLSKey); err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	} else {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logbuch.Fatal(err.Error())
		}
	}
}
