package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type SessionHandler struct{}

type sessionReq struct {
	RequestID    *string  `json:"requestId"`
	StudentID    *string  `json:"studentId"`
	Date         time.Time `json:"date" binding:"required"`
	Duration     int      `json:"durationMinutes" binding:"required,min=1"`
	Type         string   `json:"type" binding:"required"`
	Notes        *string  `json:"notes"`
	SupervisorCRP *string `json:"supervisorCrp"`
}

func (SessionHandler) Create(c *gin.Context) {
	var body sessionReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	profID := uuid.MustParse(c.GetString("uid"))
	var reqID *uuid.UUID
	if body.RequestID != nil { id := uuid.MustParse(*body.RequestID); reqID = &id }
	var studID *uuid.UUID
	if body.StudentID != nil { id := uuid.MustParse(*body.StudentID); studID = &id }
	s := models.Session{
		ID: uuid.New(), RequestID: reqID, ProfessionalID: profID, StudentID: studID,
		Date: body.Date, DurationMinutes: body.Duration, Type: body.Type, Notes: body.Notes, SupervisorCRP: body.SupervisorCRP,
	}
	if err := db.DB.Create(&s).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, s)
}
