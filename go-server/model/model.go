package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	Id            int       `json:"id,omitempty"`
	Username      string    `json:"username,omitempty"`
	Email         string    `json:"email,omitempty"`
	Password_hash string    `json:"password_hash,omitempty"`
	Created_at    time.Time `json:"created_at,omitempty"`
	Updated_at    time.Time `json:"updated_at,omitempty"`
	Is_active     bool      `json:"is_active,omitempty"`
	Token         string    `json:"token,omitempty"`
}

type Task struct {
	Id          int       `json:"id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Completed   bool      `json:"completed,omitempty"`
	Created_at  time.Time `json:"created_at,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
	UserId      int       `json:"user_id,omitempty"`
	Uuid        uuid.UUID `json:"uuid,omitempty"`
}
