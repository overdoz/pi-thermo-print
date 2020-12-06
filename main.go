package main

import (
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/disintegration/imaging"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"unicode/utf8"
)

const FILENAME_TXT = "test.txt"
const FILENAME_PNG = "test.png"
const LINE_WIDTH = 28
const PORT = ":8001"



type Message struct {
	Text string
}


func main() {
	// init firebase
	opt := option.WithCredentialsFile("./service-account-file.json")

	// link to firebase project
	config := &firebase.Config{ProjectID: "airprinter-8c2ee"}

	// Use a service account
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		_ = fmt.Errorf("error initializing app: %v", err)
	}

	// init firestore client
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// store := Session{Fs: client, Ctx: ctx}

	// handle incoming requests
	http.HandleFunc("/", homeHandler)

	log.Println("listening on port 8001")

	err = http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// check for http method
	switch r.Method {

	// serve React PWA
	// https://github.com/overdoz/airprinter
	case http.MethodGet:
		ServeFile(w, r)

	// two options to send files
	// @param text: send formatted text
	// @param file: send multipart/form-data
	// TODO: right now you either send a file or plain text, but not combined
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			log.Printf("ParseForm() err: %v", err)
			return
		}

		// POST Request params
		responseType := r.FormValue("type")

		// determine if type =? text or file
		switch responseType {
		// text will be formatted correctly
		// print timestamp at the end of the sheet
		case "text":
			var m Message

			decoder := json.NewDecoder(r.Body)

			err := decoder.Decode(&m)
			if err != nil {
				log.Fatal("JSON Decoder failed: ", err)
			}

			msg := FormatString(m.Text, LINE_WIDTH, 0)
			date := time.Now().Format("2006-01-02 15:04:05")
			o := ConcatTextAndDate(msg, date)

			// send a txt file to the printer
			SendToPrinter(o, FILENAME_TXT)
			newName := RenameFile(date)
			StoreToArchiv(newName, "archiv")

			// send data to firestore database
			// s.SendToFirebase(msg, date)

			// redirect to previous page
			http.Redirect(w, r, r.Header.Get("Referer"), 302)

		// POST Request params (type = "files")
		case "files":
			// Parse our multipart form, 10 << 20 specifies a maximum
			// upload of 10 MB files.
			err := r.ParseMultipartForm(10 << 20)

			// FormFile returns the first file for the given key `myFile`
			// it also returns the FileHeader so we can get the Filename,
			// the Header and the size of the file
			file, _, err := r.FormFile("file")
			if err != nil {
				log.Fatal(err)
			}

			// create an image file and send it to printer
			SendFileToPrinter(file)

			// redirect to previous page
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
		default:
			// redirect to previous page
			http.Redirect(w, r, r.Header.Get("Referer"), 302)
			log.Print("couldn't read query")
		}
	default:
		log.Print("couldn't handle request")
	}
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dir, _ := os.Getwd()
	http.FileServer(http.Dir(dir)).ServeHTTP(w,r)
}

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

func RenameFile(n string) string {
	oldName := "test.txt"
	newName := n + ".txt"
	err := os.Rename(oldName, newName)
	if err != nil {
		log.Fatal(err)
	}
	return newName
}

func StoreToArchiv(name, dir string) {
	err :=  os.Rename("./" + name, "./" + dir + "/" + name)

	if err != nil {
		fmt.Println(err)
		return
	}
}



