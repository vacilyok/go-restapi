package logger

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Logger interface {
	// LogRest(string, string, string)
	Error(string)
	Warning(string)
	Info(string)
	LogRest(string, string, string)
}

type LogService struct {
	connection   *sql.DB
	isConnection bool
}

func NewLogger(connection *sql.DB, isConnection bool) Logger {
	return &LogService{
		connection:   connection,
		isConnection: isConnection,
	}
}

//Logs errors message
func (l *LogService) Error(msg string) {
	log.Println("ERROR:", msg)
	if l.isConnection {
		l.insert(msg, 1)
	}
}

func (l *LogService) Warning(msg string) {
	log.Println("WARNING: ", msg)
	if l.isConnection {
		l.insert(msg, 2)
	}
}

func (l *LogService) Info(msg string) {
	log.Println("INFO: ", msg)
	if l.isConnection {
		l.insert(msg, 3)
	}
}

func (l *LogService) insert(msg string, lvlmsg int) {
	cur_time := time.Now().Format("2006-01-02 15:04:05")
	ctx := context.Background()
	err := l.connection.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	tx, err := l.connection.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Logger: ", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO logger (date, source, message, level) VALUES (?,?,?,?)", cur_time, "", msg, lvlmsg)
	if err != nil {
		log.Println("Logger: Insert into logger ", err)
	}
	tx.Commit()
}

// ******************************************************************************************************************************************
//Logs REST API event
func (l *LogService) LogRest(msg string, route string, method string) {
	if !l.isConnection {
		return
	}
	cur_time := time.Now().Format("2006-01-02 15:04:05")
	ctx := context.Background()
	err := l.connection.PingContext(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	if len(msg) == 0 {
		msg = "{}"
	}
	tx, err := l.connection.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Logger: ", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "INSERT INTO rest_logger (date, source, route, method,message) VALUES (?,?,?,?,?)", cur_time, "", route, method, msg)
	if err != nil {
		log.Println("Logger: Insert into rest_logger", err)
	}
	tx.Commit()
}
