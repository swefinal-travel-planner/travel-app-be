package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/bean"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
)

type HealthHandler struct {
	db          database.Db
	redisClient bean.RedisClient
}

func NewHealthHandler(db database.Db, redisClient bean.RedisClient) *HealthHandler {
	return &HealthHandler{
		db:          db,
		redisClient: redisClient,
	}
}

// @Summary Health Check
// @Description Checks the health of the application by verifying database and Redis connections
// @Tags Health
// @Accept json
// @Produce json
// @Router /health [get]
// @Success 200 {object} map[string]string "Healthy status"
// @Failure 503 {object} map[string]string "Unhealthy status with error details"
func (h *HealthHandler) Check(c *gin.Context) {
	// Check database connection
	if err := h.db.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	// Check Redis connection
	if err := h.redisClient.Set(c, "health_check", "ok"); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  "redis connection failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}
