package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"sociacolhe/internal/config"
	"sociacolhe/internal/handlers"
	"sociacolhe/internal/middleware"
	"sociacolhe/internal/ws"
)

func SetupRouter(cfg config.Config) *gin.Engine {
	r := gin.Default()
	jwt := middleware.JWTAuth{Secret: cfg.JWTSecret}

	// Health
	r.GET("/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"ok"}) })

	// Auth
	auth := handlers.AuthHandler{JWT: jwt, AllowSignup: cfg.AllowSignup}
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)

	// Admin
	admin := handlers.AdminHandler{}
	r.PUT("/admin/users/:userID/approve", jwt.Require("ADMIN"), admin.ApproveUser)

	// Requests
	req := handlers.RequestHandler{}
	r.POST("/requests", jwt.Require("PATIENT"), req.Create)
	r.GET("/requests/open", jwt.Require("PSYCHOLOGIST", "STUDENT", "ADMIN"), req.ListOpen)

	// Triage
	tri := handlers.TriageHandler{}
	r.PUT("/requests/:requestID/triage", jwt.Require("PSYCHOLOGIST", "ADMIN"), tri.Manual)

	// Assignment (aceitar atendimento)
	ass := handlers.AssignmentHandler{}
	r.POST("/requests/:requestID/assignments", jwt.Require("PSYCHOLOGIST", "STUDENT"), ass.Accept)

	// Sessions
	ses := handlers.SessionHandler{}
	r.POST("/sessions", jwt.Require("PSYCHOLOGIST", "STUDENT"), ses.Create)

	// Feedback
	fb := handlers.FeedbackHandler{}
	r.POST("/feedbacks", fb.Create) // aberto para pacientes anônimos também

	// Chat (room + websocket broadcast)
	chat := handlers.ChatHandler{}
	r.POST("/chat/rooms", chat.CreateRoom)
	h := ws.NewHub()
	r.GET("/ws/chat/:roomID", h.Join)

	return r
}
