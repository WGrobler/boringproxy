package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/boringproxy/boringproxy"
	"github.com/joho/godotenv"
)

const usage = `Usage: %s [command] [flags]

Commands:
    version      Prints version information.
    server       Start a new server.
    client       Connect to a server.

Use "%[1]s command -h" for a list of flags for the command.
`

var Version string

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func loadEnvFile(filePath string) {
	// loads values from .env into the system
	var err error
	if filePath != "" {
		err = godotenv.Load(filePath)
		log.Println(fmt.Sprintf("Loading .env file from '%s'", filePath))
	} else {
		err = godotenv.Load()
		log.Println("Loading .env file from working directory")
	}
	if err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	var env_file string
	var flags []string
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, os.Args[0]+": Need a command")
		fmt.Printf(usage, os.Args[0])
		os.Exit(1)
	} else {
		flags = os.Args[2:]
		if len(os.Args) > 2 {
			if os.Args[2] == "config" {
				env_file = os.Args[3]
				flags = os.Args[4:]
			}
		}
	}
	loadEnvFile(env_file)

	command := os.Args[1]
	switch command {
	case "version":
		fmt.Println(Version)
	case "help", "-h", "--help", "-help":
		fmt.Printf(usage, os.Args[0])
	case "server":
		config := boringproxy.SetServerConfig(flags)
		boringproxy.Listen(config)
	case "client":
		config := boringproxy.SetClientConfig(flags)

		ctx := context.Background()

		client, err := boringproxy.NewClient(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = client.Run(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

	default:
		fail(os.Args[0] + ": Invalid command " + command)
	}
}
