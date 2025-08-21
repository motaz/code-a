package main

import (
	"code-a/model"
	"fmt"
	"testing"
)

func Test(a *testing.T) {

	err := model.InitDB()
	if err == nil {
		domain, err := model.GetDomainInfo(4)
		if err == nil {
			fmt.Printf("domain: %+v\n", domain)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())

	}

}
