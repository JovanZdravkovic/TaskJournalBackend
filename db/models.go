package db

import (
	"time"

	"github.com/google/uuid"
)

// Struct names with DB suffix are models that are mapped from a row in a database table
// Pointers are used for fields that can have null values, because pointers can have null values

type TaskDB struct {
	Id          uuid.UUID  `json:"id"`
	TaskName    string     `json:"taskName"`
	TaskIcon    string     `json:"taskIcon"`
	TaskDesc    string     `json:"taskDesc"`
	Deadline    *time.Time `json:"deadline"`
	Starred     bool       `json:"starred"`
	Exec_status string     `json:"execStatus"`
	Created_at  time.Time  `json:"createdAt"`
	Created_by  uuid.UUID  `json:"createdBy"`
}

type TaskHistoryDB struct {
	Id          uuid.UUID `json:"id"`
	ExecRating  *int      `json:"execRating"`
	ExecComment *string   `json:"execComment"`
	TaskId      uuid.UUID `json:"taskId"`
	TaskName    string    `json:"taskName"`
	TaskIcon    string    `json:"taskIcon"`
}

type TaskHistoryPut struct {
	ExecRating  *int    `json:"execRating"`
	ExecComment *string `json:"execComment"`
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

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserPost struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserGet struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserPut struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type TaskPost struct {
	TaskName  string     `json:"taskName"`
	TaskIcon  string     `json:"taskIcon"`
	TaskDesc  string     `json:"taskDesc"`
	Deadline  *time.Time `json:"deadline"`
	Starred   bool       `json:"starred"`
	CreatedBy uuid.UUID  `json:"createdBy"`
}

type TaskPut struct {
	TaskName string     `json:"taskName"`
	TaskIcon string     `json:"taskIcon"`
	TaskDesc string     `json:"taskDesc"`
	Deadline *time.Time `json:"deadline"`
	Starred  bool       `json:"starred"`
}

type Id struct {
	Id uuid.UUID `json:"id"`
}

type Success struct {
	Success bool `json:"success"`
}
