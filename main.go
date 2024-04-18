package main

import (
	"strings"
)

func main() {
	c := Shell()
	for c.S.Scan() {

		command := strings.ToUpper(c.S.Text())
		switch c.Mode {
		case edit_mode:
			Edit_mode(command, c)
		case dictant_mode:

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
		c.initMode(dictant_mode)
	}
}

func Edit_mode(command string, c *CLI) {
	switch {
	case command == "LIST" || command == "L":
		c.ListW()
	case command == "ADD" || command == "A":
		c.AddW()
	case command == "EDIT" || command == "E":
		c.EditW()
	case command == "DELETE" || command == "D":
		c.DeleteW()
	case command == "EXIT" || command == "E":
		c.initMode(0)
	}
}
