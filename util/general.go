package util

import (
	"fmt"
	"strings"

	"github.com/motaz/codeutils"
)

func WriteLog(text ...string) {

	if GetConfigValue("debug", "") == "yes" {
		fmt.Println(text)
	}
	codeutils.WriteToLog(fmt.Sprint(text), "info")
}

func WriteErrorLog(text ...string) {

	if GetConfigValue("debug", "") == "yes" {
		fmt.Println(text)
	}
	codeutils.WriteToLog(fmt.Sprint(text), "error")
}

func IsConfigFileExist() bool {

	return codeutils.IsFileExists("config.ini")
}

func GetConfigValue(key, defaulValue string) string {

	value := codeutils.GetConfigValue("config.ini", key)
	if value == "" {
		return defaulValue
	}
	return value
}

func SetConfigValue(key, value string) bool {
	return codeutils.SetConfigValue("config.ini", key, value)
}

func GetMD5(text string) string {

	return strings.ToLower(codeutils.GetMD5(text))
}

func CallURL(aurl string, body []byte) (ResultBody []byte, err error) {

	result := codeutils.CallURLAsPost(aurl, body, 10)
	ResultBody = result.Content
	err = result.Err
	return
}
