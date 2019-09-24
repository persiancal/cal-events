package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var all Commands

type Command struct {
	Flags       *flag.FlagSet
	Command     string
	Description string
	Run         func(*Command, *File) error
}

type Commands []*Command

func (c Commands) Len() int {
	return len(c)
}

func (c Commands) Less(i, j int) bool {
	return strings.Compare(c[i].Command, c[j].Command) < 0
}

func (c Commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "\n  %s : %s\n", c.Command, c.Description)
	c.Flags.PrintDefaults()
}

func (c Commands) Usage() {
	fmt.Fprintf(os.Stderr, "%s <global options> command <command options> : is the command to validate/generate calendar events\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nGlobal Options:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nCommands:")
	for i := range c {
		c[i].Usage()
	}
}

func (c Commands) Find(cmd string, args []string) (*Command) {
	for i := range c {
		if c[i].Command == cmd {
			if err := c[i].Flags.Parse(args); err != nil {
				c[i].Flags.Usage()
				os.Exit(1)
			}

			return c[i]
		}
	}

	c.Usage()
	os.Exit(1)
	return nil
}

func Register(cmd *Command) {
	all = append(all, cmd)
	sort.Sort(all)
}

func NewCommand(name, desc string, fn func(*Command, *File) error) *Command {
	return &Command{
		Flags:       flag.NewFlagSet(name, flag.ExitOnError),
		Command:     name,
		Description: desc,
		Run:         fn,
	}
}
