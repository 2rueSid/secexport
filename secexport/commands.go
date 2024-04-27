package secexport

// 1. Input a set of strings that identify names
// 2. flag for secrets manager
// 3. flag for parameter store
// 4. create command - create and exports vars for the current dir aka project and saves them locally encrypted
// 5. delete command - deletets secrets from current dir and unset them
// 6. apply check if secrets are created for the current dir and apply them if yes
// 7. refresh compares and writes secrets for the current dir
// 8. bind to tmux to set the keys for current dir
// 9. get secrets from the lambda function
// 10. get secrets from the ecs task def
type CommandType string

const (
	CreateType   CommandType = "create"
	DeleteType   CommandType = "delete"
	RetrieveType CommandType = "retrieve"
	RefreshType  CommandType = "refresh"
)

type Command interface {
	Execute() (string, error)
	Parse(args []string) error
	Help() string
}

type commandContext struct {
	flags   map[string]bool
	command CommandType
}
