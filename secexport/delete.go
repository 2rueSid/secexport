package secexport

import (
	"errors"
	"fmt"
)

// DeleteCommand
//
// Check if record exists for $(pwd) and deletes encrypted string if password is valid.
// If no record exists returns error.
//
// Usage:
//
// command [OPTIONS]
//
// Options:
//
// -p  <password> Password that was used for creating record. required.
type deleteCommand struct {
	Password string
	context  *commandContext
}

func (c *deleteCommand) Execute() (*string, error) {
	data, err := ReadFile()
	if err != nil {
		return nil, err
	}

	decrypted, err := Decrypt(data, c.Password)
	if err != nil {
		return nil, err
	}

	if !IsByteJSON(decrypted) {
		return nil, errors.New("password invalid")
	}

	err = DeleteFile()
	if err != nil {
		return nil, err
	}

	res := "record deleted"
	return &res, nil
}

func (c *deleteCommand) Parse(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if args[0] != "-p" {
		return fmt.Errorf("option %v is unknown", args[0])
	}

	c.Password = args[1]

	return nil
}

func (c *deleteCommand) Help() string {
	return `
Delete Command:

Check if record exists for $(pwd) and deletes encrypted string if password is valid.
If no record exists returns error.

Usage:

command [OPTIONS]

Options:

-p  <password> Password that was used for creating record. required.
  `
}

func DeleteCommand() *deleteCommand {
	return &deleteCommand{
		Password: "",
		context: &commandContext{
			flags:   map[string]bool{"-p": true},
			command: DeleteType,
		},
	}
}
