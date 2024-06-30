package model

import (
	"code-a/types"
	"code-a/util"
	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func InsertNewUser(user types.Users) (id int, err error) {

	query := `INSERT INTO users
                (login, password, FullName, isEnabled, info, isAdmin, domainID)
                VALUES (?, md5(?), ?, 1, ?, ?, ?)`

	res, err := db.Exec(query, user.Login, user.Password, user.Fullname,
		user.Info, user.Isadmin, user.DomainID)
	if err != nil {
		util.WriteErrorLog("Error in InsertNewUser: " + err.Error())
		return
	} else {
		var userid int64
		userid, err = res.LastInsertId()
		id = int(userid)
		return
	}
}

func GetUserID(login string, DomainID int) (id int) {

	query := `select userid from users where lower(login) = ? and domainID = ?`
	err := db.QueryRow(query, login, DomainID).Scan(&id)
	if err != nil {
		util.WriteErrorLog("error in GetUserID" + err.Error())
		return 0
	}
	return

}

func ThereIsUser() (success bool, err error) {

	query := "select * from users"
	var users []types.Users
	rows, err := db.Query(query)
	if err != nil {
		util.WriteErrorLog("Error in getDomains: " + err.Error())
		return false, err

	}
	defer rows.Close()

	for rows.Next() {
		var user types.Users
		err = rows.Scan(&user.Userid, &user.Login, &user.Password,
			&user.Fullname, &user.Info, &user.Isenabled, &user.Isadmin, &user.DomainID)
		users = append(users, user)
	}
	if len(users) > 0 {
		return true, err
	}

	return
}

func GetUserInfo(userid int) (user types.UserInfo, err error) {

	query := `SELECT userID, login, fullname, info, users.isEnabled, isAdmin, 
		users.domainID, domainName FROM CodeA.users
		INNER JOIN CodeA.domains ON domains.domainid = users.domainid ` +
		`where userID = ? `

	rows := db.QueryRow(query, userid)
	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err = rows.Scan(&user.Userid, &user.Login, &user.Fullname, &user.Info,
		&user.IsEnabled, &user.Isadmin, &user.DomainID, &user.DomainName)

	if remoteURL.Valid {
		user.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		user.DefaultPage = defaultPage.String
	}
	if err != nil {
		util.WriteErrorLog(" error in getUserInfo " + err.Error())
		return
	}
	return
}

func CheckUser(domain, login, password string) (types.Login, error) {
	var operation types.Login

	stmt := `select * from users inner join domains on domains.domainid = users.domainid
	where login = ? and lower(domainname) = ? and users.isEnabled = 1 and (domains.isEnabled=1 or isAdmin=1)`

	rows := db.QueryRow(stmt, login, strings.ToLower(domain))

	var user types.UserInfo

	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err := rows.Scan(&user.Userid, &user.Login, &user.Password,
		&user.Fullname, &user.Info, &user.IsEnabled, &user.Isadmin, &user.DomainID, &user.DomainID, &user.DomainName, &user.IsLocal,
		&user.DefaultDomain, &user.IsEnabled, &remoteURL, &defaultPage)
	if remoteURL.Valid {
		user.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		user.DefaultPage = defaultPage.String
	}

	md5password := util.GetMD5(password)
	if user.Password == md5password {
		operation.ErrorCode = 0
		operation.Success = true
		operation.Message = "Successful login"
		operation.Domain = user.RemoteURL
		operation.Id = user.Userid
	} else {
		operation.ErrorCode = 2
		operation.Success = false
		operation.Message = "Invalid password"
	}
	if err != nil {
		util.WriteErrorLog("error in CheckUser " + err.Error())
		operation.Success = false
		operation.ErrorCode = 1
		operation.Message = "Invalid login"
		return operation, err
	}
	return operation, err
}

func UserChangePassword(userID, password string) error {

	query := `UPDATE users SET password = md5(?) WHERE userid = ?;`
	_, err := db.Exec(query, password, userID)
	if err != nil {
		util.WriteErrorLog("Error in UserChangePassword:", err.Error())
		return err
	}
	return err
}

func GetAllUsers(showall bool) ([]types.UserInfo, error) {

	stmt := ""
	if !showall {
		stmt = ` where users.isenabled =1 `
	}
	query := `SELECT userID, login, fullname, info, users.isEnabled, isAdmin, 
		users.domainID, domainName FROM CodeA.users
		INNER JOIN CodeA.domains ON domains.domainid = users.domainid ` +
		stmt +
		` order by userid`

	rows, err := db.Query(query)
	if err != nil {
		util.WriteErrorLog("Error in getAllUsers: %v", err.Error())
		return nil, err
	}
	defer rows.Close()

	users := []types.UserInfo{}
	for rows.Next() {
		var user types.UserInfo
		var info sql.NullString
		err = rows.Scan(&user.Userid, &user.Login, &user.Fullname, &info, &user.IsEnabled,
			&user.Isadmin, &user.DomainID, &user.DomainName)

		if info.Valid {
			user.Info = info.String
		}
		if err != nil {
			util.WriteErrorLog("error in getAllUsers " + err.Error())
		}
		users = append(users, user)
	}

	return users, nil
}

func SearchUsers(searchText string, domainid, usertype int) ([]types.UserInfo, error) {
	var domid string = ""
	var usrtype string = ""
	var err error
	var rows *sql.Rows
	if usertype == 2 {
		usertype = 0
		usrtype = "and (users.isadmin = ?)"

	} else if usertype == 1 {
		usrtype = "and (users.isadmin = ?)"
	}
	if domainid == 0 {
		domid = "1"
	}

	query := `SELECT userID, login, fullname, info, users.isEnabled, isAdmin, 
		 users.domainID, domainName FROM users
		 INNER JOIN domains ON domains.domainid = users.domainid 
		WHERE (login = ? or login like '%` + searchText + `%'` +
		`OR fullName LIKE '%` + searchText + `%') and( users.domainID = ? or ? )  ` +
		usrtype
	if usrtype != "" {
		rows, err = db.Query(query, searchText, domainid, domid, usertype)
		if err != nil {
			util.WriteErrorLog("error in SearchUsers 1" + err.Error())
			return nil, err
		}
	} else {
		rows, err = db.Query(query, searchText, domainid, domid)
		if err != nil {
			util.WriteErrorLog("error in SearchUsers 2" + err.Error())
			return nil, err
		}
	}
	defer rows.Close()
	users := []types.UserInfo{}
	for rows.Next() {
		var user types.UserInfo
		var remoteURL sql.NullString
		var defaultPage sql.NullString
		var info sql.NullString
		err = rows.Scan(&user.Userid, &user.Login,
			&user.Fullname, &info, &user.IsEnabled, &user.Isadmin, &user.DomainID,
			&user.DomainName)
		if remoteURL.Valid {
			user.RemoteURL = remoteURL.String
		}
		if defaultPage.Valid {
			user.DefaultPage = defaultPage.String
		}
		if info.Valid {
			user.Info = info.String
		}
		if err != nil {
			util.WriteErrorLog("error in getAllUsers " + err.Error())
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func ModifyUserInfo(userinfo types.UserInfo) (id int64, err error) {

	sql := "update users set login = ?, fullname = ?, isEnabled = ?, info = ?, isadmin =?, domainID = ? where userid = ?;"
	result, err := db.Exec(sql, userinfo.Login, userinfo.Fullname, userinfo.IsEnabled, userinfo.Info, userinfo.Isadmin, userinfo.DomainID, userinfo.Userid)
	if err != nil {
		util.WriteErrorLog("Error in ModifyUserInfo: " + err.Error())
		return 0, err
	}
	id, err = result.LastInsertId()
	return
}

func ResetPassword(login, password string, userid int) (success bool) {
	if userid > 0 {

		sql := "update users set password = md5(?) where userid = ?"
		_, err := db.Exec(sql, password, userid)
		if err != nil {
			util.WriteErrorLog("Error in ResetPassword: " + err.Error())
			return false
		}
		return true
	} else {
		return false
	}
}
