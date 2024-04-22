package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// 0 - chose mode
// 1 - edit mode
// 2 - dictation mode
const (
	edit_mode int8 = iota + 1
	dictation_mode
)

// 1 - alert
// 2 - panic with exit
const (
	unpleasant int8 = iota + 1
	critical
)

var infoWords map[int8]string

type CLI struct {
	Mode        int8
	clear       map[string]func()
	Words       map[string]Word `json:"words"`
	S           *bufio.Scanner
	haveToSay   string
	ElementsInR int
}

func Shell() *CLI {
	c := new(CLI)
	c.init()
	c.writeInfo()
	c.readFromDisk(nameOfMainFile)
	return c
}

func (c *CLI) initMode(mode int8) {
	c.Mode = mode
	c.CallClear()
}

func (c *CLI) writeInfo() {
	fmt.Println("-------INFO-------")
	fmt.Println(infoWords[c.Mode])
	if c.haveToSay != "" {
		fmt.Println(c.haveToSay)
		c.haveToSay = ""
	}

}

func (c CLI) closeApp() {
	c.writeToDisk()
	os.Exit(1)
}

func (c *CLI) setHaveToSay(s string) {
	c.haveToSay = s
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
		0:              " (E)M - to enter edit mode \n (D)M to enter dictation mode \n (EX)IT to exit form app",
		edit_mode:      " (L)IST - to list all words and translates \n (E)DIT - to start editing words \n (A)DD - to add new word and translate \n (D)ELETE to delete word or translate \n (S)AVE - to save on disk \n (C)OPY to coppy form disk \n(EX)IT for exit from EDIT mode",
		dictation_mode: " (S)STAR or press Enter - to start dictation \n (C)HANGE - change the number of words in a dictation \n (E)XIT - to exit dictation mode ",
	}
	c.Words = make(map[string]Word)

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

func (c *CLI) readFromDisk(nameOfFile string) {
	file := getFile(nameOfFile)
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

	w, err := json.Marshal(c.Words)
	errorHandler(err, 1)

	if !bytes.Equal(text, w) {
		replica.Write(text)
		file.Write(w)
		c.setHaveToSay("recorded")
	} else {
		c.setHaveToSay("data has not been changed, nothing has been recorded")
	}

}
func (c *CLI) CopyFormReplica() {
	if strings.ToUpper(c.makeFillShift("are you shure ? Y/N")) == "Y" {
		c.readFromDisk(nameOfReplicaFile)
		c.setHaveToSay("copied from replica")
	}
}

func roughtExit(err error) {
	fmt.Println("rought_exit", err)
	time.Sleep(5 * time.Second)
	os.Exit(10)
}
