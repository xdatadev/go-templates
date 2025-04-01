package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/xdatadev/{{ .Project }}/internal/config"
	"github.com/xdatadev/{{ .Project }}/internal/handlers"
	"github.com/xdatadev/{{ .Project }}/internal/models"
	"github.com/xdatadev/{{ .Project }}/internal/web"

	"github.com/xdatadev/superapp-packages/superapp-common/logger"
	samDB "github.com/xdatadev/superapp-packages/superappdb/database/postgres"
)

const traceIDKey = "traceID"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
			time.Sleep(2 * time.Hour)
			os.Exit(12) // Panic Recover
		}
	}()

	var log *logger.Logger

	// Declare Events Processed by Logger
	events := logger.LogEvents{}
	// Error: func(ctx context.Context, r logger.LogRecord) {
	// 	log.Info(ctx, "************* ALERT *************")
	// },

	ctx := context.Background()

	traceIDFunc := func(ctx context.Context) string {
		if traceID, ok := ctx.Value(traceIDKey).(string); ok {
			return traceID
		}
		return "NO-TRACE"
	}

	hostname, _ := os.Hostname()

	log = logger.New(logger.Config{
		ServiceName: "{{ .ProjectKebab }}",
		MinLevel:    logger.LevelDebug,
		TraceIDFunc: traceIDFunc,
		Events:      events,
		Outputs: logger.OutputConfig{
			UseStdout:     true,
			UseCloudWatch: false,
			CloudWatch: logger.CloudWatchConfig{
				LogGroup:    "{{ .Project }}",
				LogStream:   fmt.Sprintf("{{ .Project }}-%s", hostname),
				Region:      "us-east-1",
				Loglevel:    slog.LevelInfo,
				Credentials: nil, // Use Service Account Role
			},
		},
	}, events, true)

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "Startup", "Error", err)
		return
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	log.Info(ctx, "STARTUP", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ----------------------- LOAD CONFIGURATIONS -----------------------
	log.Info(ctx, "Loading Parameters")
	config := config.LoadParameters()
	log.Debug(ctx, "Startup", "Config", config)

	// ----------------------- DATABASE CONNECTION -----------------------
	dbConfig := samDB.DbPoolConfig{
		Host:              config.DB.Host,
		Port:              config.DB.Port,
		User:              config.DB.User,
		Password:          config.DB.Password,
		DbName:            config.DB.DbName,
		Region:            config.DB.Region,
		MaxPoolSize:       config.DB.MaxPoolSize,
		MinPoolSize:       1,
		MaxIdleTime:       30 * time.Minute,
		MaxLifeTime:       30 * time.Minute,
		HealthCheckPeriod: 30 * time.Second,
	}

	log.Info(ctx, "Creating DB Pool Connections")

	db, err := samDB.New(ctx, dbConfig, log)
	if err != nil {
		log.Error(ctx, "Startup DB Connection", "Error", err)
		os.Exit(1)
	}
	defer db.Close()

	// ----------------------- INIT SERVICE -----------------------
	appServices := models.AppServices{}
	log.Debug(ctx, "Services created", "services", appServices)

	// ----------------------- START API SERVER -----------------------
	handler := handlers.NewAppHandler(config, log, &appServices)
	srv := web.NewServer(":8080", handler)

	go func() {
		log.Info(ctx, "Servidor iniciado na porta 8080")
		err = srv.Start()
		if err != nil {
			log.Error(ctx, "Startup API Server", "Error", err)
			os.Exit(1)
		}
	}()

	// ----------------------- SHUTDOWN GRACEFULLY -----------------------
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	sig := <-shutdown

	log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(context.Background(), fmt.Sprintf("Erro durante o shutdown: %v", err))
	}

	defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

	return nil
}
