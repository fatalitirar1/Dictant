package main

import (
	"fmt"
	"os"
)

func findFile() *os.File {
	file, err := os.OpenFile("words.json", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		roughtExit(err)
	}
	return file
}

type Word struct {
	KeyWord string   `json:"keyWord"`
	T       []string `json:"translate"`
}

func (c *CLI) AddW() {
	newWord := new(Word)

	newWord.KeyWord = c.makeFillShift("pls type key word:")
	fmt.Println("enter 'q' to stop entering translations ")

	for c.S.Scan() {

		t := c.S.Text()
		fmt.Println("typed", t)
		if t == "q" {

			break
		} else {
			newWord.T = append(newWord.T, t)
		}
	}
	c.Words = append(c.Words, *newWord)
	c.CallClear()
	fmt.Println("added", newWord)
	c.writeInfo()
}

func (c *CLI) makeFillShift(msg string) string {
	fmt.Println(msg)
	c.S.Scan()
	return c.S.Text()
}

func (c CLI) ListW() {
	fmt.Println(c.Words)
}

func (c *CLI) EditW() {
	//fmt.Println(c.Words)
}
