package runner

import (
	"github.com/ideal-rucksack/workflow-glolang-plugin/cmd/command"
	"log"
	"os"
)

func Run() {
	register := command.Registry

	if len(os.Args) < 2 {
		log.Fatalln("No command provided")
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	log.Printf("Running command: %s", cmd)

	result, err := register.RunCommand(cmd, args)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}

	log.Printf("Result: %s", result)
}
