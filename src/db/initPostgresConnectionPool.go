package db

import (
	"context"
	"fmt"
	"github.com/dmitry-kostin/go-rest/src/pkg"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InitPostgresConnectionPool(config *pkg.Config, logger *pkg.Logger) *pgxpool.Pool {
	uri := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", config.DatabaseUser, config.DatabasePassword, config.DatabaseHostname, config.DatabaseName)
	conn, err := pgxpool.Connect(context.Background(), uri)
	if err != nil {
		logger.Fatal().Err(err).Msg("[database] unable to connect to database")
	}
	err = conn.Ping(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("[database] database not available")
	}
	logger.Info().Msg("[database] database connected")
	return conn
}
