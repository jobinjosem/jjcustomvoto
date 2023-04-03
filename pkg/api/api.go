package api


import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

)

type Config struct {
	HttpClientTimeout     time.Duration `mapstructure:"http-client-timeout"`
	HttpServerTimeout     time.Duration `mapstructure:"http-server-timeout"`
	ServerShutdownTimeout time.Duration `mapstructure:"server-shutdown-timeout"`
	BackendURL            []string      `mapstructure:"backend-url"`
	UILogo                string        `mapstructure:"ui-logo"`
	UIMessage             string        `mapstructure:"ui-message"`
	UIColor               string        `mapstructure:"ui-color"`
	UIPath                string        `mapstructure:"ui-path"`
	DataPath              string        `mapstructure:"data-path"`
	ConfigPath            string        `mapstructure:"config-path"`
	CertPath              string        `mapstructure:"cert-path"`
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	SecurePort            string        `mapstructure:"secure-port"`
	PortMetrics           int           `mapstructure:"port-metrics"`
	Hostname              string        `mapstructure:"hostname"`
	H2C                   bool          `mapstructure:"h2c"`
	RandomDelay           bool          `mapstructure:"random-delay"`
	RandomDelayUnit       string        `mapstructure:"random-delay-unit"`
	RandomDelayMin        int           `mapstructure:"random-delay-min"`
	RandomDelayMax        int           `mapstructure:"random-delay-max"`
	RandomError           bool          `mapstructure:"random-error"`
	Unhealthy             bool          `mapstructure:"unhealthy"`
	Unready               bool          `mapstructure:"unready"`
	JWTSecret             string        `mapstructure:"jwt-secret"`
	CacheServer           string        `mapstructure:"cache-server"`
}

type Api struct {
	Router         *mux.Router
	Logger         *zap.Logger
	pool           *redis.Pool
	Config         *Config
	handler        http.Handler
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}
