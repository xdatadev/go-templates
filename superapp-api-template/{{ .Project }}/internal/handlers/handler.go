package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xdatadev/superapp-packages/superapp-common/logger"
	"github.com/xdatadev/{{ .Project }}/internal/config"
	"github.com/xdatadev/{{ .Project }}/internal/models"
)

type AppHandler struct {
	config   *config.Config
	logger   *logger.Logger
	services *models.AppServices
}

func NewAppHandler(config *config.Config, logger *logger.Logger, services *models.AppServices) *AppHandler {
	return &AppHandler{
		config:   config,
		logger:   logger,
		services: services,
	}
}

func (h *AppHandler) HealthCheck(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"message": "{{ .Project }} is running",
	})
}
