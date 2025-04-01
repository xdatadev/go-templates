package web

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/xdatadev/{{ .Project }}/internal/handlers"
)

type Server struct {
	handlers *handlers.AppHandlers
	server   *http.Server
}

func NewServer(addr string, handlers *handlers.AppHandlers) *Server {

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(skipLoggingMiddleware())
	r.Use(traceIDMiddleware())
	r.Use(gin.Recovery())

	rootGroup := r.Group("/{{.Scaffold.Resource}}")
	{
		rootGroup.GET("/health", s.handlers.HealthCheck)
		v1 := rootGroup.Group("/v1")
		{
			//
		}
	}

	return &Server{
		handlers: handlers,
		server:   &http.Server{Addr: addr, Handler: r},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func skipLoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/metrics", "/readiness", "/"},
	})
}

func traceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("X-Amzn-Trace-Id") // AWS Trace ID
		if traceID == "" {
			traceID = c.Request.Header.Get("apigw-requestid")
			if traceID == "" {
				traceID = uuid.NewString()
			}
		}
		ctx := context.WithValue(c.Request.Context(), traceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func TimerDecorator(fn interface{}) interface{} {
	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Cria uma função que encapsula a lógica de medição
	wrapped := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		start := time.Now()

		// Chama a função original
		results := fnValue.Call(args)

		// Calcula o tempo de execução
		duration := time.Since(start).Milliseconds()
		fmt.Printf("Execution time: %d\n", duration)

		return results
	})

	return wrapped.Interface()
}
