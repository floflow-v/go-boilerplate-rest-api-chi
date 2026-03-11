package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/database/sqlc"
)

type Database struct {
	*sql.DB
	Queries *sqlc.Queries
}

func Init(cfg config.Config, logger zerolog.Logger) (*Database, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.DBPassword,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.MaxLifetimeConn)
	db.SetConnMaxIdleTime(cfg.Database.MaxIdleTimeConn)

	if err := db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Failed to ping database")
		return nil, err
	}

	logger.Info().Msg("Successfully connected to MariaDB")

	return &Database{
		DB:      db,
		Queries: sqlc.New(db),
	}, nil
}

func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
