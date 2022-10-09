package main

import (
	"fmt"
	"math"

	"github.com/unidoc/unidoc/pdf/creator"
)

func drawImage(watermarkImg *creator.Image, c *creator.Creator) {
	watermarkImg.SetOpacity(opacity)
	watermarkImg.SetAngle(angle)
	if tiles {
		repeatTiles(watermarkImg, c)
		return
	}
	watermarkImg.SetPos(offsetX, offsetY)
	_ = c.Draw(watermarkImg)
}

func repeatTiles(watermarkImg *creator.Image, c *creator.Creator) {
	w := watermarkImg.Width()
	h := watermarkImg.Height()
	pw := c.Context().PageWidth
	ph := c.Context().PageHeight

	nw := math.Ceil(pw / w)
	nh := math.Ceil(ph / h)

	debugInfo(fmt.Sprintf("Settings tiles of %v x %v", nw, nh))
	for i := 0; i < int(nw); i++ {
		x := w * float64(i)
		for j := 0; j < int(nh); j++ {
			y := h * float64(j)
			watermarkImg.SetPos(x, y)
			_ = c.Draw(watermarkImg)
		}
	}
}

func adjustImagePosition(watermarkImg *creator.Image, c *creator.Creator) {
	debugInfo(fmt.Sprintf("Watermark Width  : %v", watermarkImg.Width()))
	debugInfo(fmt.Sprintf("Watermark Height : %v", watermarkImg.Height()))


	if scaleImage != 100 {
		debugInfo(fmt.Sprintf("Scaling to %v", scaleImage))
		watermarkImg.ScaleToHeight(scaleImage * watermarkImg.Width() / 100)
	}
	if tiles {
		offsetX, offsetY = 0, 0
		return
	}
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
