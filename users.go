package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type Users struct {
	users []User
}

const (
	errorUsersOpenFileTemplate      = "open file with users error: %w"
	errorUsersReadFileTemplate      = "read file with users error: %w"
	errorUsersUnmarshalJsonTemplate = "unmarshal users from json error: %w"
	errorUsersMarshalJsonTemplate   = "marshal users to json error: %w"
	errorUsersWriteFileTemplate     = "write file with users error: %w"
	errorUsersUserNotFound          = "user not found"
	errorUsersUserAlreadyExists     = "user already exists"
)

type UsersUserAlreadyExistsError struct {
	Message string
}

func (err UsersUserAlreadyExistsError) Error() string {
	return err.Message
}

func NewUsersUserAlreadyExistsError() UsersUserAlreadyExistsError {
	return UsersUserAlreadyExistsError{Message: errorUsersUserAlreadyExists}
}

type UsersUserNotFoundError struct {
	Message string
}

func (err UsersUserNotFoundError) Error() string {
	return err.Message
}

func NewUsersUserNotFoundError() UsersUserNotFoundError {
	return UsersUserNotFoundError{Message: errorUsersUserNotFound}
}

func (users *Users) LoadFromFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf(errorUsersOpenFileTemplate, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil && err != io.EOF {
		return fmt.Errorf(errorUsersReadFileTemplate, err)
	}

	users.users = make([]User, 0)
	if len(data) > 0 {
		err = json.Unmarshal(data, &users.users)
		if err != nil {
			return fmt.Errorf(errorUsersUnmarshalJsonTemplate, err)
		}
	}

	return nil
}

func (users Users) SaveToFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf(errorUsersOpenFileTemplate, err)
	}
	defer file.Close()

	data, err := json.Marshal(users.users)
	if err != nil {
		return fmt.Errorf(errorUsersMarshalJsonTemplate, err)
	}

	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf(errorUsersWriteFileTemplate, err)
	}

	_, err = file.WriteAt(data, 0)
	if err != nil {
		return fmt.Errorf(errorUsersWriteFileTemplate, err)
	}

	return nil
}

func (users Users) GetById(id string) (user User, err error) {
	for _, u := range users.users {
		if u.Id == id {
			return u, nil
		}
	}

	return EmptyUser(), NewUsersUserNotFoundError()
}

func (users *Users) RemoveById(id string) error {
	for i, u := range users.users {
		if u.Id == id {
			users.users[i] = users.users[len(users.users)-1]
			users.users = users.users[:len(users.users)-1]
			return nil
		}
	}

	return NewUsersUserNotFoundError()
}

func (users *Users) Add(user User) error {
	_, err := users.GetById(user.Id)
	if err != nil {
		var errS UsersUserNotFoundError
		if errors.As(err, &errS) {
			users.users = append(users.users, user)
			return nil
		}
		return err
	}

	return NewUsersUserAlreadyExistsError()
}

func (users Users) AsJsonBytes() ([]byte, error) {
	data, err := json.Marshal(users.users)
	if err != nil {
		return []byte(nil), fmt.Errorf(errorUsersMarshalJsonTemplate, err)
	}
	return data, nil
}
