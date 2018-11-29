package main

import (
	"fmt"
	"math"
	"os"

	flag "github.com/ogier/pflag"
	unicommon "github.com/unidoc/unidoc/common"
	"github.com/unidoc/unidoc/pdf/creator"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

var offsetX, offsetY float64
var scaleH, scaleW, scaleHCenter, scaleWCenter, center, verbose, version bool
var opacity, angle float64

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
	flag.Float64VarP(&opacity, "opacity", "o", 0.5, "Opacity of watermark. float between 0 to 1.")
	flag.Float64VarP(&angle, "angle", "a", 0, "Angle of rotation. between 0 to 360, counter clock-wise.")

	flag.BoolVarP(&verbose, "verbose", "v", false, "Display debug information.")
	flag.BoolVarP(&version, "version", "V", false, "Display Version information.")

	flag.Usage = func() {
		fmt.Println("markpdf <source> <watermark> <output> [options...]")
		fmt.Println("<source> and <output> should be path to a PDF file and <watermark> can be a text or path of an image.")
		fmt.Println("Example: markpdf \"path/to/083.pdf\" \"img/logo.png\" \"path/to/voucher_083.pdf\" --position=10,-10 --opacity=0.4")
		fmt.Println("Available Options: ")
		flag.PrintDefaults()
	}
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
	addImageMark(sourcePath, outputPath, watermark)

	fmt.Printf("SUCCESS: Output generated at : %s \n", outputPath)
	os.Exit(0)
}

func addImageMark(inputPath string, outputPath string, watermarkPath string) error {
	debugInfo(fmt.Sprintf("Input PDF: %v", inputPath))
	debugInfo(fmt.Sprintf("Watermark image: %s", watermarkPath))

	c := creator.New()

	watermarkImg, err := creator.NewImageFromFile(watermarkPath)
	fatalIfError(err, fmt.Sprintf("Failed to load watermark image. [%s]", err))

	// Read the input pdf file.
	f, err := os.Open(inputPath)
	fatalIfError(err, fmt.Sprintf("Failed to open the source file. [%s]", err))
	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	fatalIfError(err, fmt.Sprintf("Failed to parse the source file. [%s]", err))

	numPages, err := pdfReader.GetNumPages()
	fatalIfError(err, fmt.Sprintf("Failed to get PageCount of the source file. [%s]", err))

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		// Read the page.
		page, err := pdfReader.GetPage(pageNum)
		fatalIfError(err, fmt.Sprintf("Failed to read page from source. [%s]", err))

		// Add to creator.
		c.AddPage(page)

		// Calculate the position on first page
		if pageNum == 1 {
			debugInfo(fmt.Sprintf("Page Width       : %v", c.Context().PageWidth))
			debugInfo(fmt.Sprintf("Page Height      : %v", c.Context().PageHeight))
			debugInfo(fmt.Sprintf("Watermark Width  : %v", watermarkImg.Width()))
			debugInfo(fmt.Sprintf("Watermark Height : %v", watermarkImg.Height()))

			if scaleWCenter {
				watermarkImg.ScaleToWidth(c.Context().PageWidth)
				offsetX = 0
				offsetY = (c.Context().PageHeight - watermarkImg.Height()) / 2
			} else if scaleHCenter {
				watermarkImg.ScaleToHeight(c.Context().PageHeight)
				offsetX = (c.Context().PageWidth - watermarkImg.Width()) / 2
				offsetY = 0
			} else if scaleW {
				watermarkImg.ScaleToWidth(c.Context().PageWidth)
				offsetX = 0
			} else if scaleH {
				watermarkImg.ScaleToHeight(c.Context().PageHeight)
				offsetY = 0
			} else if center {
				offsetX = (c.Context().PageWidth - watermarkImg.Width()) / 2
				offsetY = (c.Context().PageHeight - watermarkImg.Height()) / 2
			}

			// None of the above logic is setting negative position
			// So, if found, that must be set from command line flags
			if offsetX < 0 {
				offsetX = c.Context().PageWidth - (watermarkImg.Width() + math.Abs(offsetX))
			}
			if offsetY < 0 {
				offsetY = c.Context().PageHeight - (watermarkImg.Height() + math.Abs(offsetY))
			}
		}

		watermarkImg.SetPos(offsetX, offsetY)
		watermarkImg.SetOpacity(opacity)
		watermarkImg.SetAngle(angle)

		_ = c.Draw(watermarkImg)
	}

	err = c.WriteToFile(outputPath)
	return err
}
