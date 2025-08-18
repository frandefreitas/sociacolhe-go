package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sociacolhe/internal/db"
	"sociacolhe/internal/middleware"
	"sociacolhe/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	JWT middleware.JWTAuth
	AllowSignup bool
}

type registerReq struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role string `json:"role" binding:"required,oneof=PSYCHOLOGIST STUDENT PATIENT"`
	CRP *string `json:"crp"`
	Institution *string `json:"institution"`
}

type loginReq struct { Email string `json:"email" binding:"required,email"`; Password string `json:"password" binding:"required"` }

func (a AuthHandler) Register(c *gin.Context) {
	if !a.AllowSignup { c.JSON(http.StatusForbidden, gin.H{"error":"signup disabled"}); return }
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	pwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	u := models.User{ID: uuid.New(), Name: req.Name, Email: strings.ToLower(req.Email), PasswordHash: string(pwd), Role: models.UserRole(req.Role)}
	if req.CRP != nil { u.CRP = req.CRP }
	if req.Institution != nil { u.Institution = req.Institution }
	// Paciente é auto-aprovado; demais aguardam admin
	if u.Role == models.Patient { u.IsApproved = true }
	if err := db.DB.Create(&u).Error; err != nil { c.JSON(http.StatusConflict, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "isApproved": u.IsApproved})
}

func (a AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	var u models.User
	if err := db.DB.Where("email = ?", strings.ToLower(req.Email)).First(&u).Error; err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return }
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid credentials"}); return }
	if !u.IsApproved && u.Role != models.Patient { c.JSON(http.StatusForbidden, gin.H{"error":"awaiting approval"}); return }
	t, _ := a.JWT.Issue(u.ID.String(), string(u.Role))
	c.JSON(http.StatusOK, gin.H{"token": t, "role": u.Role, "id": u.ID})
}
