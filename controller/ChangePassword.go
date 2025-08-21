package controller

import (
	"code-a/model"
	"code-a/util"
	"net/http"
	"strconv"
)

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	home := setHeader(w, r, "myadmin")
	// fmt.Println("Home User:", home.User)
	if r.FormValue("resetpassword") != "" {
		oldPassword := r.FormValue("oldpassword")
		newPassword := r.FormValue("newpassword")
		confirmPassword := r.FormValue("confirmpassword")
		if oldPassword == "" || newPassword == "" || confirmPassword == "" {
			home.ResponseMessage = "empty passwords"
			home.AlertType = "alert-danger"

		} else if newPassword != confirmPassword {
			home.ResponseMessage = "mismatch passwords"
			home.AlertType = "alert-danger"

		} else if home.User.Password != util.GetMD5(oldPassword) {
			// fmt.Println("old password:", home.User.Password, "new password:", util.GetMD5(oldPassword))
			home.ResponseMessage = "old password is wrong"
			home.AlertType = "alert-danger"

		} else {
			op, _ := model.CheckUser(home.User.DomainName, home.User.Login, oldPassword)
			if op.Success {
				err := model.UserChangePassword(strconv.Itoa(home.UserID), newPassword)
				if err != nil {
					home.ResponseMessage = "Unable to change password: " + err.Error()
					home.AlertType = "alert-danger"

				} else {
					home.AlertType = "alert-success"
					home.ResponseMessage = "password has been changed successfully"
					model.DeleteSessionByUsername(home.User.Login, home.User.DomainName)
				}
			}
		}
	}

	err := tpl.ExecuteTemplate(w, "ChangePassword.html", home)
	if err != nil {
		util.WriteErrorLog("error in ChangePassword template: " + err.Error())
	}
}
