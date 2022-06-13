package main

import (
	"encoding/json"
	"fmt"
)

const (
	errorUserMarshalJsonTemplate   = "marshal user to json error: %w"
	errorUserUnmarshalJsonTemplate = "unmarshal user from json error: %w"
)

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func EmptyUser() User {
	return User{Id: "", Email: "", Age: 0}
}

func (user *User) Set(value string) error {
	err := json.Unmarshal([]byte(value), user)
	if err != nil {
		return fmt.Errorf(errorUserUnmarshalJsonTemplate, err)
	}
	return nil
}

func (user *User) String() string {
	return fmt.Sprint(user.Id, user.Email, user.Age)
}

func (user User) AsJsonBytes() ([]byte, error) {
	data, err := json.Marshal(user)
	if err != nil {
		return []byte(nil), fmt.Errorf(errorUserMarshalJsonTemplate, err)
	}
	return data, nil
}
