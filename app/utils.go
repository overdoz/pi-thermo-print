package app

import (
	"fmt"
	"github.com/disintegration/imaging"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode/utf8"
)


// edit immage
func AdjustImage(file string) error {
	src, _ := imaging.Open(file)
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)
	src = imaging.AdjustBrightness(src, 30)
	src = imaging.Grayscale(src)
	src = imaging.AdjustContrast(src, -20)
	src = imaging.AdjustGamma(src, 0.75)
	return imaging.Save(src, FILENAME_PNG)
}

// file should have the ending .png
func PrintImage(file string) {
	sh := "lp " + file + " -d LKT"

	args := strings.Split(sh, " ")

	cmd := exec.Command(args[0], args[1:]...)

	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
	}
}


func ConcatTextAndDate(message, date string) string {
	log.Print("\n " + message)
	output := message + "\n\n\n\n         " + date + "\n\n   "
	return output
}

func FindNames(s string) string {
	longestString := 0
	outputString := ""
	stringArray := strings.Split(s, " ")

	// slice without predefined length
	var in []int

	// [name]quote
	q := make(map[string]string)

	// find longest name
	for i, s := range stringArray {
		if strings.Contains(s, ":") {
			// save index of names
			in = append(in, i)
			if utf8.RuneCountInString(s) > longestString {
				longestString = utf8.RuneCountInString(s)
			}
		}
	}

	// connect quote to name
	for i := 0; i < len(in); i++ {
		if i < len(in)-1 {
			currentName := in[i]
			nextName := in[i+1]
			// connect quote to person
			// +1 to cut name at the beginning
			q[stringArray[currentName]] = strings.Join(stringArray[currentName+1:nextName], " ")
		} else {
			// if it's the last person, take the rest of the string
			currentName := in[i]
			q[stringArray[currentName]] = strings.Join(stringArray[currentName+1:], " ")
		}
	}

	// concat names and quotes
	for i, s := range q {
		outputString = outputString + i + "\n" + FormatString(s, LINE_WIDTH, longestString) + "\n\n"
	}
	return outputString
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
	fmt.Println(output)
	return output
}
// return the wordcount of word after new line
func digitsAfterNewLine(w string) int {
	s := strings.Split(w, "\n")
	return utf8.RuneCountInString(s[1])
}

func RenameFile(n string) {
	oldName := "test.txt"
	newName := n + ".txt"
	err := os.Rename(oldName, newName)
	if err != nil {
		log.Fatal(err)
	}
}

