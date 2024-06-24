package main

import (
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
    }, []string{"path", "method"})
)

func PrometheusMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Increment the counter on request
        route := c.FullPath()  // This gets the registered route path for metrics
        method := c.Request.Method
        httpRequestsTotal.With(prometheus.Labels{"path": route, "method": method}).Inc()

        c.Next() // Process the request
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
