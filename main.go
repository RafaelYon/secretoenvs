package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	flagQuotationMarks    string
	flagKeyPrefix         string
	flagKeyValueSeparator string
	flagRawValue          bool
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s [RAW_JSON or ENVIRONMENT_VARIATION_NAME]:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(
		&flagQuotationMarks,
		"quotation-marks",
		"",
		"Specifies a combination of characters to be placed before and after the ENV value in the output.",
	)

	flag.StringVar(
		&flagKeyPrefix,
		"key-prefix",
		"",
		"Specifies a prefix for each key in the output.",
	)

	flag.StringVar(
		&flagKeyValueSeparator,
		"key-value-separator",
		"=",
		"Specifies the characters to use to separate key values",
	)

	flag.BoolVar(
		&flagRawValue,
		"raw-value",
		false,
		"Specifies that JSON was passed as an argument. So it is not necessary to specify an ENV as input.",
	)

	flag.Parse()
}

func main() {
	var values map[string]string

	if err := json.Unmarshal([]byte(retriveInputJson()), &values); err != nil {
		fatalPrint(fmt.Sprintf("Can't unmarshal specified json : %s", err.Error()))
	}

	for key := range values {
		fmt.Printf(
			"%s%s%s%s%s%s\n",
			flagKeyPrefix,
			key,
			flagKeyValueSeparator,
			flagQuotationMarks,
			values[key],
			flagQuotationMarks,
		)
	}
}

func retriveInputJson() string {
	input := flag.Arg(0)
	if len(input) < 1 {
		fatalPrint(`The posicional argument "INPUT" must be a non empty string.`)
	}

	if flagRawValue {
		return input
	}

	envValue := os.Getenv(input)
	if len(envValue) < 1 {
		fatalPrint(`The positional argument "INPUT" must reference an existing environment variable name.`)
	}

	return envValue
}

func fatalPrint(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
