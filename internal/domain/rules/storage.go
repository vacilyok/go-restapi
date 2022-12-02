package rules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"mediator/internal/adapters/db/rulestorage"
	"mediator/internal/config"
)

type dbs struct {
	conn         *sql.DB
	isConnection bool
}

func NewRuleStorage(dbConn *sql.DB, isConnection bool) rulestorage.RuleStorage {
	return &dbs{
		conn:         dbConn,
		isConnection: isConnection,
	}
}

// ****************************************************************************************************************************
func (db *dbs) inactiveAllRules() error {
	ctx := context.Background()
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, "UPDATE rules SET `active` = ?", 0)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// ****************************************************************************************************************************
//  save json body to db
// ****************************************************************************************************************************
func (db *dbs) SaveRulesToDB(ruleList map[string]interface{}) error {
	if !db.isConnection {
		err := errors.New(" Query save rules to db failed. No connection to db")
		config.Mysqllog.Error(err.Error())
		return err
	}

	jsonRules, err := json.Marshal(ruleList)
	if err != nil {
		config.Mysqllog.Error(err.Error())
		return err
	}
	ctx := context.Background()
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	db.inactiveAllRules()

	_, err = tx.ExecContext(ctx, "INSERT INTO  rules (body, `active`) VALUES (?,?) ", string(jsonRules), 1)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// ****************************************************************************************************************************
//  Get  rules json from db
// ****************************************************************************************************************************
func (db *dbs) GetRulesFromDB() (string, error) {
	var rulesBody string
	// err := db.conn.QueryRow("SELECT body FROM rules where active=1 order by id desc limit 1").Scan(&rulesBody)
	// if err != nil {
	// 	return "", err
	// }
	// return rulesBody, nil

	result, err := db.conn.Query("SELECT body FROM rules where active=1 order by id desc limit 1")
	if err != nil {
		return "", err
	}

	if result.Next() {
		result.Scan(&rulesBody)
	}
	result.Close()
	return rulesBody, nil

}
