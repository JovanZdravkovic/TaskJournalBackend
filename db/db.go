package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseService struct {
	pool *pgxpool.Pool
}

func NewDatabaseService(dbPool *pgxpool.Pool) *DatabaseService {
	return &DatabaseService{
		pool: dbPool,
	}
}

func (dbService *DatabaseService) GetTasks() ([]TaskDB, error) {
	rows, err := dbService.pool.Query(context.Background(), "SELECT t.* FROM task t WHERE t.starred = false")
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

func (dbService *DatabaseService) GetStarredTasks() ([]TaskDB, error) {
	rows, err := dbService.pool.Query(context.Background(), "SELECT t.* FROM task t WHERE t.starred = true")
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

// func (dbService *DatabaseService) GetTask() (TaskDB, error) {}

// func (dbService *DatabaseService) CreateTask() (uuid.UUID, error) {}

// // TODO: Updating a single task will have multiple functions for different kinds of updates
// func (dbService *DatabaseService) PutTask() (string, error) {}

// func (dbService *DatabaseService) DeleteTask() (string, error) {}
