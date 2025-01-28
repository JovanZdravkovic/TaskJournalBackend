package db

import (
	"time"

	"github.com/google/uuid"
)

// Struct names with DB suffix are models that are mapped from a row in a database table
type TaskDB struct {
	Id          uuid.UUID `json:"id"`
	TaskName    string    `json:"taskName"`
	TaskIcon    string    `json:"taskIcon"`
	TaskDesc    string    `json:"taskDesc"`
	Deadline    time.Time `json:"deadline"`
	Starred     bool      `json:"starred"`
	Exec_status string    `json:"execStatus"`
	Created_at  time.Time `json:"createdAt"`
	Created_by  uuid.UUID `json:"createdBy"`
}

type TaskHistoryDB struct {
	Id          uuid.UUID `json:"id"`
	ExecRating  string    `json:"execRating"`
	ExecComment string    `json:"execComment"`
	TaskId      uuid.UUID `json:"taskId"`
}

type UserDB struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type AuthDB struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"userId"`
	ExpiresAt time.Time `json:"expiresAt"`
}
