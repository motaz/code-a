package controller

import (
	"code-a/model"
	"code-a/types"
	"code-a/util"
	"net/http"
	"strconv"
)

func ListUsers(w http.ResponseWriter, r *http.Request) {

	var err error
	home := setHeader(w, r, "listusers")

	if home.IsAdmin {
		if r.FormValue("resetpassword") != "" {

			if r.FormValue("password") == "" {
				home.AlertType = "alert-danger"
				home.ResponseMessage = "Empty password"
			} else if r.FormValue("password") != r.FormValue("confirmpassword") {
				home.AlertType = "alert-danger"
				home.ResponseMessage = "Passwords do not match"
			} else {
				login := r.FormValue("login")
				userid, _ := strconv.Atoi(r.FormValue("userid"))
				if model.ResetPassword(login, r.FormValue("password"), userid) {
					home.ResponseMessage = "Password has been reset"
					home.AlertType = "alert-success"
				} else {
					home.AlertType = "alert-danger"
					home.ResponseMessage = "Error while resetting password"
				}
			}
		}

		if r.FormValue("update") != "" {
			var userinfo types.UserInfo
			userinfo.Userid, _ = strconv.Atoi(r.FormValue("userid"))
			userinfo.Isadmin = r.FormValue("isadmin") == "1"
			userinfo.IsEnabled = r.FormValue("isenabled") == "1"
			userinfo.DomainID, _ = strconv.Atoi(r.FormValue("domain"))
			userinfo.Login = r.FormValue("login")
			userinfo.Fullname = r.FormValue("fullname")
			userinfo.Info = ""
			_, err = model.ModifyUserInfo(userinfo)

			if err == nil {
				home.AlertType = "alert-success"
				home.ResponseMessage = "Information Updated"
				if !userinfo.IsEnabled {
					deleteSessionByUsername(userinfo.Login, r.FormValue("domain"))
				}
			} else {
				home.AlertType = "alert-danger"
				home.ResponseMessage = err.Error()

			}
		}

		home.Showall = r.FormValue("showall")
		home.Search = r.FormValue("search")
		home.SearchButton = r.FormValue("searchButton")
		home.Modify = r.FormValue("modify")
		DomainID, _ := strconv.Atoi(r.FormValue("domain"))

		UserType, _ := strconv.Atoi(r.FormValue("usertype"))

		if home.Modify != "" {
			userid, _ := strconv.Atoi(r.FormValue("userid"))
			home.User, _ = model.GetUserInfo(userid)
			home.Domains, _ = model.GetDomains(true)
		}
		if home.Search == "" && UserType == 0 && DomainID == 0 {
			showall := r.FormValue("showall") != ""
			home.UserInfo, _ = model.GetAllUsers(showall)

		} else {
			home.DomainID = DomainID
			home.UserType = UserType
			home.UserInfo, _ = model.SearchUsers(home.Search, home.DomainID, home.UserType)

		}

		home.Domains, _ = model.GetDomains(false)

		err = tpl.ExecuteTemplate(w, "ListUsers.html", home)
		if err != nil {
			util.WriteErrorLog("error in ListUsers template: " + err.Error())
		}
	} else {
		http.Redirect(w, r, "Login", http.StatusFound)
	}
}

func deleteSessionByUsername(username, domainname string) types.OperationResult {
	var op types.OperationResult
	var err error
	op.Success, _ = model.DeleteSessionByUsername(username, domainname)
	if !op.Success {
		op.ErrorCode = 5
		op.Message = err.Error()
	} else {
		op.ErrorCode = 0
		op.Message = ""
	}
	return op
}
