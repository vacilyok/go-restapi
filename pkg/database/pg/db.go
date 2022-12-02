package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mediator/internal/config"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4"
	mg "github.com/jackc/tern/migrate"
)

var (
	dbHost = fmt.Sprintf("postgres://%s:%s@%s:%d/?sslmode=disable", config.Params.DBLogin, config.Params.DBPass, config.Params.DBHost, config.Params.DBPort)
	dbName = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", config.Params.DBLogin, config.Params.DBPass, config.Params.DBHost, config.Params.DBPort, config.Params.DBName)
)

func createDB() error {
	ctx := context.Background()
	ctxConnTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	db, err := sql.Open("postgres", dbHost)
	if err != nil {
		return err
	}

	_, err = db.Conn(ctxConnTimeout)
	if err != nil {
		config.Logging.Error("Database connect timeout")
		return err
	}

	var dbId int
	req_dbname := fmt.Sprintf("SELECT oid FROM pg_database WHERE datname = '%s'", config.Params.DBName)
	db.QueryRow(req_dbname).Scan(&dbId)
	if dbId > 0 {
		return nil
	}
	createDb := fmt.Sprintf("CREATE DATABASE %s", config.Params.DBName)
	_, err = db.Exec(createDb)
	if err != nil {
		return err
	}
	defer db.Close()
	return nil

}

func IsConnected(conn *pgx.Conn) bool {
	if conn == nil {
		return false
	}
	ctx := context.Background()
	ctxConnTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	err := conn.Ping(ctxConnTimeout)
	return err == nil
}

func ConnectDB() (*pgx.Conn, error) {
	ctx := context.Background()
	ctxConnTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctxConnTimeout, dbName)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func DBInit() (*pgx.Conn, error) {
	err := createDB()
	if err != nil {
		return nil, err
	}
	conn, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	err = migrateDatabase(conn)
	if err != nil {
		return nil, err
	}
	return conn, nil

}

func migrateDatabase(conn *pgx.Conn) error {
	ctx := context.Background()
	migrator, err := mg.NewMigrator(ctx, conn, "migrate_version")
	if err != nil {
		err_msg := fmt.Sprintf("Unable to create a migrator: %v", err.Error())
		return errors.New(err_msg)

	}

	err = migrator.LoadMigrations("./migrations")
	if err != nil {
		err_msg := fmt.Sprintf("Unable to load migrations: %v", err.Error())
		return errors.New(err_msg)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		err_msg := fmt.Sprintf("Unable to migrate: %v", err.Error())
		return errors.New(err_msg)
	}
	_, err = migrator.GetCurrentVersion(ctx)
	if err != nil {
		err_msg := fmt.Sprintf("Unable to get current schema version: %v", err.Error())
		return errors.New(err_msg)
	}
	return nil
}
