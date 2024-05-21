package controller

import (
	"code-a/model"
	"code-a/types"
	"code-a/util"
	"net/http"
	"strconv"
	"strings"
)

func Domains(w http.ResponseWriter, r *http.Request) {
	var err error
	home := setHeader(w, r, "domains")

	if home.IsAdmin {
		if r.FormValue("add") != "" {
			var domain types.DomainType
			domain.DomainName = strings.TrimSpace(r.FormValue("domain"))
			domain.DefaultPage = r.FormValue("defaultpage")
			domain.RemoteURL = r.FormValue("remoteurl")
			domain.IsLocal, err = strconv.ParseBool(r.FormValue("islocal"))
			if err != nil {
				util.WriteErrorLog(err.Error(), " in bool parse")
			}

			_, err = model.InsertNewDomain(domain)
			if err == nil {
				home.AlertType = "alert-success"
				home.ResponseMessage = "Domain (" + domain.DomainName + ") has been added"
			} else {
				home.AlertType = "alert-danger"
				home.ResponseMessage = err.Error()

			}
		}
		if r.FormValue("update") != "" {
			var domain types.DomainType
			domain.DomainName = strings.TrimSpace(r.FormValue("domain"))
			domain.DefaultPage = r.FormValue("defaultpage")
			domain.RemoteURL = r.FormValue("remoteurl")
			domain.IsLocal, err = strconv.ParseBool(r.FormValue("islocal"))
			if err != nil {
				util.WriteErrorLog(err.Error(), " in bool parse")
			}
			domain.IsEnabled, err = strconv.ParseBool(r.FormValue("isenabled"))
			if err != nil {
				util.WriteErrorLog(err.Error(), " in bool parse")
			}

			domain.DomainID, _ = strconv.Atoi(r.FormValue("domainid"))
			_, err = model.UpdateDomain(domain)
			if err == nil {
				home.AlertType = "alert-success"
				home.ResponseMessage = "New Domain (" + domain.DomainName + ") has been updated"
			} else {
				home.AlertType = "alert-danger"
				home.ResponseMessage = err.Error()

			}
		}
		if r.FormValue("remove") != "" {
			if r.FormValue("removedomain") != "" && r.FormValue("removedomain") == "1" {
				domainId, _ := strconv.Atoi(r.FormValue("domainid"))
				success, err := model.DeleteDomain(domainId)
				if err != nil {
					util.WriteErrorLog("error in DeleteDomain " + err.Error())
				}
				if success {
					home.AlertType = "alert-success"
					home.ResponseMessage = "Domain has been deleted"
				} else {
					home.AlertType = "alert-danger"
					home.ResponseMessage = err.Error()
				}
			}
		}
		home.View = r.FormValue("view")
		home.Domains, _ = model.GetDomains(true)

		if home.View == "edit" {
			domainid, _ := strconv.Atoi(r.FormValue("domainid"))
			home.DomainInfo, _ = model.GetDomainInfo(domainid)

			if (home.DomainInfo == types.DomainType{}) {
				http.Redirect(w, r, "Domains", http.StatusFound)

			}
		}
		err = tpl.ExecuteTemplate(w, "Domains.html", home)
		if err != nil {
			util.WriteErrorLog("error in Domains template: " + err.Error())
		}
	} else {
		http.Redirect(w, r, "Login", http.StatusFound)
	}
}
