package text_extractor

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/otiai10/gosseract"
	"gopkg.in/gographics/imagick.v2/imagick"
)

// Dependencies structure for holding tesseract client and
// imagemagick magickwand
type Dependencies struct {
	Client    *gosseract.Client
	MagicWand *imagick.MagickWand
}

// ExtractTextFromPDF will take a fbase64 string of a pdf file and convert the file into an
// image which has its contents extracted via tesseract. It will create the image as a
// high resolution jpg file with minimal compression.
func ExtractDataFromPDF(base64PDF string, callback func(dependencies *Dependencies) (interface{}, error)) (interface{}, error) {
	dependencies := Dependencies{}

	if err := SetupDependencies(&dependencies, base64PDF); err != nil {
		return "", err
	}

	// Removing Setup Structures when method returns
	defer imagick.Terminate()
	defer dependencies.MagicWand.Destroy()

	data, err := callback(&dependencies)
	if err != nil {
		return "", err
	}

	return data, nil
}

func TextArrayFromImages(dependencies *Dependencies) (interface{}, error) {
	var imageName string
	var data []string

	// Iterate over PDF pages
	for i := 0; i < int(dependencies.MagicWand.GetNumberImages()); i++ {
		// Set Page Index
		dependencies.MagicWand.SetIteratorIndex(i)
		imageName = fmt.Sprintf("pdf_page_%v.jpg", i)
		// Save Image
		if err := dependencies.MagicWand.WriteImage(imageName); err != nil {
			return make([]string, 0), err
		}

		text, err := ExtractTextFromImage(dependencies.Client, imageName)
		if err != nil {
			return make([]string, 0), err
		}

		data = append(data, text)
	}

	return data, nil
}

func TextFromImages(dependencies *Dependencies) (interface{}, error) {
	var imageName string
	var data string = ""

	// Iterate over PDF pages
	for i := 0; i < int(dependencies.MagicWand.GetNumberImages()); i++ {
		// Set Page Index
		dependencies.MagicWand.SetIteratorIndex(i)
		imageName = fmt.Sprintf("pdf_page_%v.jpg", i)
		// Save Image
		if err := dependencies.MagicWand.WriteImage(imageName); err != nil {
			return "", err
		}

		text, err := ExtractTextFromImage(dependencies.Client, imageName)
		if err != nil {
			return "", err
		}

		data += text
	}

	return data, nil
}

func ZippedImages(dependencies *Dependencies) (interface{}, error) {
	var imageName string

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Iterate over PDF pages
	for i := 0; i < int(dependencies.MagicWand.GetNumberImages()); i++ {
		// Set Page Index
		dependencies.MagicWand.SetIteratorIndex(i)
		imageName = fmt.Sprintf("pdf_page_%v.jpg", i)
		// Save Image
		if err := dependencies.MagicWand.WriteImage(imageName); err != nil {
			return "", err
		}

		file, err := os.Open(imageName)
		if err != nil {
			return "", err
		}

		defer file.Close()

		f, err := w.Create(imageName)
		if err != nil {
			return "", err
		}

		if _, err := io.Copy(f, file); err != nil {
			return "", err
		}
	}

	err := w.Close()
	if err != nil {
		return "", err
	}
	//TODO search for memory leak
	return buf.Bytes(), nil

}

// SetupDependencies will take a Dependencies structure and populate it
func SetupDependencies(dependencies *Dependencies, base64PDF string) error {
	// Initializing Tesseract Client
	dependencies.Client = gosseract.NewClient()
	dependencies.Client.SetLanguage("por")

	imagick.Initialize()

	// creates new imagimmagick magiwand
	dependencies.MagicWand = imagick.NewMagickWand()
	//adding default config to mw image
	if err := SetupImage(base64PDF, dependencies.MagicWand); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
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
