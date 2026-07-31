package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Énumérations ---

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleTech     Role = "tech"
	RoleStandard Role = "standard"
)

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMid      Priority = "mid"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// --- Tables ---

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Firstname    string    `gorm:"size:255"`
	Lastname     string    `gorm:"size:255"`
	Email        string    `gorm:"size:255;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         Role      `gorm:"type:role;default:'standard';not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"size:255;not null"`
}

type Ticket struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string     `gorm:"size:255;not null"`
	Description string     `gorm:"type:text"`
	Status      Status     `gorm:"type:status;default:'open';not null"`
	Priority    Priority   `gorm:"type:priority;default:'low';not null"`
	CategoryID  *uint      `gorm:"index"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid;index"`
	AssignedTo  *uuid.UUID `gorm:"type:uuid;index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ResolvedAt  *time.Time

	// Relations (GORM)
	Category *Category `gorm:"foreignKey:CategoryID"`
	Creator  *User     `gorm:"foreignKey:CreatedBy"`
	Assignee *User     `gorm:"foreignKey:AssignedTo"`
}

type TicketAssignmentHistory struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	TicketID           uuid.UUID  `gorm:"type:uuid;not null;index"`
	AssignedToUserID   *uuid.UUID `gorm:"type:uuid;index"`
	AssignedByUserID   *uuid.UUID `gorm:"type:uuid;index"`
	AssignedAt         time.Time  `gorm:"default:now()"`
	EndedAt            *time.Time
	StatusAtAssignment Status `gorm:"type:status"`

	// Relations
	Ticket         Ticket `gorm:"foreignKey:TicketID"`
	AssignedToUser *User  `gorm:"foreignKey:AssignedToUserID"`
	AssignedByUser *User  `gorm:"foreignKey:AssignedByUserID"`
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"column:token_hash;not null"`
	Revoked   bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time

	User User `gorm:"foreignKey:UserID"`
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
