package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type AssignmentHandler struct{}

func (AssignmentHandler) Accept(c *gin.Context) {
	reqID, err := uuid.Parse(c.Param("requestID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }
	uid := uuid.MustParse(c.GetString("uid"))
	var r models.ServiceRequest
	if err := db.DB.First(&r, "id = ?", reqID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	ass := models.Assignment{ID: uuid.New(), RequestID: reqID, AssigneeID: uid}
	if err := db.DB.Create(&ass).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	r.Status = models.ReqAssigned
	_ = db.DB.Save(&r)
	c.JSON(http.StatusCreated, ass)
}
