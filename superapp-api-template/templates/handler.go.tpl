package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xdatadev/superapp-assistant/internal/config"
	"github.com/xdatadev/superapp-assistant/internal/models"
	"github.com/xdatadev/superapp-packages/superapp-common/logger"
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
		"message": "{{.ProjectName}} is running",
	})
}

