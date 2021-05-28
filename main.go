package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	flag "github.com/ogier/pflag"
	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var offsetX, offsetY, fontSize float64
var scaleH, scaleW, scaleHCenter, scaleWCenter, center, verbose, version bool
var opacity, angle float64
var font, color string

const (
	VERSION = "1.0.0"
)

func init() {
	flag.Float64VarP(&offsetX, "offset-x", "x", 0, "Offset from left (or right for negative number).")
	flag.Float64VarP(&offsetY, "offset-y", "y", 0, "Offset from top (or bottom for negative number).")
	flag.BoolVarP(&center, "center", "c", false, "Set position at page center. Offset X and Y will be ignored.")
	flag.BoolVarP(&scaleW, "scale-width", "w", false, "Scale Image to page width. If set, offset X will be ignored.")
	flag.BoolVarP(&scaleH, "scale-height", "h", false, "Scale Image to page height. If set, top offset Y will be ignored.")
	flag.BoolVarP(&scaleWCenter, "scale-width-center", "W", false, "Scale Image to page width and Y will be set at middle.")
	flag.BoolVarP(&scaleHCenter, "scale-height-center", "H", false, "Scale Image to page height and X will be set at middle.")

	flag.StringVarP(&font, "font", "f", "helvetica_bold", "Set font. Check --help for allowed name list.")
	flag.StringVarP(&color, "color", "l", "333333", "Set font color as 6 or 3 character hexadecimal code (without '#'). See https://www.color-hex.com")
	flag.Float64VarP(&fontSize, "font-size", "s", 18.0, "Font-size in points.")

	flag.Float64VarP(&opacity, "opacity", "o", 0.5, "Opacity of watermark. float between 0 to 1.")
	flag.Float64VarP(&angle, "angle", "a", 0, "Angle of rotation. between 0 to 360, counter clock-wise.")

	flag.BoolVarP(&verbose, "verbose", "v", false, "Display debug information.")
	flag.BoolVarP(&version, "version", "V", false, "Display Version information.")

	flag.Usage = func() {
		fmt.Println("markpdf <source> <watermark> <output> [options...]")
		fmt.Println("<source> and <output> should be path to a PDF file and <watermark> can be a text or path of an image.")
		fmt.Println("text <watermark> can be used with the following variable:")
		fmt.Println("{{.Page}} current page number")
		fmt.Println("{{.Pages}} total page numbers")
		fmt.Println("{{.Filename}} source file name")
		fmt.Println("Example: markpdf \"path/to/083.pdf\" \"img/logo.png\" \"path/to/voucher_083.pdf\" --position=10,-10 --opacity=0.4")
		fmt.Println("Example: markpdf \"path/to/083.pdf\" \"File: {{.Filename}} Page {{.Page}} of {{.Pages}}\" \"path/to/voucher_083.pdf\" --position=10,-10 --opacity=0.4")
		fmt.Println("Available Options: ")
		flag.PrintDefaults()
	}
}

type Recipient struct {
	Page, Pages int
	Filename    string
}

func main() {
	flag.Parse()
	if version {
		fmt.Println("markpdf version ", VERSION)
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 3 {
		flag.Usage()
		os.Exit(0)
	}

	if verbose {
		unicommon.SetLogger(unicommon.NewConsoleLogger(unicommon.LogLevelDebug))
	}

	sourcePath := args[0]
	watermark := args[1]
	outputPath := args[2]
	markPDF(sourcePath, outputPath, watermark)

	fmt.Printf("SUCCESS: Output generated at : %s \n", outputPath)
	os.Exit(0)
}

func markPDF(inputPath string, outputPath string, watermark string) error {
	debugInfo(fmt.Sprintf("Input PDF: %v", inputPath))

	c := creator.New()
	var watermarkImg *creator.Image
	var para *creator.Paragraph

	isImageMark := isImageMark(watermark)

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	fatalIfError(err, fmt.Sprintf("Failed to open the source file. [%s]", err))
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	fatalIfError(err, fmt.Sprintf("Failed to parse the source file. [%s]", err))

	numPages, err := pdfReader.GetNumPages()
	fatalIfError(err, fmt.Sprintf("Failed to get PageCount of the source file. [%s]", err))

	// Prepare data to insert into the template.
	rec := Recipient{
		Pages:    numPages,
		Filename: filepath.Base(inputPath),
	}

	// Create a new template and parse the watermark into it.
	t := template.Must(template.New("watermark").Parse(watermark))

	for i := 0; i < numPages; i++ {
		pageNum := i + 1
		rec.Page = pageNum

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		fatalIfError(err, fmt.Sprintf("Failed to read page from source. [%s]", err))

		// Add to creator.
		c.AddPage(page)

		// Calculate the position on first page
		if pageNum == 1 {
			debugInfo(fmt.Sprintf("Page Width       : %v", c.Context().PageWidth))
			debugInfo(fmt.Sprintf("Page Height      : %v", c.Context().PageHeight))
		}

		if isImageMark {
			if pageNum == 1 {
				watermarkImg, err = creator.NewImageFromFile(watermark)
				fatalIfError(err, fmt.Sprintf("Failed to load watermark image. [%s]", err))
				adjustImagePosition(watermarkImg, c)
			}

			drawImage(watermarkImg, c)

		} else {

			// Execute the template for each page.
			buf := new(bytes.Buffer)
			err := t.Execute(buf, rec)
			fatalIfError(err, fmt.Sprintf("Failed to execute watermark template: [%s]", err))

			para = creator.NewParagraph(buf.String())
			adjustTextPosition(para, c)

			drawText(para, c)
		}

	}

	err = c.WriteToFile(outputPath)
	return err
}

func isImageMark(watermark string) bool {
	_, err := os.Stat(watermark)
	if err == nil {
		debugInfo(fmt.Sprintf("Watermark Image: %s", watermark))
		return true
	} else if os.IsNotExist(err) {
		debugInfo(fmt.Sprintf("No file exists at: %s, assuming Text Watermark.", watermark))
	} else {
		fmt.Printf("ERROR: File %s stat error: %v", watermark, err)
		os.Exit(1)
	}

	return false
}
