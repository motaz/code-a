package main

import (
	"fmt"
	"testing"
)

func Test(a *testing.T) {

	err := InitDB()
	if err == nil {
		domain, err := getDomainInfo(4)
		if err == nil {
			fmt.Printf("domain: %+v\n", domain)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())

	}

}
