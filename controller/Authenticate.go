package controller

import (
	"codea/model"
	"codea/types"
	"codea/util"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {
	var authenticate types.AuthenticateTemplate

	authenticate.Version = getVersion()
	authenticate.Key = r.FormValue("key")
	authenticate.Returnto = r.FormValue("returnto")
	if r.FormValue("authenticate") != "" {
		authenticate.IsAuthenticate = true
		sessionID := GetCookieValue(r, "codea-sessionid")
		currentKey := GetCookieValue(r, "codea-key")
		result := checkSession(sessionID, currentKey)

		if result.Success {
			authenticate.Success = true
			authenticate.NewSessionID = generateNewSession(authenticate.Key, 0, result.DomainName, result.Username, 4, "")
		}
	} else {
		authenticate.IsAuthenticate = false
	}

	err := tpl.ExecuteTemplate(w, "Authenticate.html", authenticate)
	if err != nil {
		util.WriteErrorLog("error in Authenticate template: " + err.Error())
	}
}

func checkSession(sessionID string, currentKey string) types.SessionResult {
	// implementation of checkSession function
	hashedSessionID := hashSession(currentKey, sessionID)
	result, _ := model.CheckSession(hashedSessionID)

	return result

}
