package server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var port = 8080

type Server struct {
	port   int
	logger log.Logger
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: port,
		logger: *log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller: true,
		}),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
