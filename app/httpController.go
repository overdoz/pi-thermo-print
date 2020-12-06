package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

func ServeFile(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	dir, _ := os.Getwd()
	http.FileServer(http.Dir(dir)).ServeHTTP(w,r)
}

func (s Session)homeHandler(w http.ResponseWriter, r *http.Request) {

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
			RenameFile(date)

			// send data to firestore database
			s.SendToFirebase(msg, date)

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
			RenameFile(date)

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
