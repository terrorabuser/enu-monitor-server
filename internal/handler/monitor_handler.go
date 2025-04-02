package handler

import (
	"golang_gpt/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MonitorHandler struct {
	monitorRepo *repository.MonitorRepository
}

// Конструктор
func NewMonitorHandler(monitorRepo *repository.MonitorRepository) *MonitorHandler {
	return &MonitorHandler{monitorRepo: monitorRepo}
}

func (h *MonitorHandler) GetAllMonitors(c *gin.Context) {
	monitors, err := h.monitorRepo.GetAllMonitors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, monitors)

}

