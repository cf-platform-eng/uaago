package sample_cli

import (
	"fmt"
	"os"

	"github.com/cf-platform-eng/uaago"
	"github.com/cf-platform-eng/uaago/management"
)

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	// constructre *real* version of things and kick off flow
	if len(args[1:]) != 3 {
		fmt.Fprintf(os.Stderr, "Usage %s [URL] [USERNAME] [PASS]", args[0])
		return 1
	}

	uaa, err := uaago.NewClient(args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client: %s", err.Error())
		return 1
	}

	token, err := uaa.GetAuthToken(args[2], args[3], false)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Faild to get auth token: %s", err.Error())
		return 1
	}

	fmt.Fprintf(os.Stdout, "TOKEN: %s\n", token)
	return 0
}

// for testing, fake interfaces
func RegisterUsers(registrar management.UaaRegistrar, users []User) error {
	for _, user := range users {
		id, err := registrar.RegisterUser(user.Username, user.Password)
		if err != nil {
			return err
		}
		for _, group := range user.Groups {
			err = registrar.AddUserToGroup(id, group)
			if err != nil {
				return err
			}
		}
	}
	return nil
}