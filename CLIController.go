package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	edit_mode int8 = iota + 1
	dictant_mode
)

const (
	unpleasant int8 = iota + 1
	critical
)

var infoWords map[int8]string

type CLI struct {
	Mode  int8
	clear map[string]func()
	Words []Word
	S     *bufio.Scanner
}

func Shell() *CLI {
	c := new(CLI)
	c.init()
	c.writeInfo()
	c.readFromDisk()
	return c
}

func (c *CLI) initMode(mode int8) {
	c.Mode = mode
	c.writeInfo()
}

func (c CLI) writeInfo() {
	fmt.Println(infoWords[c.Mode])
}

func (c CLI) closeApp() {
	c.writeToDisk()
	os.Exit(1)
}

func (c *CLI) init() {
	c.S = bufio.NewScanner(os.Stdin)
	c.clear = make(map[string]func())
	c.clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c.clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c.clear["darwin"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c.Mode = 0
	infoWords = map[int8]string{
		0:            "EM - to enter edit mode \n DM to enter dictant mode \n EXIT to exit form app",
		edit_mode:    "LIST - to list all words and translates \n EDIT {word} - to start editing translate or key word \n ADD - to add new word and translate \n DELETE {word} to delete word or translate \n END to exit EDIT mode",
		dictant_mode: "",
	}

}

func errorHandler(err error, criticality int8) {
	if err != nil {
		switch criticality {
		case unpleasant:
			fmt.Println(err)
		case critical:
			roughtExit(err)
		}
	}
}

func (c CLI) CallClear() {
	value, ok := c.clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                            //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		roughtExit(errors.New("your platform is unsupported! I can't clear terminal screen :("))
	}
}

func (c *CLI) readFromDisk() {
	file := findFile()
	s, err := file.Stat()
	errorHandler(err, 2)
	if s.Size() > 0 {
		errorHandler(json.NewDecoder(file).Decode(&c.Words), 2)
	}
}

func (c *CLI) writeToDisk() {
	file := findFile()
	w, err := json.Marshal(c.Words)
	if err != nil {
		fmt.Println(err)
	}
	file.Write(w)
}

func roughtExit(err error) {
	fmt.Println("rought_exit", err)
	time.Sleep(5 * time.Second)
	os.Exit(10)
}
