package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/cf-platform-eng/uaago/cli"
	"github.com/cf-platform-eng/uaago/management"
	"github.com/cf-platform-eng/uaago/uaago"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	// constructre *real* version of things and kick off flow
	if len(args[1:]) != 4 {
		fmt.Fprintf(os.Stderr, "Usage %s [URL] [USERNAME] [PASS] [PATH TO USER FILE]\n", args[0])
		return 1
	}

	file, err := ioutil.ReadFile(args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "File error: %s\n", err.Error())
		return 1
	}

	users := []cli.User{}
	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling json file: %s\n", err.Error())
		return 1
	}

	tokenFetcher := uaago.NewUAATokenFetcher(args[1], args[2], args[3], true)

	logger := lager.NewLogger("")
	registrar, err := management.NewUaaRegistrar(args[1], tokenFetcher, true, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new UAA registrar: %s\n", err.Error())
		return 1
	}

	err = cli.RegisterUsers(registrar, users)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error registering users: %s\n", err.Error())
		return 1
	}

	return 0
}
