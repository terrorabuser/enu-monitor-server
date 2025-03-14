package handler

import (
	"golang_gpt/internal/repository"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ApiHandler struct {
	repo *repository.ApiRepository
}

func NewApiHandler(repo *repository.ApiRepository) *ApiHandler {
	return &ApiHandler{repo: repo}
}

// 📌 Получение списка зданий
func (h *ApiHandler) GetBuildings(c *gin.Context) {
	buildings, err := h.repo.GetBuildings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}
	c.JSON(http.StatusOK, buildings)
	// c.JSON(http.StatusOK, gin.H{"buildings": buildings})

}

// 📌 Получение этажей по зданию
func (h *ApiHandler) GetFloors(c *gin.Context) {
	building := c.Param("building")

	floors, err := h.repo.GetFloors(building)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}
	c.JSON(http.StatusOK, floors)
}

// 📌 Получение примечаний по зданию и этажу
func (h *ApiHandler) GetNotes(c *gin.Context) {
	building := c.Param("building")
	floorStr := c.Param("floor")

	floor, err := strconv.Atoi(floorStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат этажа"})
		return
	}




	notes, err := h.repo.GetNotes(building, floor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}
	c.JSON(http.StatusOK, notes)
}
