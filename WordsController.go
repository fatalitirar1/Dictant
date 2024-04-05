package main

import (
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
	keyWord string            `json:"keyWord"`
	T       map[string]string `json:"translate"`
}
