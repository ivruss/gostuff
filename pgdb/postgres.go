package pgdb

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type DB struct {
	log  *zap.Logger
	conn *sqlx.DB
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
	log.Info("connected to postgreSQL database")
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}
	return &DB{conn: db, log: log}, nil
}

func (s *DB) Close() error {
	s.log.Info("closing postgresql database connection")
	err := s.conn.Close()
	if err != nil {
		return fmt.Errorf("error closing connection: %w", err)
	}
	return nil
}
