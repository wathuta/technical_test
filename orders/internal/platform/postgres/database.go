package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/wathuta/technical_test/orders/internal/config"
)

func PostgreSQLConnection() (*sqlx.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv(config.DBMaxConnectionsEnvVar))
	maxIdleConn, _ := strconv.Atoi(os.Getenv(config.DBMaxIdleConnectionsEnvVar))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv(config.DBMaxLifetimeConnectionsEnvVar))

	// Build PostgreSQL connection URL.
	postgresConnURL, err := config.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}
	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", *postgresConnURL)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings:
	// 	- SetMaxOpenConns: the default is 0 (unlimited)
	// 	- SetMaxIdleConns: defaultMaxIdleConns = 2
	// 	- SetConnMaxLifetime: 0, connections are reused forever
	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn))

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
func OpenDBConnection() (*sqlx.DB, error) {
	// Define Database connection variables.
	var (
		db  *sqlx.DB
		err error
	)

	db, err = PostgreSQLConnection()
	if err != nil {
		return nil, err
	}

	return db, nil
}
