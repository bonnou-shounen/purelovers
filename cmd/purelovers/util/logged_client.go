package util

import (
	"context"
	"fmt"
	"os"

	libnetrc "github.com/jdxcode/netrc"

	"github.com/bonnou-shounen/purelovers"
)

func NewLoggedClient(ctx context.Context) (*purelovers.Client, error) {
	id, password := getCredential()
	if id == "" || password == "" {
		return nil, fmt.Errorf("missing credentials")
	}

	client := purelovers.NewClient()
	if err := client.Login(ctx, id, password); err != nil {
		return nil, fmt.Errorf("on NewClient(): %w", err)
	}

	return client, nil
}

func getCredential() (string, string) {
	var id, password string

	getters := []func() (string, string){
		fromEnv,
		fromNetrc,
	}
	for _, getter := range getters {
		if id != "" && password != "" {
			break
		}

		i, p := getter()

		if id == "" {
			id = i
		}

		if password == "" {
			password = p
		}
	}

	return id, password
}

func fromEnv() (string, string) {
	id := os.Getenv("PURELOVERS_LOGIN")
	password := os.Getenv("PURELOVERS_PASSWORD")

	return id, password
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

	id := machine.Get("login")
	password := machine.Get("password")

	return id, password
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
