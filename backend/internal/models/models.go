package models

import (
	"time"

	"github.com/google/uuid"
)

//go:generate tygo generate

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

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Firstname    string    `gorm:"size:255" json:"firstname"`
	Lastname     string    `gorm:"size:255" json:"lastname"`
	Email        string    `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash string    `gorm:"column:password_hash;not null" json:"-"`
	Role         Role      `gorm:"type:role;default:'standard';not null" json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Category struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
}

type Ticket struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      Status     `gorm:"type:status;default:'open';not null" json:"status"`
	Priority    Priority   `gorm:"type:priority;default:'low';not null" json:"priority"`
	CategoryID  *uint      `gorm:"index" json:"categoryId,omitempty"`
	CreatedBy   *uuid.UUID `gorm:"type:uuid;index" json:"createdBy,omitempty"`
	AssignedTo  *uuid.UUID `gorm:"type:uuid;index" json:"assignedTo,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	ResolvedAt  *time.Time `json:"resolvedAt,omitempty"`

	Category            *Category                 `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Creator             *User                     `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Assignee            *User                     `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	AssignmentHistories []TicketAssignmentHistory `gorm:"foreignKey:TicketID" json:"assignmentHistories,omitempty"`
}

type TicketAssignmentHistory struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID           uuid.UUID  `gorm:"type:uuid;not null;index" json:"ticketId"`
	AssignedToUserID   *uuid.UUID `gorm:"type:uuid;index" json:"assignedToUserId,omitempty"`
	AssignedByUserID   *uuid.UUID `gorm:"type:uuid;index" json:"assignedByUserId,omitempty"`
	AssignedAt         time.Time  `gorm:"default:now()" json:"assignedAt"`
	EndedAt            *time.Time `json:"endedAt,omitempty"`
	StatusAtAssignment Status     `gorm:"type:status" json:"statusAtAssignment"`

	Ticket         Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	AssignedToUser *User  `gorm:"foreignKey:AssignedToUserID" json:"assignedToUser,omitempty"`
	AssignedByUser *User  `gorm:"foreignKey:AssignedByUserID" json:"assignedByUser,omitempty"`
}

func (TicketAssignmentHistory) TableName() string {
	return "ticket_assignment_history"
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"userId"`
	TokenHash string    `gorm:"column:token_hash;not null" json:"-"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	ExpiresAt time.Time `gorm:"not null" json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (RefreshToken) TableName() string {
	return "refresh_token"
}
