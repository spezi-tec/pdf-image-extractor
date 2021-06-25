package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	text_extractor "gitlab.com/spezi/services/pdf_text_extractor/pkg"
)

func main() {

	log.Print("Listening")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		base64Handler(w, r)
	case "POST":
		fileHandler(w, r)
	}
}

func fileHandler(w http.ResponseWriter, r *http.Request) {

	dataFormat := r.URL.Query().Get("data-format")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 20 MB files.
	r.ParseMultipartForm(20 << 20)

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// read all of the contents of our uploaded file into a
	// byte array

	pdf := text_extractor.EncodeFileB64(file)
	var data interface{}

	switch dataFormat {
	case "zip":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.ZippedImages)
		if err != nil {
			log.Print(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='images.zip'")
		w.Write(data.([]byte))
		// io.Copy(w, data.([]byte))

	case "array":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextArrayFromImages)
	case "text":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	default:
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	}

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

func base64Handler(w http.ResponseWriter, r *http.Request) {

	pdf := r.URL.Query().Get("pdf")
	dataFormat := r.URL.Query().Get("data-format")
	var data interface{}
	var err error = nil

	switch dataFormat {
	case "zip":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.ZippedImages)
		if err != nil {
			log.Print(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='images.zip'")
		w.Write(data.([]byte))

	case "array":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextArrayFromImages)
	case "text":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	default:
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	}

	if err != nil {
		log.Print(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
