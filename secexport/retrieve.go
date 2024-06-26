package secexport

import (
	"errors"
	"fmt"
)

type retrieveCommand struct {
	Password string
	context  *commandContext
}

func (c *retrieveCommand) Execute() (*string, error) {
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

	secrets := AWSSecrets{}

	res, err := secrets.Parse(decrypted)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *retrieveCommand) Parse(args []string) error {
	if len(args) <= 1 {
		return errors.New("not enough arguments")
	}
	if args[0] != "-p" {
		return fmt.Errorf("option %v is unknown", args[0])
	}

	c.Password = args[1]

	return nil
}

func (c *retrieveCommand) Help() string {
	return `
  Retrieve Command:

  Check if record exists for $(pwd) and returns encrypted string if password is valid.
If no record exists returns error.

  Usage:

    command [OPTIONS]

  Options:

    -p  <password> Password that was used for creating record. required.

  Examples:

    refresh -p 123456
  `
}

// RetrieveCommand
//
// Check if record exists for $(pwd) and returns encrypted string if password is valid.
// If no record exists returns error.
//
// Usage:
//
// command [OPTIONS]
//
// Options:
//
// -p  <password> Password that was used for creating record. required.
//
// Examples:
//
// refresh -p 123456
func RetrieveCommand() *retrieveCommand {
	return &retrieveCommand{
		Password: "",
		context: &commandContext{
			flags:   map[string]bool{"-p": true},
			command: RetrieveType,
		},
	}
}
