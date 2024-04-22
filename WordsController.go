package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const nameOfMainFile string = "words.json"
const nameOfReplicaFile string = "old_words.json"

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
		return
	}
	delete(c.Words, wordToDell)
	c.setHaveToSay(fmt.Sprint("word: ", wordToDell, " deleted"))

}

func (c *CLI) EditW() {
	defer c.CallClear()
	wordToEddit := c.makeFillShift("type word which you want to found, or type 'a' for exit:")
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

func (c *CLI) StartDictant() {
	c.CallClear()
	if c.ElementsInR == 0 {
		c.ChangeRange()
	}

	if c.ElementsInR >= len(c.Words) {
		for k := range c.Words {
			fmt.Println(k)
		}
	} else {
		r := rand.New(rand.NewSource(int64(time.Now().Second())))

		Words := []string{}
		for k := range c.Words {
			Words = append(Words, k)
		}
		cof := len(Words) / c.ElementsInR
		ids := []int{}
		for i := rand.Intn(cof) + 1; c.ElementsInR > len(ids); i += r.Intn(cof-1) + 1 {
			ids = append(ids, i)
		}

		for ii := 0; ii < c.ElementsInR; ii++ {
			id1 := r.Intn(c.ElementsInR - 1)
			id2 := r.Intn(c.ElementsInR - 1)
			forTime := ids[id2]
			ids[id2] = ids[id1]
			ids[id1] = forTime
		}

		for _, k := range ids {
			fmt.Println(Words[k])
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
