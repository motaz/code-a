package controller

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"codea/model"
	"codea/types"
	"codea/util"

	"github.com/motaz/codeutils"
)

func getVersion() string {
	return "1.2.0"
}

func Login(w http.ResponseWriter, r *http.Request) {
	var data types.Home
	if !util.IsConfigFileExist() {
		http.Redirect(w, r, "Setup", http.StatusFound)
		return
	}
	success, _ := model.ThereIsUser()

	if success {
		if r.FormValue("submitlogin") != "" {
			domainID, err := strconv.Atoi(r.FormValue("domain"))
			if err != nil {
				data.ResponseMessage = "Invalid input"
				data.AlertType = "alert-danger"
			} else {
				login := r.FormValue("login")
				password := r.FormValue("password")
				domaininfo, err := model.GetDomainInfo(domainID)
				if err != nil {
					data.ResponseMessage = "Invalid input"
					data.AlertType = "alert-danger"
				}
				now := time.Now()
				ip := GetRemoteAdd(r)
				key := util.GetMD5(r.Header.Get("user-agent") +
					now.String() + "==7B")
				result := doLogin(domaininfo.DomainName, login, password, ip, key, 24, "")
				if result.ErrorCode == 401 || result.ErrorCode == 400 {
					util.WriteErrorLog("error in  doLogin " + result.Message)

				}
				var userInfo types.UserInfo
				if result.Success {

					userInfo, _ = model.GetUserInfo(result.Id)

					if (userInfo != types.UserInfo{}) {
						successfulLogin(w, r, userInfo, result, key)

						page := r.FormValue("page")
						if page != "" {
							if page == "Authenticate" {
								page += "?key=" + r.FormValue("key") +
									"&returnto=" + r.FormValue("returnto")
							}
							http.Redirect(w, r, page, http.StatusFound)

						} else {
							if userInfo.Isadmin == 1 {
								http.Redirect(w, r, "Home", http.StatusFound)

							} else {
								http.Redirect(w, r, domaininfo.DefaultPage, http.StatusFound)

							}
						}
					} else {
						util.WriteErrorLog("Error: " + err.Error())
					}
				} else {
					data.AlertType = "alert-danger"
					data.ResponseMessage = result.Message

				}
			}

		}
	} else {
		http.Redirect(w, r, "AddAdmin", http.StatusFound)
	}
	data.Domains, _ = model.GetDomains(false)

	err := tpl.ExecuteTemplate(w, "login.html", data)
	if err != nil {
		util.WriteErrorLog("error in login template: " + err.Error())
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	RemoveCookie(w, r, "userid")
	RemoveCookie(w, r, "spices")
	RemoveCookie(w, r, "user")
	RemoveCookie(w, r, "codea-sessionid")
	RemoveCookie(w, r, "codea-key")
	http.Redirect(w, r, "Login", http.StatusFound)
}

func doLogin(domainName, userLogin, userPassword, ip, key string, hoursToLive int, sessionInfo string) (loginRes types.Login) {
	if (userLogin == "") || (userPassword == "") {
		loginRes = setError("Empty username or password", 400)
		return
	}

	invalidLoginIPCount, err := model.GetInvalidLoginIPCount(ip)
	if err != nil {
		util.WriteErrorLog("error in GetInvalidLoginIPCount" + err.Error())
	}

	if invalidLoginIPCount > 30 {
		loginRes = setError("IP has been blocked", 401)
		return
	}

	invalidLoginCount, _ := model.GetInvalidLoginCount(userLogin, domainName)
	if invalidLoginCount > 10 {
		loginRes = setError("Username has been blocked", 401)
		return
	}

	userID := 0

	var domainInfo types.DomainType
	if domainName == "" {
		domainInfo, _ = model.GetDefaultDomain()
	} else {
		domainInfo, _ = model.GetDomainInfoByName(domainName)
	}

	if (domainInfo == types.DomainType{}) {
		loginRes = setError("Domain not accessible", 401)
		return
	}

	loginRes.Domain = domainInfo.DomainName
	if domainInfo.IsLocal {
		loginRes = localAuthentication(domainInfo, userLogin, userPassword, loginRes, domainName, ip, hoursToLive, key, sessionInfo)
	} else {
		loginRes = remoteAuthentication(userLogin, userPassword, domainInfo, loginRes, domainName, ip, hoursToLive, key, userID, sessionInfo)
	}

	return
}

func setError(errorMessage string, errorCode int) types.Login {
	result := types.Login{}
	result.Success = false
	result.ErrorCode = errorCode
	result.Message = errorMessage
	return result
}

func localAuthentication(domaininfo types.DomainType, username, password string, loginRes types.Login, domain, ip string, hoursTolive int, key, sessionInfo string) types.Login {
	var success bool
	var userID int
	op, _ := model.CheckUser(domaininfo.DomainName, username, password)

	loginRes.ErrorCode = op.ErrorCode
	loginRes.Id = op.Id
	loginRes.Success = op.Success
	loginRes.Message = op.Message
	success = op.Success
	userID = op.Id
	if !op.Success {
		_, err := LogInvalidAuthentication(op.Message, username, domain, ip)
		if err != nil {
			op.Message = err.Error()
			return loginRes
		}
	}

	loginRes.SessionID = newSession(success, hoursTolive, loginRes, key, userID, domaininfo, username, sessionInfo)
	return loginRes
}

func newSession(success bool, hoursToLive int, loginRes types.Login, key string, userID int, domaininfo types.DomainType, username, sessionInfo string) string {
	if success && hoursToLive > 0 {
		loginRes.SessionID = generateNewSession(key, userID, domaininfo.DomainName, username, hoursToLive, sessionInfo)
		return loginRes.SessionID
	}
	return ""
}

func LogInvalidAuthentication(resultText, username, domain, ip string) (bool, error) {
	obj := make(map[string]interface{})
	obj["result"] = resultText
	sucess, err := model.InsertLog(3, username, domain, ip, obj)
	return sucess, err
}

func generateNewSession(key string, userID int, domainname, username string, hoursToLive int, sessionInfo string) string {
	sessionID := getRandomSessionID(key)
	sessionExpiration := time.Now().Add(time.Duration(int(time.Hour) * hoursToLive)).Format("2006-01-02 15:04:05")
	hashedSessionID := hashSession(key, sessionID)
	success, _ := model.InsertNewSession(hashedSessionID, userID, domainname, username, sessionExpiration, sessionInfo)

	if success {
		return sessionID
	} else {
		return ""
	}
}

func hashSession(key string, sessionID string) string {
	hashedSessionID := codeutils.GetMD5(key + sessionID + "Con;st")
	return hashedSessionID
}

func getRandomSessionID(key string) string {
	leftSide := getLeftPart(key)
	ran := rand.Int()

	randomDigest := util.GetMD5(strconv.Itoa(ran) + "guest")
	sessionID := codeutils.GetMD5(leftSide + randomDigest[:10])
	return sessionID
}

func getLeftPart(key string) string {
	leftSide := codeutils.GetMD5(key + "^%FP")
	leftSide = leftSide[:10]
	return leftSide
}

func remoteAuthentication(username string, password string, domainInfo types.DomainType, loginRes types.Login, domain, ip string, hoursToLive int, key string, userID int, sessionInfo string) types.Login {
	success := false
	req := make(map[string]interface{})
	req["username"] = username
	req["password"] = password

	reqJSON, _ := json.Marshal(req)

	result, err := util.CallURL(domainInfo.RemoteURL, reqJSON)
	if err != nil {
		util.WriteErrorLog("Error calling URL:", err.Error())
		return loginRes
	}
	var codeaRes Response
	err = json.Unmarshal(result, &codeaRes)
	if err != nil {
		util.WriteErrorLog("Error in remoteLogin:", string(result))
		return types.Login{}
	}

	loginRes.Success = codeaRes.Success
	loginRes.ErrorCode = codeaRes.ErrorCode
	loginRes.Message = codeaRes.Message

	if !codeaRes.Success {
		LogInvalidAuthentication(loginRes.Message, username, domain, ip)

	} else {
		loginRes = setError(loginRes.Message, loginRes.ErrorCode)
		loginRes.SessionID = newSession(success, hoursToLive, loginRes, key, userID, domainInfo, username, sessionInfo)
	}

	return loginRes
}

func successfulLogin(w http.ResponseWriter, r *http.Request, userinfo types.UserInfo, result types.Login, key string) {
	SetCookieValue(w, "userid", strconv.Itoa(userinfo.Userid))
	spices := GetSpices(r, result.Id, userinfo)
	SetCookieValue(w, "spices", spices)
	SetCookieValue(w, "user", r.FormValue("login"))
	SetCookieValue(w, "codea-sessionid", result.SessionID)
	SetCookieValue(w, "codea-key", key)
}
