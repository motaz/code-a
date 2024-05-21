package controller

import (
	"code-a/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type SessionInfoType struct {
	UserID int `json:"userid"`
}

type LoginRequest struct {
	Key         string          `json:"key"`
	Username    string          `json:"username"`
	Password    string          `json:"password"`
	Domain      string          `json:"domain"`
	Hours       any             `json:"hours"`
	SessionInfo SessionInfoType `json:"sessioninfo"`
}

type Response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	ErrorCode int    `json:"errorcode"`
}

type LoginResponse struct {
	Response
	Domain    string `json:"domain"`
	Sessionid string `json:"sessionid"`
}

func CheckLoginAPI(w http.ResponseWriter, r *http.Request) {

	setJSONHeader(w)
	var req LoginRequest
	var res LoginResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success, res.ErrorCode, res.Message = false, 400, "invalid request json"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if strings.TrimSpace(req.Username) == "" || strings.TrimSpace(req.Password) == "" {
		res.Success, res.ErrorCode, res.Message = false, 400, "missing parameters: username, password"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if strings.TrimSpace(req.Domain) == "" {
		req.Domain = "Default"
	}

	if req.Key == "" {
		req.Key = r.Header.Get("USER-AGENT")
	}
	var hours int = 4
	if req.Hours != nil {
		hoursStr := fmt.Sprintf("%v", req.Hours)
		hoursInt, _ := strconv.Atoi(hoursStr)
		if hoursInt > 0 {
			hours = hoursInt
		}
	}

	ip := GetRemoteAdd(r)

	loginres := doLogin(req.Domain, req.Username, req.Password, ip, req.Key, hours, req.SessionInfo)
	if loginres.Success {
		res.Success, res.Message, res.ErrorCode = true, "successful login", 0
		res.Sessionid = loginres.SessionID
		res.Domain = loginres.Domain
	} else {
		res.Success, res.Message, res.ErrorCode = false, loginres.Message, loginres.ErrorCode
		w.WriteHeader(http.StatusUnauthorized)
	}
	json.NewEncoder(w).Encode(res)
}

type SessionRequest struct {
	Sessionid string `json:"sessionid"`
	Key       string `json:"key"`
}

type SessionResponse struct {
	Response
	Sessioninfo string `json:"sessioninfo"`
	Domain      string `json:"domain"`
	Username    string `json:"username"`
}

func CheckSessionAPI(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	var req SessionRequest
	var res SessionResponse

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success, res.ErrorCode, res.Message = false, 400, "invalid request json"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if strings.TrimSpace(req.Sessionid) == "" {
		res.Success, res.ErrorCode, res.Message = false, 400, "missing parameters: sessionid"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if req.Key == "" {
		req.Key = r.Header.Get("USER-AGENT")
	}

	res, _ = validateSession(req.Sessionid, req.Key)
	if !res.Success {
		w.WriteHeader(res.ErrorCode)
	}
	json.NewEncoder(w).Encode(res)
}

type RemoveSessionRequest struct {
	All       string `json:"all"`
	Key       string `json:"key"`
	Sessionid string `json:"sessionid"`
}

func RemoveSessionAPI(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	var req RemoveSessionRequest
	var res Response

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		res.Success, res.ErrorCode, res.Message = false, 400, "invalid request json"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	req.Sessionid = strings.TrimSpace(req.Sessionid)
	req.All = strings.TrimSpace(req.All)

	if req.Sessionid == "" {
		res.Success, res.ErrorCode, res.Message = false, 400, "missing parameters: sessionid"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	if req.Key == "" {
		req.Key = r.Header.Get("USER-AGENT")
	}

	sessionRes, userID := validateSession(req.Sessionid, req.Key)
	if !sessionRes.Success {
		res.Success, res.ErrorCode, res.Message = false, sessionRes.ErrorCode, sessionRes.Message
		w.WriteHeader(sessionRes.ErrorCode)
		json.NewEncoder(w).Encode(res)
		return
	}

	var success bool
	if req.All == "yes" {
		success = model.DeleteSessionByUserID(userID)
	} else {
		success = model.DeleteSessionByID(hashSession(req.Key, req.Sessionid))
	}

	if success {
		res = Response{Success: true, Message: "", ErrorCode: 0}
	} else {
		res = Response{Success: false, Message: "error while removing session", ErrorCode: 500}
		if req.All == "yes" {
			res.Message += "s"
		}
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(res)
}

func validateSession(sessionid, key string) (res SessionResponse, userID int) {
	hashedSession := hashSession(key, sessionid)
	sessionObj, err := model.CheckSession(hashedSession)
	if err != nil {
		res.Success, res.ErrorCode, res.Message = false, 500, "database access error"
		return
	}

	if sessionObj.SessionID == "" {
		res.Success, res.ErrorCode, res.Message = false, 401, "invalid session"
		return
	}

	res.Success, res.ErrorCode, res.Message = true, 0, "valid"
	res.Sessioninfo = sessionObj.SessionInfo
	res.Username = sessionObj.Username
	res.Domain = sessionObj.DomainName
	userID = sessionObj.UserID
	return
}

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
