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

// Maybe needed in future
var icons [33]string = [33]string{
	"job",
	"doctor_appointment",
	"mechanic",
	"electrician",
	"transport",
	"cleaning",
	"swimming",
	"gym",
	"basketball",
	"football",
	"american_football",
	"volleyball",
	"concert",
	"movie",
	"meeting",
	"reading",
	"writing",
	"payment",
	"message",
	"photography",
	"moving",
	"running",
	"drive",
	"shopping",
	"coffee",
	"sailing",
	"church",
	"pets",
	"plants",
	"lunch",
	"phone_call",
	"computer",
	"party",
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

// TASK

func (dbService *DatabaseService) GetTasks(userId uuid.UUID, searchName *string, searchIcons []string, searchOrderBy *string) ([]TaskDB, error) {
	query := "SELECT t.* FROM task t WHERE t.created_by = @userId AND t.exec_status = 'ACTIVE'"
	if len(searchIcons) > 0 && searchIcons[0] != "null" {
		query += " AND t.task_icon = ANY(@searchIcons::text[])"
	}
	if searchName != nil && *searchName != "null" && *searchName != "" {
		query += " AND t.task_name ILIKE concat('%', @searchName::text, '%')"
	}
	if searchOrderBy != nil {
		if *searchOrderBy == "starred" {
			query += " ORDER BY t.starred DESC"
		} else if *searchOrderBy == "deadline" {
			query += " ORDER BY t.deadline ASC"
		}
	}
	rows, err := dbService.pool.Query(
		context.Background(),
		query,
		pgx.NamedArgs{
			"userId":        userId,
			"searchName":    *searchName,
			"searchIcons":   searchIcons,
			"searchOrderBy": *searchOrderBy,
		},
	)
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
	var task TaskDB
	err := dbService.pool.QueryRow(context.Background(), "SELECT t.* FROM task t WHERE t.id = $1 AND t.created_by = $2", taskId, userId).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("task doesn't exist")
		}
		return nil, errors.New("unexpected error")
	}
	return &task, nil
}

func (dbService *DatabaseService) CreateTask(task TaskPost) (*uuid.UUID, error) {
	var taskId uuid.UUID
	err := dbService.pool.QueryRow(
		context.Background(),
		"INSERT INTO task(task_name, task_icon, task_desc, deadline, starred, exec_status, created_by) VALUES ($1, $2, $3, $4, $5, 'ACTIVE', $6) RETURNING id",
		task.TaskName,
		task.TaskIcon,
		task.TaskDesc,
		task.Deadline,
		task.Starred,
		task.CreatedBy,
	).Scan(&taskId)
	if err != nil {
		return nil, errors.New("error while creating task")
	}
	return &taskId, nil
}

func (dbService *DatabaseService) CompleteTask(taskId uuid.UUID, userId uuid.UUID) (bool, error) {
	tx, err := dbService.pool.BeginTx(context.Background(), pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return false, err
	}
	defer tx.Rollback(context.Background())

	var owner uuid.UUID
	var execStatus string
	err = tx.QueryRow(context.Background(), "SELECT created_by, exec_status FROM task WHERE id = $1", taskId).Scan(&owner, &execStatus)
	if err != nil {
		return false, err
	}

	if (owner != userId) || execStatus != "ACTIVE" {
		return false, errors.New("invalid request")
	}

	_, err = tx.Exec(context.Background(), "UPDATE task SET exec_status = 'INACTIVE' WHERE id = $1", taskId)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO task_history(task_id) VALUES ($1)", taskId)
	if err != nil {
		return false, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return false, err
	}

	return true, nil
}

func (dbService *DatabaseService) UpdateTask(taskId uuid.UUID, task TaskPut, userId uuid.UUID) error {
	cmdTag, err := dbService.pool.Exec(
		context.Background(),
		"UPDATE task SET task_name = $1, task_icon = $2, task_desc = $3, starred = $4, deadline = $5 WHERE id = $6 AND created_by = $7",
		task.TaskName,
		task.TaskIcon,
		task.TaskDesc,
		task.Starred,
		task.Deadline,
		taskId,
		userId,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("task doesn't exist")
	}
	return nil
}

func (dbService *DatabaseService) DeleteTask(taskId uuid.UUID, userId uuid.UUID) error {
	cmdTag, err := dbService.pool.Exec(context.Background(), "DELETE FROM task t WHERE t.id = $1 AND t.created_by = $2", taskId, userId)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("task doesn't exist")
	}
	return nil
}

// TASK HISTORY

func (dbService *DatabaseService) GetTaskHistory(taskHistoryId uuid.UUID, userId uuid.UUID) (*TaskHistoryDB, error) {
	var taskHistory TaskHistoryDB
	err := dbService.pool.QueryRow(context.Background(), "SELECT th.*, t.task_name, t.task_icon FROM task_history th JOIN task t ON th.task_id = t.id WHERE t.created_by = $1 AND th.id = $2", userId, taskHistoryId).
		Scan(
			&taskHistory.Id,
			&taskHistory.ExecRating,
			&taskHistory.ExecComment,
			&taskHistory.TaskId,
			&taskHistory.TaskName,
			&taskHistory.TaskIcon,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("task history doesn't exist")
		}
		return nil, errors.New("unexpected error")
	}
	return &taskHistory, nil
}

func (dbService *DatabaseService) GetTasksHistory(userId uuid.UUID, searchName *string, searchIcons []string, searchRating int) ([]TaskHistoryDB, error) {
	query := "SELECT th.*, t.task_name, t.task_icon FROM task_history th JOIN task t ON th.task_id = t.id WHERE t.created_by = @userId"
	if len(searchIcons) > 0 && searchIcons[0] != "null" {
		query += " AND t.task_icon = ANY(@searchIcons::text[])"
	}
	if searchName != nil && *searchName != "null" && *searchName != "" {
		query += " AND t.task_name ILIKE concat('%', @searchName::text, '%')"
	}
	if searchRating >= 1 && searchRating <= 3 {
		query += " AND th.exec_rating = @searchRating::int"
	}
	rows, err := dbService.pool.Query(
		context.Background(),
		query,
		pgx.NamedArgs{
			"userId":       userId,
			"searchName":   *searchName,
			"searchIcons":  searchIcons,
			"searchRating": searchRating,
		},
	)
	if err != nil {
		return nil, errors.New("error while getting tasks history from database")
	} else {
		var tasksHistory []TaskHistoryDB
		for rows.Next() {
			var taskHistory TaskHistoryDB
			err := rows.Scan(
				&taskHistory.Id,
				&taskHistory.ExecRating,
				&taskHistory.ExecComment,
				&taskHistory.TaskId,
				&taskHistory.TaskName,
				&taskHistory.TaskIcon,
			)
			if err != nil {
				return nil, errors.New("error while iterating dataset")
			}
			tasksHistory = append(tasksHistory, taskHistory)
		}
		return tasksHistory, nil
	}
}

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

// USER

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
