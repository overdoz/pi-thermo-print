package app

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)



// file should have the ending .txt
func SendToPrinter(t, file string) {
	printText := []byte(t)

	err := ioutil.WriteFile("./" + file, printText, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// print command to LKT
	sh := "lp " + file + " -d LKT"

	args := strings.Split(sh, " ")

	cmd := exec.Command(args[0], args[1:]...)

	// execute command
	b, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
	}
	log.Printf("%s \n", b)
}

func SendFileToPrinter(file io.Reader) {
	log.Println("processing image...")

	// create png file
	tempFile, _ := os.Create(FILENAME_PNG)

	// read request body from client
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Cannot read request body: ", err)
	}

	// write incoming file to new file
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Printf("write error", err)
	}

	// close temporary file
	err = tempFile.Close()
	if err != nil {
		log.Printf("could not open", err)
	}

	// Read image from file that already exists
	err = AdjustImage(FILENAME_PNG)
	if err != nil {
		log.Fatalf("failed to save image after edit: %v", err)
	}

	// print the png file
	PrintImage(FILENAME_PNG)
}
