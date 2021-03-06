package util

import (
	"errors"
	"os"

	libnetrc "github.com/jdxcode/netrc"

	"github.com/bonnou-shounen/purelovers"
)

func NewLoggedClient() (*purelovers.Client, error) {
	id, password := getCredential()
	if id == "" || password == "" {
		return nil, errors.New("missing credentials")
	}

	client := purelovers.NewClient()
	err := client.Login(id, password)

	return client, err
}

func getCredential() (id, password string) {
	getters := []func() (string, string){
		fromEnv,
		fromNetrc,
	}

	for _, getter := range getters {
		if id != "" && password != "" {
			return
		}

		i, p := getter()

		if id == "" {
			id = i
		}

		if password == "" {
			password = p
		}
	}

	return
}

func fromEnv() (id, password string) {
	id = os.Getenv("PURELOVERS_LOGIN")
	password = os.Getenv("PURELOVERS_PASSWORD")

	return
}

func fromNetrc() (id, password string) {
	netrc := getNetrc()
	if netrc == nil {
		return
	}

	machine := netrc.Machine("www.purelovers.com")
	if machine == nil {
		return
	}

	id = machine.Get("login")
	password = machine.Get("password")

	return
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

	dir := os.Getenv("HOME")
	if dir != "" {
		return dir + "/.netrc"
	}

	return ""
}
