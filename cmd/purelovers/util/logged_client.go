package util

import (
	"context"
	"fmt"
	"os"

	libnetrc "github.com/jdxcode/netrc"

	"github.com/bonnou-shounen/purelovers"
)

func NewLoggedClient(ctx context.Context) (*purelovers.Client, error) {
	loginID, password := getCredential()
	if loginID == "" || password == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	client := purelovers.NewClient()
	if err := client.Login(ctx, loginID, password); err != nil {
		return nil, fmt.Errorf("on NewClient(): %w", err)
	}

	return client, nil
}

func getCredential() (string, string) {
	var loginID, password string

	getters := []func() (string, string){
		fromEnv,
		fromNetrc,
	}
	for _, getter := range getters {
		if loginID != "" && password != "" {
			break
		}

		id, pwd := getter()

		if loginID == "" {
			loginID = id
		}

		if password == "" {
			password = pwd
		}
	}

	return loginID, password
}

func fromEnv() (string, string) {
	loginID := os.Getenv("PURELOVERS_LOGIN")
	password := os.Getenv("PURELOVERS_PASSWORD")

	return loginID, password
}

func fromNetrc() (string, string) {
	netrc := getNetrc()
	if netrc == nil {
		return "", ""
	}

	machine := netrc.Machine("purelovers.com")
	if machine == nil {
		return "", ""
	}

	loginID := machine.Get("login")
	password := machine.Get("password")

	return loginID, password
}

func getNetrc() *libnetrc.Netrc {
	netrcPath := getNetrcPath()
	if netrcPath == "" {
		return nil
	}

	netrc, err := libnetrc.Parse(netrcPath)
	if err != nil {
		return nil
	}

	return netrc
}

func getNetrcPath() string {
	path := os.Getenv("NETRC")
	if path != "" {
		return path
	}

	path = os.Getenv("CURLOPT_NETRC_FILE")

	if path != "" {
		return path
	}

	if home := os.Getenv("HOME"); home != "" {
		return home + "/.netrc"
	}

	return ""
}
