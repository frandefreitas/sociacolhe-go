package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/models"
)

type AdminHandler struct{}

type approveReq struct { Approve bool `json:"approve"` }

func (AdminHandler) ApproveUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("userID"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }
	var body approveReq
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var u models.User
	if err := db.DB.First(&u, "id = ?", id).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"not found"}); return }
	u.IsApproved = body.Approve
	if err := db.DB.Save(&u).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "isApproved": u.IsApproved})
}
