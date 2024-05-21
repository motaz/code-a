package controller

import (
	"code-a/util"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	home := setHeader(w, r, "home")
	home.Url = home.User.DefaultPage
	home.Domain = home.User.DomainName
	err := tpl.ExecuteTemplate(w, "home.html", home)
	if err != nil {
		util.WriteErrorLog("error in home template: " + err.Error())
	}
}
