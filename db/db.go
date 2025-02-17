package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DatabaseService struct {
	pool *pgxpool.Pool
}

func NewDatabaseService(dbPool *pgxpool.Pool) *DatabaseService {
	return &DatabaseService{
		pool: dbPool,
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func MatchPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// TASKS

func (dbService *DatabaseService) GetTasks(userId uuid.UUID) ([]TaskDB, error) {
	rows, err := dbService.pool.Query(context.Background(), "SELECT t.* FROM task t WHERE t.created_by = $1", userId)
	if err != nil {
		return nil, errors.New("error while getting tasks from database")
	} else {
		var tasks []TaskDB
		for rows.Next() {
			var task TaskDB
			err := rows.Scan(
				&task.Id,
				&task.TaskName,
				&task.TaskIcon,
				&task.TaskDesc,
				&task.Deadline,
				&task.Starred,
				&task.Exec_status,
				&task.Created_at,
				&task.Created_by,
			)
			if err != nil {
				return nil, errors.New("error while iterating dataset")
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}
}

func (dbService *DatabaseService) GetStarredTasks(userId uuid.UUID) ([]TaskDB, error) {
	rows, err := dbService.pool.Query(context.Background(), "SELECT t.* FROM task t WHERE t.starred = true AND t.created_by = $1", userId)
	if err != nil {
		return nil, errors.New("error while getting tasks from database")
	} else {
		var tasks []TaskDB
		for rows.Next() {
			var task TaskDB
			err := rows.Scan(
				&task.Id,
				&task.TaskName,
				&task.TaskIcon,
				&task.TaskDesc,
				&task.Deadline,
				&task.Starred,
				&task.Exec_status,
				&task.Created_at,
				&task.Created_by,
			)
			if err != nil {
				return nil, errors.New("error while iterating dataset")
			}
			tasks = append(tasks, task)
		}
		return tasks, nil
	}
}

func (dbService *DatabaseService) GetTask(taskId uuid.UUID, userId uuid.UUID) (*TaskDB, error) {
	row := dbService.pool.QueryRow(context.Background(), "SELECT t.* FROM task t WHERE t.id = $1 AND t.created_by = $2", taskId, userId)
	if row == nil {
		return nil, errors.New("task with given uuid doesnt exist")
	}
	var task TaskDB
	err := row.Scan(
		&task.Id,
		&task.TaskName,
		&task.TaskIcon,
		&task.TaskDesc,
		&task.Deadline,
		&task.Starred,
		&task.Exec_status,
		&task.Created_at,
		&task.Created_by,
	)
	if err != nil {
		return nil, errors.New("error while iterating dataset")
	}
	return &task, nil
}

func (dbService *DatabaseService) CreateTask(task TaskPost) (*uuid.UUID, error) {
	var taskId uuid.UUID
	err := dbService.pool.QueryRow(
		context.Background(),
		"INSERT INTO task(task_name, task_icon, task_desc, deadline, starred, exec_status, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		task.TaskName,
		task.TaskIcon,
		task.TaskDesc,
		task.Deadline,
		task.Starred,
		task.ExecStatus,
		task.CreatedBy,
	).Scan(&taskId)
	if err != nil {
		return nil, errors.New("error while creating task")
	}
	return &taskId, nil
}

// // TODO: Updating a single task will have multiple functions for different kinds of updates
// func (dbService *DatabaseService) PutTask() (string, error) {}

// func (dbService *DatabaseService) DeleteTask() (string, error) {}

// AUTH

func (dbService *DatabaseService) GetLoggedInUser(tokenId uuid.UUID) (*uuid.UUID, error) {
	var userId uuid.UUID
	err := dbService.pool.QueryRow(context.Background(), "SELECT u.id FROM \"user\" u JOIN user_auth ua ON u.id = ua.user_id WHERE ua.id = $1 AND ua.expires_at > CURRENT_TIMESTAMP", tokenId).Scan(&userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("invalid token")
		}
		return nil, errors.New("unexpected error")
	}
	return &userId, nil
}

func (dbService *DatabaseService) CreateToken(credentials Credentials) (*AuthDB, error) {
	tx, err := dbService.pool.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var userId uuid.UUID
	var passwordHash string
	err = tx.QueryRow(context.Background(), "SELECT u.id, u.password FROM \"user\" u WHERE u.username = $1", credentials.Username).Scan(&userId, &passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user doesn't exist")
		}
		return nil, errors.New("unexpected error")
	}

	passwordCheck := MatchPassword(credentials.Password, passwordHash)
	if !passwordCheck {
		return nil, errors.New("invalid credentials")
	}

	var authRow AuthDB
	err = tx.QueryRow(context.Background(), "INSERT INTO user_auth(user_id, expires_at) VALUES ($1, $2) RETURNING *", userId, time.Now().Add(7*24*time.Hour)).Scan(&authRow.Id, &authRow.UserId, &authRow.ExpiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("error creating the token")
		}
		return nil, errors.New("unexpected error")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &authRow, nil
}

func (dbService *DatabaseService) InvalidateToken(tokenId uuid.UUID) {
	cmdTag, err := dbService.pool.Exec(context.Background(), "DELETE FROM user_auth ua WHERE ua.id = $1", tokenId)
	if err != nil {
		log.Printf("Error while deleting tag")
	}
	log.Printf("%v", cmdTag)
}

// USERS

func (dbService *DatabaseService) GetUserInfo(userId uuid.UUID) (*UserGet, error) {
	var user UserGet
	err := dbService.pool.QueryRow(context.Background(), "SELECT u.username, u.email. u.created_at FROM \"user\" u WHERE u.id = $1", userId).Scan(
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user doesn't exist")
		}
		return nil, errors.New("unexpected error")
	}
	return &user, nil
}

func (dbService *DatabaseService) CreateUser(user UserPost) (*uuid.UUID, error) {
	tx, err := dbService.pool.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	var cnt int
	err = tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM \"user\" u WHERE u.username = $1", user.Username).Scan(&cnt)
	if err != nil {
		return nil, errors.New("unexpected error")
	}
	if cnt > 0 {
		return nil, errors.New("username taken")
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("error hashing password")
	}

	var userId uuid.UUID
	err = tx.QueryRow(
		context.Background(),
		"INSERT INTO \"user\"(username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username,
		user.Email,
		user.Password,
	).Scan(&userId)
	if err != nil {
		return nil, errors.New("error creating user")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	return &userId, nil
}
