package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var all commands

type command struct {
	Flags       *flag.FlagSet
	Command     string
	Description string
	Run         func(*command, *File) error
}

type commands []*command

func (c commands) Len() int {
	return len(c)
}

func (c commands) Less(i, j int) bool {
	return strings.Compare(c[i].Command, c[j].Command) < 0
}

func (c commands) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c *command) Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "\n  %s : %s\n", c.Command, c.Description)
	c.Flags.PrintDefaults()
}

func (c commands) Usage() {
	_, _ = fmt.Fprintf(os.Stderr, "%s <global options> command <command options> : is the command to validate/generate calendar events\n", os.Args[0])
	_, _ = fmt.Fprintf(os.Stderr, "\nGlobal Options:\n")
	flag.PrintDefaults()
	_, _ = fmt.Fprintf(os.Stderr, "\nCommands:")
	for i := range c {
		c[i].Usage()
	}
}

func (c commands) Find(cmd string, args []string) (*command) {
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

func registerCommand(cmd *command) {
	all = append(all, cmd)
	sort.Sort(all)
}

func newCommand(name, desc string, fn func(*command, *File) error) *command {
	return &command{
		Flags:       flag.NewFlagSet(name, flag.ExitOnError),
		Command:     name,
		Description: desc,
		Run:         fn,
	}
}
