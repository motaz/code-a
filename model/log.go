package model

import (
	"code-a/util"
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

func GetInvalidLoginIPCount(ip string) (count int, err error) {
	count, err = getInvalidCount(ip, "", "ip = ?")
	return count, err
}

func GetInvalidLoginCount(username, domain string) (count int, err error) {
	count, err = getInvalidCount(username, domain, "username = ? and domain = ?")
	if err != nil {
		util.WriteErrorLog("error in getInvalidCount " + err.Error())
		return
	}

	return count, err
}

func getInvalidCount(fieldValue, domain, sqlcondition string) (count int, err error) {
	query := `SELECT count(*) from log where typeID = 3 and logTime > ? and ` + sqlcondition
	var rows *sql.Row
	timeFormat := time.Now().Add(-time.Minute * 10).Format("2006-01-02 15:04:05")
	if domain != "" {
		rows = db.QueryRow(query, timeFormat, strings.ToLower(fieldValue), strings.ToLower(domain))
	} else {
		rows = db.QueryRow(query, timeFormat, strings.ToLower(fieldValue))
	}

	err = rows.Scan(&count)
	if err != nil {
		util.WriteErrorLog("error in getInvalidCount " + err.Error())
		return
	}

	return

}

func truncateLog() bool {
	statement := "truncate log"
	_, err := db.Exec(statement)
	if err != nil {
		util.WriteErrorLog("error in truncateLog " + err.Error())
		return false
	}
	return true
}

func InsertLog(typeID int, username, domain, ip string, details map[string]interface{}) (bool, error) {
	query := `insert into log (typeID, userName, domain, ip, details, LogTime) values (?, ?, ?, ?, ?, now())`
	detailsobj, _ := json.Marshal(details)
	_, err := db.Exec(query, typeID, username, domain, ip, detailsobj)
	if err != nil {
		util.WriteErrorLog("Error in InsertLog: " + err.Error())
	}
	return err == nil, err

}
