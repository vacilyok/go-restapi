package rules

import (
	"context"
	"encoding/json"
	"errors"
	"mediator/internal/adapters/db/rulestorage"
	"mediator/internal/config"
	"mediator/pkg/database/pg"

	"github.com/jackc/pgx/v4"
)

type dbs struct {
	conn *pgx.Conn
}

func NewRuleStorage(dbConn *pgx.Conn) rulestorage.RuleStorage {
	return &dbs{
		conn: dbConn,
	}
}

// ****************************************************************************************************************************
func (db *dbs) inactiveAllRules() error {
	ctx := context.Background()
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	deleteOldRules := "DELETE FROM rules WHERE id < (select id from (select max(id)-20 as id from rules) t)"
	_, err = tx.Exec(ctx, deleteOldRules)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE rules SET active = $1", false)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

// ****************************************************************************************************************************
//  save json body to db
// ****************************************************************************************************************************
func (db *dbs) SaveRulesToDB(ruleList map[string]interface{}) error {
	var err error
	if !pg.IsConnected(db.conn) {
		db.conn, err = pg.ConnectDB()
	}
	if err != nil {
		err := errors.New(" Query save rules to db failed. No connection to db")
		config.Logging.Error(err.Error())
		return err
	}
	jsonRules, err := json.Marshal(ruleList)
	if err != nil {
		config.Logging.Error(err.Error())
		return err
	}
	ctx := context.Background()
	tx, err := db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	db.inactiveAllRules()

	_, err = tx.Exec(ctx, "INSERT INTO  rules (body, active,created_at) VALUES ($1,$2,NOW()) ", string(jsonRules), true)
	if err != nil {
		return err
	}
	tx.Commit(ctx)
	return nil
}

// ****************************************************************************************************************************
//  Get  rules json from db
// ****************************************************************************************************************************
func (db *dbs) GetRulesFromDB() (string, error) {
	var (
		rulesBody string
		err       error
	)
	if !pg.IsConnected(db.conn) {
		db.conn, err = pg.ConnectDB()
		if err != nil {
			return "", err
		}
	}
	result, err := db.conn.Query(context.Background(), "SELECT body FROM rules where active=true order by id desc limit 1")
	if err != nil {
		return "", err
	}
	if result.Next() {
		result.Scan(&rulesBody)
	}
	result.Close()
	return rulesBody, nil
}
