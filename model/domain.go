package model

import (
	"codea/types"
	"codea/util"
	"database/sql"
	"errors"
)

func InsertNewDomain(domain types.DomainType) (id int64, err error) {

	query := `INSERT INTO domains
                (domainName, isLocal, defaultDomain, remoteURL, defaultPage, isEnabled)
                VALUES (?, ?, ?, ?, ?, 1)`

	res, err := db.Exec(query, domain.DomainName, domain.IsLocal, domain.DefaultDomain,
		domain.RemoteURL, domain.DefaultPage)
	if err != nil {
		util.WriteErrorLog("Error in insertNewDomain: " + err.Error())
		return
	} else {

		id, err = res.LastInsertId()
	}
	return
}

func GetDomainInfo(domainID int) (domain types.DomainType, err error) {

	query := `select * from domains where domainID = ?`
	rows := db.QueryRow(query, domainID)
	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err = rows.Scan(&domain.DomainID, &domain.DomainName, &domain.IsLocal,
		&domain.DefaultDomain, &domain.IsEnabled, &remoteURL, &defaultPage)
	if remoteURL.Valid {
		domain.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		domain.DefaultPage = defaultPage.String
	}
	if err != nil {
		util.WriteErrorLog("error in GetDomainInfo " + err.Error())
		return
	}
	return
}

func GetDefaultDomain() (domain types.DomainType, err error) {

	query := `select * from domains where defaultDomain `
	rows := db.QueryRow(query)
	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err = rows.Scan(&domain.DomainID, &domain.DomainName, &domain.IsLocal,
		&domain.DefaultDomain, &domain.IsEnabled, &remoteURL, &defaultPage)
	if remoteURL.Valid {
		domain.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		domain.DefaultPage = defaultPage.String
	}
	if err != nil {
		util.WriteErrorLog("error in GetDefaultDomain " + err.Error())
	}
	return

}

func GetDomainInfoByName(domainName string) (domain types.DomainType, err error) {
	query := `select * from domains where Lower(domainName) = ? `
	rows := db.QueryRow(query, domainName)
	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err = rows.Scan(&domain.DomainID, &domain.DomainName, &domain.IsLocal,
		&domain.DefaultDomain, &domain.IsEnabled, &remoteURL, &defaultPage)
	if remoteURL.Valid {
		domain.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		domain.DefaultPage = defaultPage.String
	}
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("domain not found")
			util.WriteErrorLog("error in GetDomainInfoByName: domain " + domainName + " not found")
		} else {
			util.WriteErrorLog("error in GetDomainInfoByName: " + err.Error())
			err = errors.New("system error")
		}
	}
	return

}

func ReadDomainInfo(rows *sql.Rows) (domain types.DomainType, err error) {

	var remoteURL sql.NullString
	var defaultPage sql.NullString
	err = rows.Scan(&domain.DomainID, &domain.DomainName, &domain.IsLocal,
		&domain.DefaultDomain, &domain.IsEnabled, &remoteURL, &defaultPage)
	if remoteURL.Valid {
		domain.RemoteURL = remoteURL.String
	}
	if defaultPage.Valid {
		domain.DefaultPage = defaultPage.String
	}
	if err != nil {
		util.WriteErrorLog("error in ReadDomainInfo " + err.Error())
	}
	return
}

func GetDomains(displayRemote bool) (domains []types.DomainType, err error) {

	query := "select * from domains "
	if !displayRemote {
		query += "where isLocal = 1 "
	}
	query += "order by domainID"

	rows, err := db.Query(query)
	if err != nil {
		util.WriteErrorLog("Error in getDomains: " + err.Error())

	}
	defer rows.Close()

	for rows.Next() {
		var domain types.DomainType
		domain, err = ReadDomainInfo(rows)
		if err != nil {
			util.WriteErrorLog("Error in getDomains, Scan:" + err.Error())
			continue
		}
		domains = append(domains, domain)
	}

	return
}

func DeleteDomain(domainId int) (bool, error) {
	statement := "DELETE FROM domains WHERE domainID = ?"
	_, err := db.Exec(statement, domainId)
	if err != nil {
		util.WriteErrorLog("Error in DeleteDomain: %s", err.Error())
		return false, err
	}
	return true, err
}

func UpdateDomain(domain types.DomainType) (id int64, err error) {

	sql := "update domains set domainName=?, isLocal=?, isEnabled =?, remoteURL=?, defaultPage=? where domainID = ?;"
	result, err := db.Exec(sql, domain.DomainName, domain.IsLocal, domain.IsEnabled, domain.RemoteURL, domain.DefaultPage, domain.DomainID)
	if err != nil {
		util.WriteErrorLog("Error in UpdateDomain: " + err.Error())
		return 0, err
	}
	id, err = result.LastInsertId()
	return
}
