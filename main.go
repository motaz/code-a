// code-a project main.go
// Central Authentication web service and manager tool for Code projects
// started by Motaz Abdel Azeem 17 Dec 2023
// as migration from CodeA java version

package main

import (
	"code-a/controller"
	"code-a/model"
	"code-a/util"
	"embed"
	"fmt"
	"net/http"
)

const appPath = "/codea"
const VERSION = "1.2.2 r30-June"

//go:embed view
var view embed.FS

//go:embed assets
var assets embed.FS

func main() {
	model.InitDB()
	controller.InitTemplate(view)
	mux := http.NewServeMux()

	util.WriteLog("code-a version: " + VERSION)
	mux.Handle(appPath+"/assets/", http.StripPrefix(appPath, http.FileServer(http.FS(assets))))

	// Portal Pathes
	mux.HandleFunc("/", redirect)
	mux.HandleFunc(appPath+"/", controller.CheckUser(controller.Home))
	mux.HandleFunc(appPath+"/Login", controller.Login)
	mux.HandleFunc(appPath+"/Logout", controller.Logout)
	mux.HandleFunc(appPath+"/Setup", controller.Setup)
	mux.HandleFunc(appPath+"/AddAdmin", controller.CheckUser(controller.AddAdmin))
	mux.HandleFunc(appPath+"/Home", controller.CheckUser(controller.Home))
	mux.HandleFunc(appPath+"/AddUser", controller.CheckUser(controller.AddUser))
	mux.HandleFunc(appPath+"/Authenticate", controller.CheckUser(controller.Authenticate))
	mux.HandleFunc(appPath+"/ChangePassword", controller.CheckUser(controller.ChangePassword))
	mux.HandleFunc(appPath+"/Domains", controller.CheckUser(controller.Domains))
	mux.HandleFunc(appPath+"/ListUsers", controller.CheckUser(controller.ListUsers))

	// Web service pathes
	mux.HandleFunc(appPath+"/CheckLogin", controller.CheckLoginAPI)
	mux.HandleFunc(appPath+"/CheckSession", controller.CheckSessionAPI)
	mux.HandleFunc(appPath+"/RemoveSession", controller.RemoveSessionAPI)
	fmt.Println("http://localhost:2024")
	err := http.ListenAndServe(":2024", mux)
	if err != nil {
		util.WriteLog("Error in listening: " + err.Error())
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, appPath, http.StatusTemporaryRedirect)
}
