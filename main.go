package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	text_extractor "gitlab.com/spezi/services/pdf_text_extractor/pkg"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	pdf := r.URL.Query().Get("pdf")
	dataFormat := r.URL.Query().Get("data-format")
	var data interface{}
	var err error = nil

	switch dataFormat {
	case "zip":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.ZippedImages)
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename='images.zip'")
		w.Write(data.([]byte))

		return
	case "array":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextArrayFromImages)
	case "text":
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	default:
		data, err = text_extractor.ExtractDataFromPDF(pdf, text_extractor.TextFromImages)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
