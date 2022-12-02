package test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Configuration struct {
	DBLogin      string
	DBPass       string
	DBHost       string
	DBPort       int
	RPCHost      string
	RPCPort      int
	DBName       string
	MediatorName string
	MediatorPort int
	MediatorHost string
}

var (
	uri    string
	Params Configuration
)

func init() {
	var err error
	Params, err = ReadConcfig()
	if err != nil {
		log.Println("unknown configuration. Check the config file")
		panic(err)
	}
	uri = fmt.Sprintf("http://%s:%d", Params.MediatorHost, Params.MediatorPort)

}

// func TestMain(m *testing.M) {
// 	os.Exit(testMainWrapper(m))
// }

func ReadConcfig() (Configuration, error) {
	config := Configuration{}
	var err error
	config_file := "../conf.json"
	_, err = os.Stat(config_file)
	if err == nil {
		file, _ := os.Open(config_file)
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&config)
	}
	config.MediatorHost = "127.0.0.1"
	config.MediatorPort = 4001
	return config, err
}

// func testMainWrapper(m *testing.M) int {
// 	conn_str := fmt.Sprintf("%s:%s@tcp(%s:%d)/?multiStatements=true", Params.DBLogin, Params.DBPass, Params.DBHost, Params.DBPort)
// 	conn, err := sql.Open("mysql", conn_str)
// 	if err != nil {
// 		log.Println("Connect to database is fail ", err)

// 	}
// 	// defer func() {
// 	// 	dropDb := fmt.Sprintf("DROP DATABASE  %s", Params.DBName)
// 	// 	conn.Exec(dropDb)
// 	// 	conn.Close()
// 	// }()

// 	createDb := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", Params.DBName)
// 	_, err = conn.Exec(createDb)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	migrationRun(conn)
// 	return m.Run()

// }

// func migrationRun(db *sql.DB) {
// 	if _, err := os.Stat("../migrations"); os.IsNotExist(err) {
// 		log.Println("Folder migrations does not exist. Skipping migrations.")
// 		return
// 	}
// 	_, err := db.Exec("USE " + Params.DBName)
// 	if err != nil {
// 		log.Println("Use db ", err)
// 	}
// 	driver, err := mysql.WithInstance(db, &mysql.Config{})
// 	if err != nil {
// 		log.Println("==>", err)
// 		return
// 	}
// 	m, err := migrate.NewWithDatabaseInstance(
// 		"file://../migrations/",
// 		Params.DBName,
// 		driver,
// 	)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	m.Steps(3)

// }
