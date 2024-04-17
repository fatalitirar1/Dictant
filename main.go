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

	}
}

func SwitchingMode(command string, c *CLI) {

	switch command {
	case "EXIT":
		c.closeApp()
	case "EM":
		c.initMode(edit_mode)
	case "DM":
		c.initMode(dictant_mode)
	}
}

func Edit_mode(command string, c *CLI) {

	switch {
	case command == "LIST":
		c.ListW()
	case command == "ADD":
		c.AddW()
	case command == "END":
		c.initMode(0)
	}

}
