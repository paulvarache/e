package main

import (
	"flag"
	"fmt"
	"os"
)

// These variables are provided at build time using ldflags
var (
	BuildVersion string = "dev"
	BuildHash    string = "dev"
)

type Subcommand struct {
	FlagSet *flag.FlagSet
	Command string
	Usage   string
}

func NewSubcommand(command string, name string, usage string) *Subcommand {
	return &Subcommand{
		Command: name,
		Usage:   usage,
		FlagSet: flag.NewFlagSet(command, flag.ExitOnError),
	}
}

func main() {
	versionFlag := flag.Bool("version", false, "Prints the version")

	createCommand := NewSubcommand("create", "create <name>", "Creates a new env profile")
	listCommand := NewSubcommand("list", "list", "Lists the available env profiles")
	setCommand := NewSubcommand("set", "set <key> <value>", "Sets an env variable for the current profile")
	defaultCommand := NewSubcommand("select", "<name>", "Selects an env profile to be used")

	commands := []*Subcommand{
		createCommand,
		listCommand,
		setCommand,
		defaultCommand,
	}

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		max := 0
		for _, c := range commands {
			if len(c.Command) > max {
				max = len(c.Command)
			}
		}
		for _, c := range commands {
			if len(c.Command) > max {
				max = len(c.Command)
			}
			fmt.Fprintf(flag.CommandLine.Output(), "  %s  %s\n", fmt.Sprintf("%-*s", max, c.Command), c.Usage)
		}
		flag.PrintDefaults()
	}

	flag.Parse()

	if *versionFlag {
		fmt.Printf("%s-%s", BuildVersion, BuildHash)
		os.Exit(0)
	}

	em, err := Load()
	if err != nil {
		panic(err)
	}

	args := flag.Args()

	if len(args) == 0 {
		values, err := em.Profiles[em.Selected].GetValues()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Selected profile: %s\n\n", em.Selected)
		longestKey := 0
		for k := range values {
			if len(k) > longestKey {
				longestKey = len(k)
			}
		}
		for k, v := range values {
			fmt.Printf("%s = %s\n", fmt.Sprintf("%-*s", longestKey, k), v)
		}
		os.Exit(0)
	}

	switch args[0] {
	case "create":
		err := createCommand.FlagSet.Parse(args[1:])
		if err != nil {
			panic(err)
		}
		args := createCommand.FlagSet.Args()
		name := args[0]
		_, err = em.CreateProfile(name)
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	case "list":
		err := listCommand.FlagSet.Parse(args[1:])
		if err != nil {
			panic(err)
		}
		for name, p := range em.Profiles {
			if name == em.Selected {
				fmt.Printf("*%s\n", p.Name)
			} else {
				fmt.Println(p.Name)
			}
		}
		os.Exit(0)
	case "set":
		err := setCommand.FlagSet.Parse(args[1:])
		if err != nil {
			panic(err)
		}
		args := setCommand.FlagSet.Args()
		if len(args) == 0 {
			setCommand.FlagSet.PrintDefaults()
			os.Exit(1)
		}
		var key string
		var value string
		if len(args) == 1 {
			key = args[0]
			value = ""
		} else if len(args) == 2 {
			key = args[0]
			value = args[1]
		}
		err = em.GetProfile().SetValue(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("#env")
		fmt.Printf("%s%s=%s\n", SetVarPrefix, key, value)
	default:
		if len(args) == 1 {
			name := args[0]
			var oldValues ProfileValues
			if em.Selected != "" {
				oldValues, err = em.GetProfile().GetValues()
				if err != nil {
					panic(err)
				}
			} else {
				oldValues = make(ProfileValues)
			}
			err = em.SelectProfile(name)
			if err != nil {
				panic(err)
			}
			values, err := em.GetProfile().GetValues()
			if err != nil {
				panic(err)
			}
			fmt.Println("#env")
			for k := range oldValues {
				fmt.Printf("%s%s=\"\"\n", SetVarPrefix, k)
			}
			for k, v := range values {
				fmt.Printf("%s%s=%s\n", SetVarPrefix, k, v)
			}
		} else {
			flag.Usage()
			os.Exit(1)
		}
	}
}
