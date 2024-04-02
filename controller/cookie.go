package controller

import (
	"codea/model"
	"codea/types"
	"codea/util"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func SetCookieValue(w http.ResponseWriter, name, value string) {
	expiration := time.Now().Add(time.Hour * 24)
	coo := &http.Cookie{Name: name, Value: value, Expires: expiration}

	http.SetCookie(w, coo)

}

func GetCookieValue(r *http.Request, name string) string {
	coo, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	return coo.Value
}

func RemoveCookie(w http.ResponseWriter, r *http.Request, name string) {
	coo, err := r.Cookie(name)
	if err != nil {
		return
	}
	coo.MaxAge = -1
	http.SetCookie(w, coo)
}

func CheckSession(r *http.Request, userid int) bool {
	spices := GetCookieValue(r, "spices")
	util.WriteLog("checkSession: userID: "+strconv.Itoa(userid)+", spices: "+spices, "ca")
	user, _ := model.GetUserInfo(userid)

	if (user == types.UserInfo{}) {
		return false
	}

	currentSpices := GetSpices(r, user.Userid, user)
	util.WriteLog("current: "+currentSpices+", spices: "+spices, "ca")
	return currentSpices == spices
}

func GetSpices(request *http.Request, userID int, userInfo types.UserInfo) string {

	var currentSpices string
	defer func() {
		if r := recover(); r != nil {
			util.WriteLog("getSpices error:")
		}
	}()
	remoteAddress := GetRemoteAdd(request)
	userAgent := request.Header.Get("user-agent")
	util.WriteLog("Check session IP: " + request.RemoteAddr + " cda")
	if remoteAddress == "[::0]" || remoteAddress == "[::]" || strings.Contains(remoteAddress, "0:0:0:0:0:0") {
		remoteAddress = "127.0.0.1"
	}

	currentSpices = util.GetMD5(util.GetMD5(userAgent+strconv.Itoa(userID)+"k980v"+userInfo.Password+remoteAddress+"--") + "!!!")
	return currentSpices
}

func GetRemoteAdd(r *http.Request) string {
	remoteAddress := r.Header.Get("X-REAL-IP")
	if remoteAddress == "" {
		remoteAddress = r.Header.Get("X-FORWARDED-FOR")
	}
	if remoteAddress == "" {
		remoteAddress = "127.0.0.1"
	}
	return remoteAddress
}

func isSessionValid(r *http.Request) bool {
	var result types.OperationResult
	var err error
	var userID int
	useridStr := GetCookieValue(r, "userid")
	if useridStr != "" {
		userID, err = strconv.Atoi(useridStr)
		if err != nil {
			util.WriteErrorLog("error in checkuser ", err.Error())
			return false
		}
	}

	result.Success = CheckSession(r, userID)
	result.ID = userID
	return result.Success
}
func CheckUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isSessionValid(r) {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "Login", http.StatusFound)
	}
}

func setHeader(w http.ResponseWriter, r *http.Request, tabName string) types.Home {
	var home types.Home
	userID := GetCookieValue(r, "userid")
	home.UserID, _ = strconv.Atoi(userID)
	home.User, _ = model.GetUserInfo(home.UserID)
	home.Page = tabName
	home.IsAdmin = home.User.Isadmin == 1
	home.Username = GetCookieValue(r, "user")

	return home
}
