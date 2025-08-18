package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type ChatHandler struct{}

type createRoomReq struct { RequestID *string `json:"requestId"` }

func (ChatHandler) CreateRoom(c *gin.Context) {
	var body createRoomReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var rid *uuid.UUID
	if body.RequestID != nil { id := uuid.MustParse(*body.RequestID); rid = &id }
	room := models.ChatRoom{ID: uuid.New(), RequestID: rid}
	if err := db.DB.Create(&room).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, room)
}
