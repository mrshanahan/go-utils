package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mrshanahan/go-utils/auth"
	"github.com/mrshanahan/go-utils/utils"
)

func main() {
	code := Run()
	os.Exit(code)
}

func Run() int {
	if len(os.Args) > 1 && utils.Any(os.Args[1:], func(x string) bool { return x == "-h" || x == "--help" || x == "-?" }) {
		printHelp()
		return 0
	}

	command := ""
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch strings.ToLower(command) {
	case "generate-salt":
		saltLength := 32
		if len(os.Args) > 2 {
			saltLengthStr := os.Args[2]
			l, err := strconv.Atoi(saltLengthStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: non-integer argument to generate-salt: %s\n", saltLengthStr)
				return -1
			}
			if l <= 0 {
				fmt.Fprintf(os.Stderr, "error: argument to generate-salt should be greater than 0: %d\n", l)
				return -1
			}
			saltLength = l
		}
		salt := auth.GenerateSalt(saltLength)
		saltEnc := hex.EncodeToString(salt)
		fmt.Println(saltEnc)
		return 0
	case "hash":
		if len(os.Args) <= 2 {
			fmt.Fprintln(os.Stderr, "error: must provide an argument for salt")
			return -1
		}
		if len(os.Args) <= 3 {
			fmt.Fprintln(os.Stderr, "error: must provide an argument for input")
			return -1
		}
		saltEnc, input := os.Args[2], os.Args[3]
		salt, err := hex.DecodeString(saltEnc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to decode hex salt value: %s", err)
			return -1
		}
		inputBytes := []byte(input)
		hashed, err := auth.Hash(salt, inputBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to hash input: %s", err)
			return -1
		}
		hashedEnc := hex.EncodeToString(hashed)
		fmt.Println(hashedEnc)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "error: invalid command: %s\n", command)
		printHelp()
		return 1
	}
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `
go-utils [-h|--help] [generate-salt|hash]

DESCRIPTION:
    Library functions defined in the mrshanahan/go-utils packaged exposed as CLI tools.

GLOBAL FLAGS:
    -h|--help                  Display this message and exit

COMMANDS:
    generate-salt [<length>]   Generate a new salt value with <length> bytes. Defaults to 32. Output is in hex.
    hash <salt> <value>        Perform a SHA256 hash of <value> using the salt <salt>. Output is in hex.
`)
}
