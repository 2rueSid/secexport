package secexport

import (
	"encoding/json"
	"errors"
	"log"
)

// Create Command:
//
//	Creates an encrypted record for the current directory $(pwd) and returns formatted exports as a string.
//	If secrets already exist for this directory, it returns an error.
//
// Usage:
//
//	command [OPTIONS] <filters>
//
// Parameters:
//
//	<filters>    A set of strings indicating the filters to retrieve secrets from AWS. Required.
//	             Example: tree flower grass
//
// Options:
//
//	-s           Use secrets from the Secrets Manager. Default: true.
//	-pm 0        Use secrets from the Parameter Store. Default: true.
//	-p <password> Password for encrypting and retrieving secrets, or deleting them. Required if encryption is enabled.
//
// Examples:
//
//	command -s -p YourPassword tree flower grass
//
// Notes:
//
//	This command interacts with AWS Secrets Manager and Parameter Store based on the provided options.
type createCommand struct {
	Filters        []*string
	SecretManager  bool
	ParameterStore bool
	Password       string
	context        *commandContext
}

func (c *createCommand) Execute() (*string, error) {
	file, err := CreateFile()
	if err != nil {
		return nil, err
	}

	// Retrieve AWS Data
	data, err := RetreiveSecrets(c.Filters, c.ParameterStore, c.SecretManager)
	if err != nil {
		log.Printf("Got error when retrieving secrets from the AWS. %v", err)
		return nil, err
	}

	jsonData, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}
	// Encrypt
	inp := []byte(jsonData)

	encrypted, err := Encrypt(inp, c.Password)
	if err != nil {
		return nil, nil
	}

	// encoded := base64.StdEncoding.EncodeToString(encrypted)

	// Write to a file
	err = WriteFile(file, encrypted)
	if err != nil {
		return nil, err
	}

	return data.Values(), nil
}

func (c *createCommand) Parse(args []string) error {
	flags := c.context.flags

	for i := 0; i < len(args); i++ {
		if flags[args[i]] {
			if args[i] != "-p" && i+1 < len(args) && args[i+1] == "false" {
				switch args[i] {
				case "-s":
					c.SecretManager = false
				case "-pm":
					c.ParameterStore = false
				}
				i++
			} else if args[i] == "-p" && i+1 < len(args) {
				c.Password = args[i+1]
				i++
			}
		} else {
			c.Filters = append(c.Filters, &args[i])
		}
	}
	if c.Password == "" {
		msg := "can't have password as Nil when Encryption is enabled"
		return errors.New(msg)
	}

	return nil
}

func (c *createCommand) Help() string {
	return `
Create Command:
    Creates an encrypted record for the current directory $(pwd) and returns formatted exports as a string.
    If secrets already exist for this directory, it returns an error.

Usage:
    command [OPTIONS] <filters>

Parameters:
    <filters>    A set of strings indicating the filters to retrieve secrets from AWS. Required.
                 Example: tree flower grass

Options:
    -s           Use secrets from the Secrets Manager. Default: true.
    -pm 0        Use secrets from the Parameter Store. Default: true.
    -p <password> Password for encrypting and retrieving secrets, or deleting them. Required.

Examples:
    command -s -p YourPassword tree flower grass

Notes:
    This command interacts with AWS Secrets Manager and Parameter Store based on the provided options.
  `
}

func CreateCommand() *createCommand {
	return &createCommand{
		Filters:        make([]*string, 0),
		SecretManager:  true,
		ParameterStore: true,
		Password:       "",
		context: &commandContext{
			flags:   map[string]bool{"-p": true, "-s": true, "-pm": true},
			command: CreateType,
		},
	}
}
