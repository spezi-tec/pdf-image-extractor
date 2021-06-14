package main

import (
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
	pdfText, err := text_extractor.ExtractTextFromPDF(pdf)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(pdfText))
}
