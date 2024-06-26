package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/HanssenKai/go-api-observability-test/api"
	_ "github.com/HanssenKai/go-api-observability-test/docs"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"path", "method", "code"})
)


type responseWriter struct {
	gin.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a response writer to capture the status code
		wrappedWriter := &responseWriter{ResponseWriter: c.Writer, statusCode: http.StatusOK}

		// Replace the context's writer with the wrapped writer
		c.Writer = wrappedWriter

		// Process the request
		c.Next()

		// Get the route path, method, and status code
		route := c.FullPath()
		method := c.Request.Method
		statusCode := wrappedWriter.statusCode

		// Increment the counter with the path, method, and status code
		httpRequestsTotal.With(prometheus.Labels{
			"path":   route,
			"method": method,
			"code":   fmt.Sprintf("%d", statusCode),
		}).Inc()
	}
}

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
func main() {
	baseAPI := os.Getenv("BASE_API_URL")
	if baseAPI == "" {
		baseAPI = "https://kystdatahuset.no/ws/api"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(PrometheusMiddleware())

	api.SetupRoutes(router)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.Run(":8080")
}

func setupTestRouter() *gin.Engine {
	// Disable console logging during tests
	gin.SetMode(gin.TestMode)

	router := gin.New()
	api.SetupRoutes(router)

	return router
}
