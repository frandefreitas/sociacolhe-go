package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
	"sociacolhe/internal/services"
)

type RequestHandler struct{}

type createReq struct {
	IsAnonymous bool   `json:"isAnonymous"`
	Description string `json:"description" binding:"required,min=10"`
}

func (RequestHandler) Create(c *gin.Context) {
	var body createReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	area, urg := services.AutoTriage(body.Description)
	var patientID *uuid.UUID
	if !body.IsAnonymous {
		if uid, ok := c.Get("uid"); ok { id := uuid.MustParse(uid.(string)); patientID = &id }
	}
	req := models.ServiceRequest{ID: uuid.New(), PatientID: patientID, IsAnonymous: body.IsAnonymous, Description: body.Description, Area: &area, Urgency: urg, Status: models.ReqOpen}
	if err := db.DB.Create(&req).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, req)
}

func (RequestHandler) ListOpen(c *gin.Context) {
	var list []models.ServiceRequest
	if err := db.DB.Where("status IN ?", []models.RequestStatus{models.ReqOpen, models.ReqInTriage}).Order("urgency DESC, created_at ASC").Find(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, list)
}
