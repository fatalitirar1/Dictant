package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type CLI struct {
	Mode      string
	clear     map[string]func()
	Words     []Word
	infoWords map[string]string
	S         *bufio.Scanner
}

func (c *CLI) initDictantMode() {
	fmt.Println("")
}

func (c *CLI) initEditMode() {

}

func (c CLI) writeInfo() {
	fmt.Println(c.infoWords[c.Mode])
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
	c.Mode = "none"
	c.infoWords = map[string]string{
		"none": "EM - to enter edit mode \n DM to enter dictant mode \n EXIT to exit form app",
		"EM":   "LIST - to list all words and translates \n EDIT {word} - to start editing translate or key word \n ADD - to add new word and translate \n DELETE {word} to delete word or translate \n END to exit EDIT mode",
		"DM":   "",
	}

}

func (c CLI) CallClear() {
	value, ok := c.clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                            //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func (c *CLI) readFromDisk() {
	file := findFile()
	err := json.NewDecoder(file).Decode(c.Words)
	if err != nil {
		roughtExit(err)
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
	fmt.Println(err)
	time.Sleep(5 * time.Second)
	os.Exit(10)
}
