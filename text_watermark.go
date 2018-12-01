package main

import (
	"fmt"
	"math"

	"github.com/unidoc/unidoc/pdf/creator"
	"github.com/unidoc/unidoc/pdf/model/fonts"
)

var fontList = []string{
	"courier", "courier_bold", "courier_oblique", "courier_bold_oblique",
	"helvetica", "helvetica_bold", "helvetica_oblique", "helvetica_bold_oblique",
	"times", "times_bold", "times_italic", "times_bold_italic",
}

func drawText(p *creator.Paragraph, c *creator.Creator) {
	// Change to times bold font (default is helvetica).
	p.SetFont(getFontByName(font))
	p.SetFontSize(fontSize)
	p.SetPos(offsetX, offsetY)
	p.SetColor(creator.ColorRGBFromHex("#" + color))
	p.SetAngle(angle)

	_ = c.Draw(p)
}

func adjustTextPosition(p *creator.Paragraph, c *creator.Creator) {
	p.SetTextAlignment(creator.TextAlignmentLeft)
	p.SetEnableWrap(false)

	if center {
		// To set text at middle, widthining the paragraph to page width and setting center allignment
		// then placing the box at left side of the page
		p.SetWidth(c.Context().PageWidth)
		p.SetTextAlignment(creator.TextAlignmentCenter)
		offsetX = 0
		offsetY = (c.Context().PageHeight / 2) - (p.Height() / 2)
	} else {
		if offsetX < 0 {
			// To set text from right side, widthining the paragraph to (page width - preferred offset)
			// then placing the box at left side of the page
			p.SetWidth(c.Context().PageWidth - math.Abs(offsetX))
			p.SetTextAlignment(creator.TextAlignmentRight)
			offsetX = 0
		}
		if offsetY < 0 {
			offsetY = c.Context().PageHeight - (math.Abs(offsetY) + p.Height())
		}

	}

	debugInfo(fmt.Sprintf("Paragraph width: %f", p.Width()))
	debugInfo(fmt.Sprintf("Paragrapg height: %f", p.Height()))
}

func getFontByName(fontName string) fonts.Font {
	switch fontName {
	case "courier":
		return fonts.NewFontCourier()
	case "courier_bold":
		return fonts.NewFontCourierBold()
	case "courier_oblique":
		return fonts.NewFontCourierOblique()
	case "courier_bold_oblique":
		return fonts.NewFontCourierBoldOblique()
	case "helvetica":
		return fonts.NewFontHelvetica()
	case "helvetica_bold":
		return fonts.NewFontHelveticaBold()
	case "helvetica_oblique":
		return fonts.NewFontHelveticaOblique()
	case "helvetica_bold_oblique":
		return fonts.NewFontHelveticaBoldOblique()
	case "times":
		return fonts.NewFontTimesRoman()
	case "times_bold":
		return fonts.NewFontTimesBold()
	case "times_italic":
		return fonts.NewFontTimesItalic()
	case "times_bold_italic":
		return fonts.NewFontTimesBoldItalic()
	}

	debugInfo("No allowed font name didn't match the allowed list. Using helvetica_bold as default.")
	debugInfo(fmt.Sprintf("Allowed font names: %s", fontList))

	return fonts.NewFontHelveticaBold()
}
