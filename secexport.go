package main

import (
	"fmt"
	"log"
	"os"

	"github.com/2rueSid/secexport/secexport"
)

func main() {
	commandType := os.Args[1]

	var command secexport.Command
	switch commandType {
	case string(secexport.CreateType):
		command = secexport.CreateCommand()
	case string(secexport.DeleteType):
		command = secexport.DeleteCommand()
	case string(secexport.RefreshType):
		log.Panic("Not Implemented")
	case string(secexport.RetrieveType):
		command = secexport.RetrieveCommand()
	default:
		log.Fatalln("Not Implemented")
	}

	if len(os.Args) < 3 || os.Args[2] == "help" || os.Args[2] == "-h" {
		fmt.Println(command.Help())
		os.Exit(0)
	}
	err := command.Parse(os.Args[2:])
	if err != nil {
		log.Fatalln(err.Error())
	}

	res, err := command.Execute()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(*res)
}
