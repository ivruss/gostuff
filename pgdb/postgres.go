package pgdb

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DB struct {
	Log  *zap.Logger
	Conn *sqlx.DB
}

func NewDB(
	user string,
	password string,
	host string,
	port int,
	dbname string,
	log *zap.Logger,
) (*DB, error) {
	dbConnStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)
	db, err := sqlx.Connect("pgx", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	log.Info("connected to postgreSQL database")
	return &DB{Conn: db, Log: log}, nil
}

func (s *DB) Close() error {
	s.Log.Info("closing postgresql database connection")
	err := s.Conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}
	return nil
}
