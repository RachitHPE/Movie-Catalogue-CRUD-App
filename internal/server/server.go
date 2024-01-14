package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"catalogue-app/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AppServerBase struct {
	Router    *gin.Engine
	server    *http.Server
	name      string
	isRunning bool
	mutex     sync.Mutex
}

// New implements AppServerBase.
func New(name string) *AppServerBase {
	return &AppServerBase{name: name}
}

func configureLogger() {
	log.ConfigureLogger()
}

func (app *AppServerBase) ConfigureAndStart() {
	app.Init()
	app.setupAPIWithRouter(context.Background())
	app.Start()
}

func (app *AppServerBase) Init() {
	configureLogger()

	app.Router = gin.New()
	gin.EnableJsonDecoderDisallowUnknownFields()
	app.Router.Use(gin.Recovery())
	app.Router.HandleMethodNotAllowed = true
	// Enabling Cors to allow your browser access the API.
	app.Router.Use(Cors())

	app.Router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "OK",
		})
	})

	app.Router.GET("/support/metrics", prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func (app *AppServerBase) setupAPIWithRouter(ctx context.Context) {

	router := app.Router.Group("catalogue")
	router.POST("/getData", JWTConfiguration(), nil)
}

// Start starts the Server for real.
func (app *AppServerBase) Start() {
	log.Infof(context.Background(), "Starting %s Server...", app.name)
	app.startGinServer()
	log.Infof(context.Background(), "%s server started successfully ...", app.name)

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof(context.Background(), "Shutting down %s server...", app.name)
	app.StopServer()
}

func (app *AppServerBase) startGinServer() {
	app.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", 9002),
		Handler:           app.Router,
		ReadHeaderTimeout: time.Second * time.Duration(20),
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		app.mutex.Lock()
		app.isRunning = true
		app.mutex.Unlock()

		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf(context.Background(), "listen: %s\n", err)
			os.Exit(1)
		}

		app.mutex.Lock()
		app.isRunning = false
		app.mutex.Unlock()
	}()
}

func (app *AppServerBase) StopServer() {
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Errorf(context.Background(), "Server Shutdown: %v", err)
	}

	log.Info(context.Background(), "Server stopped successfully ...")
}
