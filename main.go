package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
)

const (
	CLIExec string = "CLI"
)

func main() {
	wToExec := flag.String("WToExec", CLIExec, "way to execute programm")
	flag.Parse()

	switch *wToExec {
	case CLIExec:
		CLIStart()
	}

}

func CLIStart() {
	c := Shell()

	exitNotify := make(chan os.Signal, 1)
	signal.Notify(exitNotify, os.Interrupt)

	go func(exitNotify chan os.Signal) {
		<-exitNotify
		if len(c.Words) > 0 {
			c.writeToDisk()
		}
		os.Exit(1)
	}(exitNotify)

	for c.S.Scan() {

		command := strings.ToUpper(c.S.Text())
		switch c.Mode {
		case edit_mode:
			EditMode(command, c)
		case dictation_mode:
			DictantMde(command, c)
		default:
			SwitchingMode(command, c)
		}
		c.writeInfo()
	}
}

func SwitchingMode(command string, c *CLI) {

	switch {
	case command == "EXIT" || command == "EX":
		c.closeApp()
	case command == "EM" || command == "E":
		c.initMode(edit_mode)
	case command == "DM" || command == "D":
		c.initMode(dictation_mode)
	}
}

func EditMode(command string, c *CLI) {
	switch {
	case command == "LIST" || command == "L":
		c.ListW()
	case command == "ADD" || command == "A":
		c.AddW()
	case command == "EDIT" || command == "E":
		c.EditW()
	case command == "DELETE" || command == "D":
		c.DeleteW()
	case command == "EXIT" || command == "EX":
		c.writeToDisk()
		c.initMode(0)
	case command == "SAVE" || command == "S":
		c.writeToDisk()
	case command == "COPY" || command == "C":
		c.CopyFormReplica()
	}
}

func DictantMde(command string, c *CLI) {
	switch {
	case command == "START" || command == "S":
		c.StartDictant()
	case command == "CHANGE" || command == "C":
		c.ChangeRange()
	case command == "EXIT" || command == "E":
		c.initMode(0)
	default:
		c.StartDictant()
	}
}
