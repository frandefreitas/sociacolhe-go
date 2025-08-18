package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type FeedbackHandler struct{}

type feedbackReq struct {
	RequestID string  `json:"requestId" binding:"required"`
	RatingAcolhimento int `json:"ratingAcolhimento" binding:"required,min=1,max=5"`
	RatingTecnica int `json:"ratingTecnica" binding:"required,min=1,max=5"`
	RatingResultado int `json:"ratingResultado" binding:"required,min=1,max=5"`
	Comment *string `json:"comment"`
	FlagIssue bool `json:"flagIssue"`
}

func (FeedbackHandler) Create(c *gin.Context) {
	var body feedbackReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	reqID := uuid.MustParse(body.RequestID)
	f := models.Feedback{
		ID: uuid.New(), RequestID: reqID, RatingAcolhimento: body.RatingAcolhimento,
		RatingTecnica: body.RatingTecnica, RatingResultado: body.RatingResultado,
		Comment: body.Comment, FlagIssue: body.FlagIssue,
	}
	if err := db.DB.Create(&f).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, f)
}
