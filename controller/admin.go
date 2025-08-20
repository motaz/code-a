package controller

import (
	"code-a/model"
	"code-a/types"
	"code-a/util"
	"fmt"
	"net/http"
)

func AddAdmin(w http.ResponseWriter, r *http.Request) {

	configExists := util.IsConfigFileExist()
	added := false
	version := getVersion()
	var setup types.SetUp
	setup.Version = version

	if configExists {
		thereIs, _ := model.ThereIsUser()
		if thereIs {
			http.Redirect(w, r, AppPath+"/Login", http.StatusTemporaryRedirect)
			return

		} else {
			added = doAddAdmin(w, r)
		}

		if added {
			http.Redirect(w, r, AppPath+"/Login", http.StatusTemporaryRedirect)
		} else {
			err := tpl.ExecuteTemplate(w, "AddAdmin.html", setup)
			if err != nil {
				util.WriteErrorLog("error in setup template: " + err.Error())
			}
		}

	} else {
		http.Redirect(w, r, AppPath+"/Setup", http.StatusTemporaryRedirect)

	}

}
func doAddAdmin(w http.ResponseWriter, r *http.Request) bool {
	result := false
	if r.FormValue("set") != "" {
		login := r.FormValue("login")
		password := r.FormValue("pass")
		confirmpassword := r.FormValue("confirmpass")
		fullname := r.FormValue("fullname")
		if login == "" || password == "" || fullname == "" || (password != confirmpassword) {
			fmt.Fprint(w, "<p id=errormessage>Invalid input</p>")
		} else {
			domains, err := model.GetDomains(true)
			if err == nil {
				if len(domains) == 0 {
					var domain types.DomainType
					domain.DomainName = "Default"
					domain.IsLocal = true
					domain.DefaultDomain = true
					domain.RemoteURL = ""
					domain.DefaultPage = ""

					model.InsertNewDomain(domain)

				}
			}
			domainsdata, _ := model.GetDomains(true)

			id := 0
			if len(domainsdata) > 0 {
				s := domains[0]
				id = s.DomainID
			}
			if id != 0 {
				var user types.Users
				user.Login = login
				user.Password = password
				user.Fullname = fullname
				user.Info = ""
				user.Isadmin = 1
				user.DomainID = id

				domaininfo, _ := model.GetDomainInfo(user.DomainID)

				id := model.GetUserID(user.Login, domaininfo.DomainID)

				if id == 0 {
					_, err = model.InsertNewUser(user)
					if err == nil {
						result = true
					}
				}

			}
		}
	}
	return result
}
