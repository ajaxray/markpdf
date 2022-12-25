package main

import (
	"fmt"
	"github.com/unidoc/unidoc/pdf/creator"
	"math"
)

// Watermarkable is a common interface for watermarkable image or paragraph
type Watermarkable interface {
	creator.VectorDrawable
	SetPos(x, y float64)
}

func repeatTiles(watermark Watermarkable, c *creator.Creator) {
	w := watermark.Width()
	h := watermark.Height()
	pw := c.Context().PageWidth
	ph := c.Context().PageHeight

	nw := math.Ceil(pw / w)
	nh := math.Ceil(ph / h)

	debugInfo(fmt.Sprintf("Settings tiles of %v x %v", nw, nh))
	for i := 0; i < int(nw); i++ {
		x := w * float64(i)
		for j := 0; j < int(nh); j++ {
			y := h * float64(j)
			watermark.SetPos(x + spacing * float64(i), y + spacing * float64(j))
			err := c.Draw(watermark)
			if err != nil {
				fatalIfError(err, fmt.Sprintf("Error %s", err))
			}
		}
	}
}
