package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	errorParseUserJsonTemplate  = "parse user json error: %w"
	errorParseUsersJsonTemplate = "parse users json error: %w"
	errorParseValueTemplate     = "parse value error: %w"
	errorOpenFileTemplate       = "open file error: %w"
	errorReadFileTemplate       = "read file error: %w"
	errorWriteTemplate          = "write file error: %w"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

func ParseUserJson(value string) (user *User, err error) {
	errI := json.Unmarshal([]byte(value), user)
	if errI != nil {
		return nil, fmt.Errorf(errorParseUserJsonTemplate, errI)
	}
	return user, nil
}

func ParseUsersJson(value string) (users *[]User, err error) {
	errI := json.Unmarshal([]byte(value), users)
	if errI != nil {
		return nil, fmt.Errorf(errorParseUsersJsonTemplate, errI)
	}
	return users, nil
}

func (user *User) Set(value string) error {
	userI, err := ParseUserJson(value)
	if err != nil {
		return err
	}
	*user = *userI
	return nil
	// err := json.Unmarshal([]byte(value), user)
	// if err != nil {
	// 	return fmt.Errorf(errorParseUserJsonTemplate, err)
	// }
	// return nil
}
func (user *User) String() string {
	return fmt.Sprint(user.Id, user.Email, user.Age)
}
func GetUsersFromFile(fileName string) (users *[]User, err error) {
	file, errI := os.Open(fileName)
	if errI != nil {
		return nil, fmt.Errorf(errorOpenFileTemplate, errI)
	}
	defer file.Close()
	data, errI := io.ReadAll(file)
	if errI != nil && errI != io.EOF {
		return nil, fmt.Errorf(errorReadFileTemplate, errI)
	}
	if len(data) > 0 {
		users, errI = ParseUsersJson(string(data))
		if err != nil {
			return nil, errI
		}
	} else {
		users = new([]User)
	}
	return users, nil
}
func GetUserIndexById(users []User, id string) int {
	for i, u := range users {
		if u.Id == id {
			return i
		}
	}
	return -1
}
func Perform(args Arguments, writer io.Writer) error {
	fileName := args["fileName"]
	if fileName == "" {
		//writer.Write([]byte("-fileName flag has to be specified"))
		return errors.New("-fileName flag has to be specified")
	}
	operation := args["operation"]
	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}

	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
		//var us User
		userI, err := ParseUserJson(item)
		if err != nil {
			return err
		}
		// err := json.Unmarshal([]byte(item), &us)
		// if err != nil && err != io.EOF {
		// 	return fmt.Errorf(errorParseValueTemplate, err)
		// }
		// users, err := GetUsersFromFile(fileName)
		// if err != nil {
		// 	return err
		// }
		// if -1 != GetUserIndexById(*users, userI.Id){
		// 	writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", userI.Id)))
		// 	return nil
		// }
		// users = append(users, *userI)

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf(errorOpenFileTemplate, err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorReadFileTemplate, err)
		}
		var users []User
		if len(data) > 0 {
			err = json.Unmarshal(data, &users)
			if err != nil {
				return fmt.Errorf(errorParseValueTemplate, err)
			}
			for _, u := range users {
				if u.Id == userI.Id {
					writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", userI.Id)))
					return nil
				}
			}
		}
		users = append(users, *userI)
		data, err = json.Marshal(users)
		if err != nil {
			return fmt.Errorf(errorParseValueTemplate, err)
		}
		_, err = file.Write(data)
		if err != nil {
			return fmt.Errorf(errorWriteTemplate, err)
		}
	case "list":
		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf(errorOpenFileTemplate, err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorReadFileTemplate, err)
		}
		writer.Write(data)
	case "findById":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorOpenFileTemplate, err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorReadFileTemplate, err)
		}
		var users []User
		if len(data) > 0 {
			err = json.Unmarshal(data, &users)
			if err != nil {
				return fmt.Errorf(errorParseValueTemplate, err)
			}
		}
		for _, u := range users {
			if u.Id == id {
				data, err = json.Marshal(u)
				if err != nil {
					return fmt.Errorf(errorParseValueTemplate, err)
				}
				writer.Write(data)
				return nil
			}
		}
		writer.Write([]byte(""))
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorOpenFileTemplate, err)
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil && err != io.EOF {
			return fmt.Errorf(errorReadFileTemplate, err)
		}
		var users []User
		if len(data) > 0 {
			err = json.Unmarshal(data, &users)
			if err != nil {
				return fmt.Errorf(errorParseValueTemplate, err)
			}
		}
		for i, u := range users {
			if u.Id == id {
				users[i] = users[len(users)-1]
				users = users[:len(users)-1]
				data, err = json.Marshal(users)
				if err != nil {
					return fmt.Errorf(errorParseValueTemplate, err)
				}
				err = file.Truncate(0)
				if err != nil {
					return fmt.Errorf(errorWriteTemplate, err)
				}
				_, err = file.WriteAt(data, 0)
				if err != nil {
					return fmt.Errorf(errorWriteTemplate, err)
				}
				return nil
			}
		}
		writer.Write([]byte(fmt.Sprintf("Item with id %s not found", id)))
	default:
		return fmt.Errorf("Operation %s not allowed!", operation)
	}
	return nil
}

func parseArgs() Arguments {
	var flagOperation string
	var flagItem User
	var flagId string
	var flagFileName string
	flag.StringVar(&flagOperation, "operation", "", "allowed values: add, list, findById, remove")
	flag.Var(&flagItem, "item", "value json format: {\"id\":\"value\",\"email\":\"value\",\"age\":age}")
	flag.StringVar(&flagId, "id", "", "allowed values: >0")
	flag.StringVar(&flagFileName, "fileName", "", "file name witch content are json format: [{\"id\":\"value\",\"email\":\"value\",\"age\":value}, ...]")
	flag.Parse()
	return Arguments{"operation": flagOperation, "item": flagItem.String(), "id": flagId, "fileName": flagFileName}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
