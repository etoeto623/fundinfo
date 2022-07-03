package dbs

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Url string
	Password string
	User string
	Port string
	Database string
}

func OpenDefaultDB() (*sql.DB, error) {
	config := DBConfig{
		Url: "127.0.0.1",
		User: "root",
		Password: "xxxxxx",
		Port: "3306",
		Database: "common",
	}
	return OpenDB(&config)
}

func OpenDB(cfg *DBConfig) (*sql.DB, error) {
	port := cfg.Port
	if len(port) <= 0 {
		port = "3306"
	}
	// "root:1991623@tcp(127.0.0.1:3308)/meeting"
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.User, cfg.Password, cfg.Url, port, cfg.Database)

	db, _ := sql.Open("mysql", url)
	db.SetConnMaxLifetime(2000)
	db.SetMaxIdleConns(1)
	if !isDbConnected(db) {
		return nil, errors.New("dbs connect fail")
	}
	return db, nil
}

func isDbConnected(db *sql.DB) bool{
	if nil != db {
		if err := db.Ping(); nil == err {
			return true
		}else{
			fmt.Println("dbs connect test fail:", err.Error())
		}
	}
	return false
}