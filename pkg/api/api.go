package api


import (
	"net/http"
	_ "net/http/pprof"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)


type Api struct {
	router         *mux.Router
	logger         *zap.Logger
	pool           *redis.Pool
	handler        http.Handler
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}
