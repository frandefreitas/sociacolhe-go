package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type TriageHandler struct{}

type triageReq struct {
	Area    *string `json:"area"`
	Urgency *int    `json:"urgency"`
}

func (TriageHandler) Manual(c *gin.Context) {
	id, err := uuid.Parse(c.Param("requestID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }
	var body triageReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var r models.ServiceRequest
	if err := db.DB.First(&r, "id = ?", id).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	if body.Area != nil { r.Area = body.Area }
	if body.Urgency != nil { r.Urgency = *body.Urgency }
	r.Status = models.ReqInTriage
	if err := db.DB.Save(&r).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, r)
}
