package database

import (
	"context"
	"fmt"
	"strconv"

	"courses/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDataBase() (*pgxpool.Pool, error) {
	host := config.DBHost
	port, _ := strconv.Atoi(config.DBPort)
	user := config.DBUser
	dbname := config.DBName
	pass := config.DBPassword

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, pass, host, port, dbname)
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %v\n", err)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("Could not ping database: %v", err)
	}

	return dbpool, nil
}
