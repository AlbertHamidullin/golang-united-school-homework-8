package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	warningItemAlreadyExistsTemplate = "Item with id %s already exists"
	warningItemNotFoundTemplate      = "Item with id %s not found"
	errorFlagNotSpecifiedTemplate    = "-%s flag has to be specified"
	errorIncorrectOperationTemplate  = "Operation %s not allowed!"
)

type Arguments map[string]string

func (args Arguments) GetValue(name string) (string, error) {
	value := args[name]
	if value == "" {
		return "", fmt.Errorf(errorFlagNotSpecifiedTemplate, name)
	}
	return value, nil
}

func AddUser(fileName string, item string, writer io.Writer) error {
	var users Users

	err := users.LoadFromFile(fileName)
	if err != nil {
		return err
	}

	var user User

	err = user.Set(item)
	if err != nil {
		return err
	}

	err = users.Add(user)
	if err != nil {
		//if err.Error() == errorUsersUserAlreadyExists {
		var errS UsersUserAlreadyExistsError
		if errors.As(err, &errS) {
			writer.Write([]byte(fmt.Sprintf(warningItemAlreadyExistsTemplate, user.Id)))
			return nil
		}
		return err
	}

	err = users.SaveToFile(fileName)
	if err != nil {
		return err
	}

	return nil
}

func List(fileName string, writer io.Writer) error {
	var users Users

	err := users.LoadFromFile(fileName)
	if err != nil {
		return err
	}

	data, err := users.AsJsonBytes()
	if err != nil {
		return err
	}

	writer.Write(data)
	return nil
}

func GetUser(fileName, id string, writer io.Writer) error {
	var users Users

	err := users.LoadFromFile(fileName)
	if err != nil {
		return err
	}

	user, err := users.GetById(id)
	if err != nil {
		var errS UsersUserNotFoundError
		if errors.As(err, &errS) {
			writer.Write([]byte(""))
			return nil
		}
		return err
	}

	data, err := user.AsJsonBytes()
	if err != nil {
		return err
	}

	writer.Write(data)
	return nil
}

func RemoveUser(fileName, id string, writer io.Writer) error {
	var users Users

	err := users.LoadFromFile(fileName)
	if err != nil {
		return err
	}

	err = users.RemoveById(id)
	if err != nil {
		var errS UsersUserNotFoundError
		if errors.As(err, &errS) {
			writer.Write([]byte(fmt.Sprintf(warningItemNotFoundTemplate, id)))
			return nil
		}
		return err
	}

	err = users.SaveToFile(fileName)
	if err != nil {
		return err
	}

	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	fileName, err := args.GetValue("fileName")
	if err != nil {
		return err
	}
	operation, err := args.GetValue("operation")
	if err != nil {
		return err
	}

	switch operation {
	case "add":
		item, err := args.GetValue("item")
		if err != nil {
			return err
		}

		err = AddUser(fileName, item, writer)
		if err != nil {
			return err
		}
	case "list":
		err := List(fileName, writer)
		if err != nil {
			return err
		}
	case "findById":
		id, err := args.GetValue("id")
		if err != nil {
			return err
		}

		err = GetUser(fileName, id, writer)
		if err != nil {
			return err
		}
	case "remove":
		id, err := args.GetValue("id")
		if err != nil {
			return err
		}

		err = RemoveUser(fileName, id, writer)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf(errorIncorrectOperationTemplate, operation)
	}
	return nil
}

func parseArgs() Arguments {
	var flagOperation string
	var flagItem User
	var flagId string
	var flagFileName string

	flag.StringVar(&flagOperation, "operation", "", "allowed values: add, list, findById, remove")
	flag.Var(&flagItem, "item", "value json format: {\"id\":\"value\",\"email\":\"value\",\"age\":value}")
	flag.StringVar(&flagId, "id", "", "allowed values: >0")
	flag.StringVar(&flagFileName, "fileName", "", "file name witch content are json format: [{\"id\":\"value\",\"email\":\"value\",\"age\":value}, ...]")

	flag.Parse()

	var flagItemAsBytes []byte
	if flagItem == EmptyUser() {
		flagItemAsBytes = []byte(nil)
	} else {
		var err error
		flagItemAsBytes, err = flagItem.AsJsonBytes()
		if err != nil {
			panic(err)
		}
	}

	return Arguments{"operation": flagOperation, "item": string(flagItemAsBytes), "id": flagId, "fileName": flagFileName}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
