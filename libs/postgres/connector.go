package postgres

import (
	"context"
	"route256/libs/logger"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func CreateDBPool(ctx context.Context, connString string, attemptCount uint16) (*pgxpool.Pool, error) {
	var counts uint16

	for {
		pool, err := pgxpool.Connect(ctx, connString)
		if err != nil {
			logger.Info("postgres not yet ready...")
			counts++
		} else {
			logger.Info("connected to postgres")
			return pool, nil
		}

		if counts > attemptCount {
			return nil, err
		}

		logger.Info("waiting 1 second...")
		time.Sleep(time.Second * 1)
		continue
	}
}
