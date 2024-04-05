package main

import (
	"strings"
)

func main() {
	c := Init()
	for c.S.Scan() {
		comand := strings.ToUpper(c.S.Text())
		switch {
		case comand == "EXIT":
			c.closeApp()
		case comand == "EM":
			c.initEditMode()
		case comand == "DM":
			c.initDictantMode()
		}
	}
}

func Init() *CLI {
	c := new(CLI)
	c.init()
	c.writeInfo()
	c.readFromDisk()
	return c
}
