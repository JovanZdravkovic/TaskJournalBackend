package db

import (
	"context"
	"log"

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

func (dbService *DatabaseService) GetTasks() {
	rows, err := dbService.pool.Query(context.Background(), "SELECT * FROM TASK")
	if err != nil {
		log.Fatal("Error while getting tasks from database")
	} else {
		for rows.Next() {
			values, err := rows.Values()
			if err != nil {
				log.Fatal("Error while iterating dataset")
			}

			task_icon := values[2].(string)
			log.Println("task_icon: ", task_icon)
		}
	}
}
