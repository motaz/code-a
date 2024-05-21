package model

import (
	"code-a/types"
	"code-a/util"
	"database/sql"
)

func InsertNewSession(sessionID string, userID int, domainname string, username string, sessionExpiration string, sessionInfo string) (bool, error) {

	stmt := "insert into sessions (sessionID, userID, username, domainname, sessionTime, sessionExpiration, sessionInfo) values (?, ?, ?, ?, now(), ?, ?)"
	if sessionInfo == "" {
		sessionInfo = "{}"
	}
	_, err := db.Exec(stmt, sessionID, userID, username, domainname, sessionExpiration, sessionInfo)
	if err != nil {
		util.WriteErrorLog("Error in InsertNewSession: " + err.Error())
		return false, err
	} else {

		return true, err
	}
}

func CheckSession(sessionID string) (types.SessionResult, error) {
	stmt := `SELECT id, sessionID, sessionTime, sessionExpiration, userID, username, domainname, sessionInfo 
	FROM sessions where sessionID = ? and sessionExpiration > now()`
	row := db.QueryRow(stmt, sessionID)
	var result types.SessionResult
	err := row.Scan(&result.Id, &result.SessionID, &result.SessionTime, &result.SessionExpiration,
		&result.UserID, &result.Username, &result.DomainName, &result.SessionInfo)
	if err != nil {
		util.WriteErrorLog("error in CheckSession " + err.Error())
		if err == sql.ErrNoRows {
			err = nil
		}
		return types.SessionResult{}, err
	}
	return result, err

}

func DeleteSessionByUsername(username string, domainname string) (bool, error) {
	statement := "DELETE FROM sessions WHERE username = ?"
	_, err := db.Exec(statement, username+"."+domainname)
	if err != nil {
		util.WriteErrorLog("Error in deleteSessionByUsername: %s", err.Error())
		return false, err
	}
	return true, err
}

func removeOldSessions() bool {
	stmt := "DELETE FROM sessions WHERE sessionExpiration < now()"
	_, err := db.Exec(stmt)

	if err != nil {
		return false
	}
	return err == nil
}

func DeleteSessionByID(sessionID string) bool {
	stmt := "DELETE FROM sessions where sessionID = ?"

	_, err := db.Exec(stmt, sessionID)
	if err != nil {
		return false
	}
	return err == nil
}

func DeleteSessionByUserID(userID any) bool {
	stmt := "DELETE FROM sessions where userID = ?"

	_, err := db.Exec(stmt, userID)
	if err != nil {
		return false
	}
	return err == nil
}
