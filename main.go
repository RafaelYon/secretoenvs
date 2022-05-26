package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s [ENVIRONMENT_NAME]:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	var (
		envName, envValue = retriveSourceEnv()
		values            map[string]string
	)

	if err := json.Unmarshal([]byte(envValue), &values); err != nil {
		fatalPrint(fmt.Sprintf("Can't json unmarshal value from %s: '%s'", envName, envValue))
	}

	for key := range values {
		fmt.Printf("%s=%s\n", key, values[key])
	}
}

func retriveSourceEnv() (string, string) {
	envName := flag.Arg(0)
	if len(envName) < 1 {
		fatalPrint(`The posicional argument "ENVIRONMENT_NAME" must be a non empty string.`)
	}

	envValue := os.Getenv(envName)
	if len(envValue) < 1 {
		fatalPrint(`The positional argument "ENVIRONMENT_NAME" must reference an existing environment variable.`)
	}

	return envName, envValue
}

func fatalPrint(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
