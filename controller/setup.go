package controller

import (
	"code-a/model"
	"code-a/types"
	"code-a/util"
	"embed"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

var (
	tpl *template.Template
)

func InitTemplate(embededTemplates embed.FS) error {
	var err error
	tpl, err = template.ParseFS(embededTemplates, "view/*.html")
	if err != nil {
		util.WriteErrorLog("error in InitTemplate: " + err.Error())
		return err
	}
	return nil
}

func Setup(w http.ResponseWriter, r *http.Request) {
	var setup types.SetUp

	configExists := util.IsConfigFileExist()
	useridStr := GetCookieValue(r, "userid")
	userID, _ := strconv.Atoi(useridStr)

	if !configExists {
		setupSuccess := setupConfig(w, r)

		if setupSuccess {
			success, _ := model.ThereIsUser()
			if success {
				http.Redirect(w, r, "Login", http.StatusFound)
			} else {
				http.Redirect(w, r, "AddAdmin", http.StatusFound)
			}
		} else {
			success, _ := model.ThereIsUser()

			if configExists && success && !CheckSession(r, userID) {
				fmt.Fprintf(w, "<p id=errormessage>Please login first</p>")
			} else {
				version := getVersion()
				server := util.GetConfigValue("dbserver", "")
				database := util.GetConfigValue("dbname", "")
				dbuser := util.GetConfigValue("dbuser", "")
				dbpass := util.GetConfigValue("dbpass", "")
				if server == "" {
					server = "localhost"
				}
				if database == "" {
					database = "CodeA"
				}
				setup.Database = database
				setup.Server = server
				setup.Version = version
				setup.User = dbuser
				setup.Password = dbpass
				setup.IsConfigExists = configExists
			}
		}
	} else {
		http.Redirect(w, r, "Home", http.StatusFound)
	}

	err := tpl.ExecuteTemplate(w, "setup.html", setup)
	if err != nil {
		util.WriteErrorLog("error in setup template: " + err.Error())
	}

}

func setupConfig(w http.ResponseWriter, r *http.Request) bool {
	result := false
	if r.FormValue("set") != "" {
		result = util.SetConfigValue("server", r.FormValue("server"))
		util.SetConfigValue("database", r.FormValue("database"))
		util.SetConfigValue("dbuser", r.FormValue("dbuser"))
		util.SetConfigValue("dbpass", r.FormValue("dbpassword"))

		model.CloseDB()
		model.InitDB()
	}
	return result
}
