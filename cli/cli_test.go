package cli_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cf-platform-eng/uaago/management"
	. "github.com/cf-platform-eng/uaago/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("uaa_registrar", func() {
	inputJson := `{
		"username": "someuser",
		"password": "somepass",
		"groups": ["group1", "group2"]
	}`

	It("properly de-serializes incoming json to struct", func() {
		user := User{}
		err := json.Unmarshal([]byte(inputJson), &user)
		//fmt.Print("%+v", err)
		Expect(err).To(BeNil())
		Expect(user.Username).To(Equal("someuser"))
		Expect(user.Password).To(Equal("somepass"))
		Expect(user.Groups).To(Equal([]string{"group1", "group2"}))
	})

	Context("given two users", func() {
		users := []User{
			{Username: "Jim", Password: "secret", Groups: []string{"group1", "group2"}},
			{Username: "Jane", Password: "shhh", Groups: []string{"group2", "group3"}},
		}

		It("registers each user", func() {
			providedUsernames := []string{}
			registrar := &MockRegistrar{
				RegisterUserFn: func(uaaUser string, uaaPassword string) (string, error) {
					providedUsernames = append(providedUsernames, uaaUser)
					return fmt.Sprintf("%v-id"), nil
				},
			}
			err := RegisterUsers(registrar, users)

			Expect(err).To(BeNil())
			Expect(providedUsernames).To(Equal([]string{"Jim", "Jane"}))
		})

		It("adds user to groups", func() {
			providedGroupAssignments := map[string][]string{
				"Jim-id":  {},
				"Jane-id": {},
			}
			registrar := &MockRegistrar{
				RegisterUserFn: func(uaaUser string, uaaPassword string) (string, error) {
					return fmt.Sprintf("%v-id", uaaUser), nil
				},
				AddUserToGroupFn: func(userId string, groupName string) error {
					providedGroupAssignments[userId] = append(providedGroupAssignments[userId], groupName)
					return nil
				},
			}
			err := RegisterUsers(registrar, users)

			Expect(err).To(BeNil())
			Expect(providedGroupAssignments["Jim-id"]).To(Equal([]string{"group1", "group2"}))
			Expect(providedGroupAssignments["Jane-id"]).To(Equal([]string{"group2", "group3"}))
		})

		It("returns error if RegisterUser errors", func() {
			registrar := &MockRegistrar{
				RegisterUserFn: func(uaaUser string, uaaPassword string) (string, error) {
					return "", errors.New("asploded")
				},
			}
			err := RegisterUsers(registrar, users)
			Expect(err).NotTo(BeNil())
		})

		It("returns error if AddUserToGroup errors", func() {
			registrar := &MockRegistrar{
				AddUserToGroupFn: func(uaaId string, groupName string) error {
					return errors.New("asploded")
				},
			}
			err := RegisterUsers(registrar, users)
			Expect(err).NotTo(BeNil())
		})
	})

})

type MockRegistrar struct {
	RegisterClientFn func(uaaSecret string, client *management.Client) error
	RegisterUserFn   func(uaaUser string, uaaPassword string) (string, error)
	AddUserToGroupFn func(userId string, groupName string) error
}

func (m *MockRegistrar) RegisterClient(uaaSecret string, client *management.Client) error {
	if m.RegisterClientFn != nil {
		return m.RegisterClientFn(uaaSecret, client)
	}
	return nil
}

func (m *MockRegistrar) RegisterUser(uaaUser string, uaaPassword string) (string, error) {
	if m.RegisterUserFn != nil {
		return m.RegisterUserFn(uaaUser, uaaPassword)
	}
	return "id", nil
}

func (m *MockRegistrar) AddUserToGroup(userId string, groupName string) error {
	if m.AddUserToGroupFn != nil {
		return m.AddUserToGroupFn(userId, groupName)
	}
	return nil
}
