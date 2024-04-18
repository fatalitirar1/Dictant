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

// 0 - chose mode
// 1 - edit mode
// 2 - dictant mode
const (
	edit_mode int8 = iota + 1
	dictant_mode
)

const (
	unpleasant int8 = iota + 1
	critical
)
const nameOfMainFile = "words.json"
const nameOfReplicaFile = "old_words.json"

var infoWords map[int8]string

type CLI struct {
	Mode      int8
	clear     map[string]func()
	Words     []Word `json:"words"`
	S         *bufio.Scanner
	haveToSay string
}

func Shell() *CLI {
	c := new(CLI)
	c.init()
	c.writeInfo()
	c.readFromDisk()
	return c
}

func (c *CLI) initMode(mode int8) {
	if c.Mode == edit_mode {
		c.writeToDisk()
	}
	c.Mode = mode
	c.CallClear()
}

func (c *CLI) writeInfo() {
	fmt.Println("-------INFO-------")
	fmt.Println(infoWords[c.Mode])
	fmt.Println(c.haveToSay)
	c.haveToSay = ""
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
		0:            " (E)M - to enter edit mode \n (D)M to enter dictant mode \n (EX)IT to exit form app",
		edit_mode:    " (L)IST - to list all words and translates \n (E)DIT - to start editing words \n (A)DD - to add new word and translate \n (D)ELETE to delete word or translate \n EXIT for exit from EDIT mode",
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
	file := getFile(nameOfMainFile)
	defer file.Close()
	s, err := file.Stat()
	errorHandler(err, 2)
	if s.Size() > 0 {
		errorHandler(json.NewDecoder(file).Decode(&c.Words), 2)
	}
}

func (c *CLI) writeToDisk() {

	file := getFile(nameOfMainFile)
	defer file.Close()

	text, err := os.ReadFile(nameOfMainFile)
	errorHandler(err, 1)
	replica := getFile(nameOfReplicaFile)
	defer replica.Close()

	replica.Write(text)

	w, err := json.Marshal(c.Words)
	errorHandler(err, 1)
	file.Write(w)

}

func roughtExit(err error) {
	fmt.Println("rought_exit", err)
	time.Sleep(5 * time.Second)
	os.Exit(10)
}
