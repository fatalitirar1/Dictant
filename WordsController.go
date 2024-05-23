package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	nameOfMainFile    = *flag.String("NameOfMainFile", "Words.json", "file which use to save last data ")
	nameOfReplicaFile = *flag.String("NameOfReplica", "r_Words.json", "Replica file")
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
	if c.isQuitExit(w) {
		return
	}
	newWord.KeyWord = w

	c.setHaveToSay(fmt.Sprint("word: ", w, " added"))

	c.Words[w] = *newWord
}

func (c *CLI) makeFillShift(msg string) string {
	fmt.Println(msg)
	c.S.Scan()
	return c.S.Text()
}

func (c CLI) ListW() {
	for k := range c.Words {
		fmt.Println(k)
	}
}

func (c *CLI) isQuitExit(s string) bool {
	return strings.ToUpper(s) == "Q"
}
func (c *CLI) DeleteW() {
	defer c.CallClear()
	wordToDell := c.makeFillShift("enter word to delete, or type 'q' for exit:")
	if c.isQuitExit(wordToDell) {
		c.setHaveToSay("nothing changed")
		return
	}
	delete(c.Words, wordToDell)
	c.setHaveToSay(fmt.Sprint("word: ", wordToDell, " deleted"))

}

func (c *CLI) EditW() {
	defer c.CallClear()
	wordToEddit := c.makeFillShift("type word which you want to found, or type 'q' for exit:")
	if c.isQuitExit(wordToEddit) {
		return
	}
	var w string
	if _, ok := c.Words[wordToEddit]; ok {

		w = c.makeFillShift(fmt.Sprint("selected(", wordToEddit, ") ",
			"enter new word, or type 'q' for exit:"))

		if c.isQuitExit(w) {
			c.setHaveToSay("no changes")
			return
		}
		delete(c.Words, wordToEddit)
		c.Words[w] = Word{w}

	} else {
		fmt.Println("word doesn't exist in collection")
		answer := c.makeFillShift("try agan ? Y/N")
		if strings.ToUpper(answer) == "Y" {
			c.EditW()
		} else {
			c.setHaveToSay("no changes")
			return
		}
	}
	c.setHaveToSay(fmt.Sprint("changed: ", wordToEddit, " -> ", w))
}

func (c CLI) StartDictant() {
	c.CallClear()
	lenWords := len(c.Words)
	if c.ElementsInR == 0 {
		c.ChangeRange()
	}

	if c.ElementsInR >= lenWords {
		for k := range c.Words {
			fmt.Println(k)
		}
	} else {
		r := rand.New(rand.NewSource(int64(time.Now().Second())))

		shufled := []string{}
		for k := range c.Words {
			shufled = append(shufled, k)
		}
		for l := 0; l <= lenWords/2; l++ {
			id_1, id_2 := r.Intn(lenWords), r.Intn(lenWords)
			shufled[id_1], shufled[id_2] = shufled[id_2], shufled[id_1]
		}

		cof := len(shufled) / c.ElementsInR
		ids := make(map[int]struct{})
		for i := 0; c.ElementsInR > len(ids); i += cof {
			ids[i] = struct{}{}
		}

		for k := range ids {
			fmt.Println(shufled[k])
		}
	}
}

func (c *CLI) ChangeRange() {
	n := c.makeFillShift("how many words should there be in a dictation? (nuber) or 'q' for exit")
	if c.isQuitExit(n) {
		if c.ElementsInR == 0 {
			c.setHaveToSay("can't be 0")
			c.ChangeRange()
		} else {
			c.setHaveToSay(fmt.Sprint("no changes, words will be", c.ElementsInR))
			return
		}

	} else {
		num, err := strconv.Atoi(n)
		errorHandler(err, 1)
		if err != nil {
			if num > len(c.Words) {
				c.ElementsInR = len(c.Words)
			} else if num == 0 {
				c.setHaveToSay("can't be 0")
				c.ChangeRange()
			} else {
				c.ElementsInR = num
			}
		}
	}

}
