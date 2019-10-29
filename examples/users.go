package main

import (
	"fmt"

	grafana "github.com/skilld-labs/grafana-golang"
)

func createUserExample() {
	g, err := grafana.NewBasicAuthClient(nil, "https://yourgrafana.domain", "admin", "admin")
	if err != nil {
		fmt.Println(err)
	}

	opt := &grafana.CreateUserOptions{
		Name:     "Test Testing",
		Email:    "test@testing.test",
		Login:    "ttest",
		Password: "test",
	}

	_, err = g.Users.CreateUser(opt)
	if err != nil {
		fmt.Println(err)
	}
}
