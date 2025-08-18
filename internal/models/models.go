package models

import (
	"time"
	"github.com/google/uuid"
)

type UserRole string

const (
	Admin        UserRole = "ADMIN"
	Psychologist UserRole = "PSYCHOLOGIST"
	Student      UserRole = "STUDENT"
	Patient      UserRole = "PATIENT"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string    `json:"name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `gorm:"type:user_role" json:"role"`
	CRP          *string   `json:"crp"`
	Institution  *string   `json:"institution"`
	IsApproved   bool      `json:"isApproved"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type RequestStatus string

const (
	ReqOpen     RequestStatus = "OPEN"
	ReqInTriage RequestStatus = "IN_TRIAGE"
	ReqAssigned RequestStatus = "ASSIGNED"
	ReqClosed   RequestStatus = "CLOSED"
)

type ServiceRequest struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey" json:"id"`
	PatientID   *uuid.UUID    `gorm:"type:uuid" json:"patientId"`
	IsAnonymous bool          `json:"isAnonymous"`
	Description string        `json:"description"`
	Area        *string       `json:"area"`
	Urgency     int           `json:"urgency"`
	Status      RequestStatus `gorm:"type:request_status" json:"status"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
}

type Assignment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RequestID  uuid.UUID `gorm:"type:uuid" json:"requestId"`
	AssigneeID uuid.UUID `gorm:"type:uuid" json:"assigneeId"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Session struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RequestID        *uuid.UUID `gorm:"type:uuid" json:"requestId"`
	ProfessionalID   uuid.UUID  `gorm:"type:uuid" json:"professionalId"`
	StudentID        *uuid.UUID `gorm:"type:uuid" json:"studentId"`
	Date             time.Time  `json:"date"`
	DurationMinutes  int        `json:"durationMinutes"`
	Type             string     `json:"type"`
	Notes            *string    `json:"notes"`
	SupervisorCRP    *string    `json:"supervisorCrp"`
	CreatedAt        time.Time  `json:"createdAt"`
}

type Feedback struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RequestID         uuid.UUID `gorm:"type:uuid" json:"requestId"`
	RatingAcolhimento int      `json:"ratingAcolhimento"`
	RatingTecnica     int      `json:"ratingTecnica"`
	RatingResultado   int      `json:"ratingResultado"`
	Comment           *string   `json:"comment"`
	FlagIssue         bool      `json:"flagIssue"`
	CreatedAt         time.Time `json:"createdAt"`
}

type ChatRoom struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	RequestID *uuid.UUID `gorm:"type:uuid" json:"requestId"`
	CreatedAt time.Time  `json:"createdAt"`
}

type ChatMessage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoomID    uuid.UUID `gorm:"type:uuid;index" json:"roomId"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
