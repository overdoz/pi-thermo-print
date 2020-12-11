package main

import (
	"encoding/json"
	// firebase "firebase.google.com/go"
	// "fmt"
	// "golang.org/x/net/context"
	// "google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// init firebase
	// opt := option.WithCredentialsFile("./service-account-file.json")

	// link to firebase project
	// config := &firebase.Config{ProjectID: os.Getenv("PRINTER_ID")}

	// Use a service account
	// ctx := context.Background()

	//app, err := firebase.NewApp(ctx, config, opt)
	//if err != nil {
	//	_ = fmt.Errorf("error initializing app: %v", err)
	//}

	// init firestore client
	//client, err := app.Firestore(ctx)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//defer client.Close()

	// store := Session{Fs: client, Ctx: ctx}

	// handle incoming requests
	http.HandleFunc("/", homeHandler)

	log.Println("listening on port 8001")

	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// check for http method
	switch r.Method {

	// serve React PWA
	// https://github.com/overdoz/pi-thermo-pwa
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

			msg := FormatString(m.Text, LineWidth, 0)
			date := time.Now().Format("2006-01-02 15:04:05")
			o := ConcatTextAndDate(msg, date)

			// send a txt file to the printer
			SendToPrinter(o, FilenameTxt)
			newName := RenameFile(date)
			StoreToArchiv(newName, ArchivName)

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
	http.FileServer(http.Dir(dir)).ServeHTTP(w, r)
}
