package main

import (
	"fmt"
	"os"
)

func getFile(name string) *os.File {
	var file *os.File
	var err error
	if fileExits(name) {
		file, err = os.OpenFile(name, os.O_RDWR, 0660)
		if err != nil {
			roughtExit(err)
		}
	} else {
		os.Create(name)
		file = getFile(name)
	}

	return file
}

func fileExits(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

type Word struct {
	KeyWord string `json:"keyWord"`
}

func (c *CLI) AddW() {
	defer c.CallClear()
	newWord := new(Word)
	w := c.makeFillShift("pls type key word or q for exit:")
	if w == "q" {
		return
	}
	newWord.KeyWord = w

	c.haveToSay = fmt.Sprint("word: ", w, " added")

	c.Words = append(c.Words, *newWord)
}

func (c *CLI) makeFillShift(msg string) string {
	fmt.Println(msg)
	c.S.Scan()
	return c.S.Text()
}

func (c CLI) ListW() {
	for _, w := range c.Words {
		fmt.Println(w)
	}
}
func (c *CLI) DeleteW() {
	defer c.CallClear()
	wordToDell := c.makeFillShift("enter word to delete, or type 'q' for exit:")
	if wordToDell == "q" {
		return
	}
	for i, w := range c.Words {
		if w.KeyWord == wordToDell {
			c.Words = append(c.Words[:i], c.Words[i+1:]...)
		}
	}
	c.haveToSay = fmt.Sprint("word: ", wordToDell, " deleted")

}

func (c *CLI) EditW() {
	defer c.CallClear()
	wordToEddit := c.makeFillShift("enter word to edit, or type 'q' for exit:")
	if wordToEddit == "q" {
		return
	}
	id_cwords := -1
	for i, w := range c.Words {
		if w.KeyWord == wordToEddit {
			id_cwords = i
			break
		}
	}
	if id_cwords == -1 {
		fmt.Println("word doesn't exist in collection")
		answer := c.makeFillShift("try agan ? Y/N")
		if answer == "Y" {
			c.EditW()
		} else {
			return
		}

	} else {
		fmt.Println("enter word to edit, or type 'q' for exit:")
		fmt.Print(c.Words[id_cwords].KeyWord)
		c.S.Scan()
		w := c.S.Text()
		c.Words[id_cwords].KeyWord = w
	}
	c.haveToSay = fmt.Sprint("word: ", wordToEddit, " changed")
}
