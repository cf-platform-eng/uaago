package cli

import (
	"github.com/cf-platform-eng/uaago/management"
)

type User struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Groups   []string `json:"groups"`
}



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