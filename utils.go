package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func RenameFile(n string) string {
	oldName := FilenameTxt
	newName := n + ".txt"
	err := os.Rename(oldName, newName)
	if err != nil {
		log.Fatal(err)
	}
	return newName
}

func StoreToArchiv(name, dir string) {
	err := os.Rename("./"+name, "./"+dir+"/"+name)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func ConcatTextAndDate(message, date string) string {
	output := message + "\n\n\n\n         " + date + "\n\n   "
	return output
}

// adjust width
func FormatString(text string, maxLen, longest int) string {
	output := ""
	col := 0
	maxCols := maxLen - longest
	fixNewLine := strings.ReplaceAll(text, "\n", " \n ")
	tempText := strings.Split(fixNewLine, " ")
	_ = strings.Repeat(" ", longest)

	for _, word := range tempText {
		count := utf8.RuneCountInString(word)
		if strings.Contains(word, "\n") {
			col = 0
		}
		if col == 0 {
			if word != "\n" {
				// add word to output string
				output += word + " "

				// increase rune counter
				col += count + 1
			} else {
				// add newline to output string
				output += word
			}
			// check if new word is too big for this line
		} else if (col + count) <= maxCols {

			// if it fits perfectly into line don't add a space rune
			if (col + count) == maxCols {
				output += word
			} else {
				output += word + " "
			}

			// if it's not new line, increase counter
			if !strings.Contains(word, "\n") {
				col += count + 1
			}
			// word is too long for this line
			// better add a new line
		} else {
			output += "\n" + word + " "
			col = longest + count + 1
		}
	}
	log.Println("\r\n" + output)
	return output
}
