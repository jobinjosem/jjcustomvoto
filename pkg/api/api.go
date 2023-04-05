package api

import (
	"net/http"
	_ "net/http/pprof"
	"time"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"github.com/jobinjosem/jjcustomvoto/pkg/fscache"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	healthy int32
	ready   int32
	watcher *fscache.Watcher
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
	Pool           *redis.Pool
	Config         *Config
	Handler        http.Handler
	Tracer         trace.Tracer
	TracerProvider *sdktrace.TracerProvider
}


func NewServer(config *Config) (*Api, error) {
	srv := &Api{
		Router: mux.NewRouter(),
		Config: config,
	}

	return srv, nil
}

func (a *Api) StartMetricsServer() {
	if a.Config.PortMetrics > 0 {
		mux := http.DefaultServeMux
		mux.Handle("/metrics", promhttp.Handler())
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%v", a.Config.PortMetrics),
			Handler: mux,
		}

		srv.ListenAndServe()
	}
}

// func (a *Api) RegisterMiddlewares() {
// 	// prom := NewPrometheusMiddleware()
// 	// a.Router.Use(prom.Handler)
// 	// otel := NewOpenTelemetryMiddleware()
// 	// a.Router.Use(otel)
// 	// httpLogger := NewLoggingMiddleware(a.Logger)
// 	// a.Router.Use(httpLogger.Handler)
// 	a.Router.Use(VersionMiddleware)
// 	if a.Config.RandomDelay {
// 		randomDelayer := NewRandomDelayMiddleware(a.Config.RandomDelayMin, a.Config.RandomDelayMax, a.Config.RandomDelayUnit)
// 		a.Router.Use(randomDelayer.Handler)
// 	}
// 	if a.Config.RandomError {
// 		a.Router.Use(RandomErrorMiddleware)
// 	}
// }


// func (a *Api) registerMiddlewares() {
// 	httpLogger := NewLoggingMiddleware(a.Logger)
// 	a.Router.Use(httpLogger.Handler)
// 	a.Router.Use(VersionMiddleware)
// 	if a.Config.RandomDelay {
// 		randomDelayer := NewRandomDelayMiddleware(a.Config.RandomDelayMin, a.Config.RandomDelayMax, a.Config.RandomDelayUnit)
// 		a.Router.Use(randomDelayer.Handler)
// 	}
// 	if a.Config.RandomError {
// 		a.Router.Use(RandomErrorMiddleware)
// 	}
// }

func (a *Api) PrintRoutes() {
	a.Router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		queriesTemplates, err := route.GetQueriesTemplates()
		if err == nil {
			fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		}
		queriesRegexps, err := route.GetQueriesRegexp()
		if err == nil {
			fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		}
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
}