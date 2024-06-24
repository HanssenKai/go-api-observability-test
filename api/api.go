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

// locations updates list of ship locations and stores in database
// @Summary get list of locations and update database
// @Tags health
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
