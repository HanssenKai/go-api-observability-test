package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes setups all the routes for the application
func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.GET("/health", healthCheck)
		v1.GET("/fail", fail)
		v1.POST("/update", locations)
		v1.GET("/list", locations)
	}
}

// healthCheck responds with the health status of the application
// @Summary Check health status
// @Description Returns the health status of the application.
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthStatus
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, HealthStatus{"healthy"})
}

// fail always responds with a 503 Service Unavailable status
// @Summary Check failure status
// @Description Always returns a failure status indicating that the service is unavailable.
// @Tags failure
// @Accept  json
// @Produce  json
// @Failure 503 {object} HealthStatus
// @Router /fail [get]
func fail(c *gin.Context) {
    c.JSON(http.StatusServiceUnavailable, HealthStatus{"unhealthy"})
}

// locations updates list of ship locations and stores in database
// @Summary get list of locations and update database
// @Tags update
// @Accept  json
// @Produce  json
// @Success 200 {object} HealthStatus
// @Router /update [post]
func locations(c *gin.Context) {
	result, err := fetchLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch locations"})
		return
	}
	c.String(http.StatusOK, result)
}

// HealthStatus represents the health status of the application
type HealthStatus struct {
	Status string `json:"status"`
}
