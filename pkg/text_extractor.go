package text_extractor

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/gosseract"
	"gopkg.in/gographics/imagick.v2/imagick"
)

// ExtractTextFromPDF will take a fbase64 string of a pdf file and convert the file into an
// image which has its contents extracted via tesseract. It will create the image as a
// high resolution jpg file with minimal compression.
func ExtractTextFromPDF(base64PDF string) (string, error) {
	client := gosseract.NewClient()
	client.SetLanguage("por")

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	//adding default config to mw image
	if err := SetupImage(base64PDF, mw); err != nil {
		fmt.Println(err)
		return "", err
	}

	var imageName string
	var pdfText string = ""

	// Iterate over PDF pages
	for i := 0; i < int(mw.GetNumberImages()); i++ {
		// Set Page Index
		mw.SetIteratorIndex(i)
		imageName = fmt.Sprintf("pdf_page_%v.jpg", i)
		// Save Image
		if err := mw.WriteImage(imageName); err != nil {
			return "", err
		}

		text, err := ExtractTextFromImage(client, imageName)
		if err != nil {
			return "", err
		}

		pdfText += text
	}

	return pdfText, nil
}

// ExtractTextFromPDF will take a filename of a pdf file and convert the file into an
func SetupImage(base64PDF string, mw *imagick.MagickWand) error {
	dec, err := base64.StdEncoding.DecodeString(base64PDF)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile("", "external.*.pdf")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(dec); err != nil {
		return err
	}

	if err := file.Sync(); err != nil {
		return err
	}

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(300, 300); err != nil {
		return err
	}

	// Load the image file into imagick
	if err := mw.ReadImage(filepath.Join("", file.Name())); err != nil {
		return err
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_FLATTEN); err != nil {
		return err
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		return err
	}

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		return err
	}

	return nil
}

func ExtractTextFromImage(client *gosseract.Client, imageName string) (string, error) {
	defer os.Remove(imageName)

	client.SetImage(imageName)

	imgText, err := client.Text()
	if err != nil {
		return "", err
	}

	return imgText, nil
}
