package controller

import (
	"code-a/model"
	"code-a/types"
	"code-a/util"
	"net/http"
	"strconv"
	"strings"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	home := setHeader(w, r, "adduser")
	if home.IsAdmin {
		insertuser := r.FormValue("insertuser")
		if insertuser != "" {
			home.ResponseStatus = false
			login := strings.TrimSpace(r.FormValue("login"))
			fullname := r.FormValue("fullname")
			password := r.FormValue("password")
			confirmpassword := r.FormValue("confirmpassword")
			home.AlertType = "alert-danger"
			if login == "" {
				home.ResponseMessage = "Empty login"
				home.AlertType = "alert-danger"
			} else if fullname == "" {
				home.ResponseMessage = "Empty user name"
				home.AlertType = "alert-danger"

			} else if password == "" {
				home.ResponseMessage = "Empty password "
				home.AlertType = "alert-danger"

			} else if confirmpassword == "" {
				home.ResponseMessage = "Empty password confirmation"
				home.AlertType = "alert-danger"

			} else if confirmpassword != password {
				home.ResponseMessage = "Passwords do not match"
				home.AlertType = "alert-danger"

			} else {
				domainID, err := strconv.Atoi(r.FormValue("domain"))
				if err != nil {
					home.ResponseMessage = "Invalid Input"
					home.AlertType = "alert-danger"
				} else {
					var user types.Users
					user.Login = login
					user.Password = password
					user.Fullname = fullname
					user.Info = ""
					user.Isadmin = 0
					user.DomainID = domainID

					domaininfo, _ := model.GetDomainInfo(user.DomainID)

					id := model.GetUserID(user.Login, domaininfo.DomainID)

					if id == 0 {
						id, _ = model.InsertNewUser(user)
					}

					if id > 0 {
						home.ResponseMessage = "New user (" + login + ") has been added</p>"
						home.ResponseStatus = true
						home.AlertType = "alert-success"
					} else {
						home.ResponseMessage = "Error while adding new user"
						home.AlertType = "alert-danger"
					}
				}

			}

		}
		if (home.User != types.UserInfo{}) {
			if r.FormValue("domain") != "" {
				home.DomainID, _ = strconv.Atoi(r.FormValue("domain"))
				domain, _ := model.GetDomainInfo(home.DomainID)

				if (domain != types.DomainType{}) {
					if !domain.IsLocal {
						r.Form.Set("hidepassword", "true")
					}
				}
			}

			home.Domains, _ = model.GetDomains(false)

		}

	} else {
		http.Redirect(w, r, "Login", http.StatusFound)
		return
	}

	err := tpl.ExecuteTemplate(w, "adduser.html", home)
	if err != nil {
		util.WriteErrorLog("error in addUser template: " + err.Error())
	}
}
