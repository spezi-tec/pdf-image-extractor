package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	text_extractor "gitlab.com/spezi/services/pdf_text_extractor/pkg"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Fatal(err)
		}
	}()

	router := mux.NewRouter()

	router.HandleFunc("/", handler).Queries("pdf", "{pdf}", "data-format", "{data-format}").Methods("GET")

	log.Print("Listening")

	// servePath := fmt.Sprintf(":8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}

func handler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	pdf := params["pdf"]
	dataFormat := params["data-format"]

	var data interface{}
	var err error = nil

	switch dataFormat {
	case "zip":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.ZippedImages)
		if err != nil {
			log.Fatal(err)
			return
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
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
