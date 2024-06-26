package model

import (
	"code-a/util"
	"database/sql"
	"fmt"
	"time"
)

var db *sql.DB

func InitDB() (err error) {

	dpip := util.GetConfigValue("dbserver", "localhost")
	dbname := util.GetConfigValue("dbname", "CodeA")
	dbuser := util.GetConfigValue("dbuser", "")
	dbpass := util.GetConfigValue("dbpass", "")
	connectionString := fmt.Sprintf("%v:%v@tcp(%s:3306)/%v?parseTime=false",
		dbuser, dbpass, dpip, dbname)

	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		util.WriteErrorLog("Error in InitDB connection: " + err.Error())
		return
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)
	db.SetConnMaxLifetime(time.Minute * 2)

	go checklog()
	return
}

func CloseDB() (err error) {
	err = db.Close()
	if err != nil {
		util.WriteErrorLog("error in CloseDB " + err.Error())
	}
	return
}

func checklog() {
	for {
		removeOldSessions()
		truncateLog()
		time.Sleep(time.Hour * 24)
	}
}
