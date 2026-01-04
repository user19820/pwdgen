package args

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
)

const (
	CmdGen  string = "gen"
	CmdGet  string = "get"
	CmdInit string = "init"

	defaultPwdLen int = 25
)

type Cmd struct {
	Type   string
	Name   string
	Length int
	Debug  bool
}

//nolint:mnd,gocognit // it's fine
func Init() (Cmd, error) {
	var cmd Cmd
	help := flag.Bool("h", false, "prints help text")
	debug := flag.Bool("d", false, "debug mode (for development)")

	flag.Parse()

	if help != nil && *help {
		printHelp()
	}

	if debug != nil && *debug {
		cmd.Debug = true
	}

	args := flag.Args()
	numArgs := len(args)

	if numArgs == 0 {
		printHelp()
		return Cmd{}, errors.New("command not specified")
	}

	switch args[0] {
	case CmdGen:
		if cmd.Debug {
			fmt.Println("Gen:")
			fmt.Printf("\tName: %s\n", args[1])
			if numArgs == 3 {
				fmt.Printf("\tLength: %s\n", args[2])
			}
		}

		name := args[1]
		if name == "" {
			return Cmd{}, errors.New("gen: name argument must be provided")
		}

		var length int
		if numArgs == 3 {
			// an optional length has been provided
			l, err := strconv.Atoi(args[2])
			if err != nil {
				return Cmd{}, errors.New("gen: invalid length argument provided")
			}
			length = l
		}
		if length == 0 {
			length = defaultPwdLen
		}

		return Cmd{
			Type:   CmdGen,
			Name:   name,
			Length: length,
		}, nil
	case CmdGet:
		if cmd.Debug {
			fmt.Println("Get:")
			fmt.Printf("\tName: %s\n", args[1])
		}
		name := args[1]
		if name == "" {
			return Cmd{}, errors.New("gen: name argument must be provided")
		}

		return Cmd{
			Type: CmdGet,
			Name: name,
		}, nil
	case CmdInit:
		if cmd.Debug {
			fmt.Println("Init")
		}

		return Cmd{
			Type: CmdInit,
		}, nil
	default:
		return Cmd{}, errors.New("invalid command provided")
	}
}

func printHelp() {
	fmt.Println(`
pwdgen, a minimal password generator

commands:
	- init: initialize pwdgen so you can start using it.

	- gen <name> <length>: generate a random password with the provided length.
	  Must provide a name for the thing we're creating a password for.

	- get <name>: retrieve a password for <name>.
	  Password is not actually printed, but copied to clipboard

flags:
	- h: help, prints this message
	- d: debug, used for development purposes only
		`)
	os.Exit(1)
}
