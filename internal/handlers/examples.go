package handlers

import (
	"net/http"
	"time"

	"go-api-template/internal/models"

	"github.com/gin-gonic/gin"
)

// In-memory storage for examples
var examples = []models.Example{
	{ID: 1, Title: "Example 1", Description: "This is the first example", CreatedAt: time.Now()},
	{ID: 2, Title: "Example 2", Description: "This is the second example", CreatedAt: time.Now()},
}

var nextExampleID = 3

// GetExamples godoc
// @Summary Get all examples
// @Description Get a list of all examples
// @Tags examples
// @Accept json
// @Produce json
// @Success 200 {array} models.Example
// @Router /api/v1/examples [get]
func GetExamples(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":  examples,
		"count": len(examples),
	})
}

// CreateExample godoc
// @Summary Create a new example
// @Description Create a new example with the provided information
// @Tags examples
// @Accept json
// @Produce json
// @Param example body models.CreateExampleRequest true "Example information"
// @Success 201 {object} models.Example
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/examples [post]
func CreateExample(c *gin.Context) {
	var req models.CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	example := models.Example{
		ID:          nextExampleID,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	examples = append(examples, example)
	nextExampleID++

	c.JSON(http.StatusCreated, gin.H{"data": example})
}