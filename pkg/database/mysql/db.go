package mysqldb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	cfg "mediator/internal/config"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//  *****************************************************************************************************
var (
	createDb = fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s`, cfg.Params.DBName)
)

func migrationRun(db *sql.DB) {
	if _, err := os.Stat("./migrations"); os.IsNotExist(err) {
		log.Println("ERROR: Folder migrations does not exist. Skipping migrations.")
		return
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Println("ERROR: migration ", err)
		return
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		cfg.Params.DBName,
		driver,
	)
	if err != nil {
		log.Println("ERROR: migration ", err)
		return
	}
	m.Steps(3)

}

func DBconnect() (*sql.DB, error) {
	conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/", cfg.Params.DBLogin, cfg.Params.DBPass, cfg.Params.DBHost, cfg.Params.DBPort)

	ctxBG := context.Background()
	ctxConnTimeout, cancel := context.WithTimeout(ctxBG, 2*time.Second)
	defer cancel()
	conn, err := sql.Open("mysql", conn_str)
	if err != nil {
		log.Println("ERROR: Connect to database is fail ", err)
		return nil, err
	}
	_, err = conn.Conn(ctxConnTimeout)
	if err != nil {
		log.Println("ERROR: Database connect timeout ", err)
	}
	conn_err := conn.Ping()
	if conn_err != nil {
		log.Println("ERROR: Database connection check is fail", conn_err)
		return nil, conn_err
	}
	_, err = conn.Exec(createDb)
	if err != nil {
		log.Println("ERROR: Create database is fail", err)
		return nil, err
	}
	conn.Close()
	conn_str = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?multiStatements=true", cfg.Params.DBLogin, cfg.Params.DBPass, cfg.Params.DBHost, cfg.Params.DBPort, cfg.Params.DBName)
	conn, err = sql.Open("mysql", conn_str)
	if err != nil {
		log.Println("ERROR: Connect to database is fail ", err)
		return nil, err
	}
	return conn, err

}

func InitDB(Connection *sql.DB) {
	migrationRun(Connection)
}
