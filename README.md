# markpdf - Watermark PDF files using image or text

A tiny command line tool for watermarking PDF files using image or text. 
With simple options to configure position, opacity, rotation, stretch etc.

Highlights -

- Very simple and easy to use 
- Extreamly fast!
- Stretching watermark image to height or weight proportionately 
- Options to adjust position, opacity, rotation of image
- Free and open source

## Install

It's just a single binary file, no external dependencies. 
Just download the appropriate version of [executable from latest release](https://github.com/ajaxray/markpdf/releases) for your OS. Then rename and give it execute permission.
```bash
mv markpdf_linux-amd64 markpdf  
sudo chmod +x markpdf
```

If you want to install it globally (run from any directory of your system), put it in your systems $PATH directory.
```bash
sudo mv markpdf /usr/local/bin/markpdf
```
Done! 

## How to use

### Image watermarking

Command options are shown in both, shorthand and full name.

```bash
# watermark with all default options (on top left corner with 50% opacity)
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf"

# watermark at center
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --center
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" -c

# watermark at right top with 20px offset from edge and full opaque
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --offset-x=-20 --offset-y=20 --opacity=1.0
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" -x -20 -y 20 -o 1.0

# watermark at left bottom with 100px offset and 45 degree rotation
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --offset-x=100 --offset-y=-100 --angle=45
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" -x 100 -y -100 -a 45

# stretch full with of page at page middle, with 30% opacity
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --scale-width-center --opacity=0.3
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" -Wo 0.3
# Note the capital "W" 

# stretch full with of page at page bottom
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --scale-width --offset-y=-10
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" -wy -10

# Scale the image to desired percentage
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --scale=30

# Add image as tiles all over the page
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --tiles

# Add image as tiles with interleaved spacing
markpdf "path/to/source.pdf" "img/logo.png" "path/to/output.pdf" --tiles --spacing=20
``


### Text watermarking

```bash
# watermark text at right top with 20px offset from edge
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" --offset-x=-20 --offset-y=20
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" -x -20 -y 20

# Place text at center with bold-italic "Times Roman" font in blue color
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" --center --font=times_bold_italic --color=0000FF
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" -cf times_bold_italic -l 0000FF

# Place text at center with large bold-italic "Times Roman" font in blue color
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" --center --font=times_bold_italic --font-size=24.0 --color=0000FF
markpdf "path/to/source.pdf" "The Company Name" "path/to/output.pdf" -cf times_bold_italic -s 24.0 -l 0000FF
```

#### Using placeholders with text watermark

The following placeholder can be used in text watermark:
- `{{.Page}}` prints the current page number
- `{{.Pages}}` prints the total page numbers
- `{{.Filename}}` prints name of the source file

```bash
# Using placeholders in text watermark
markpdf "path/to/083.pdf" "File: {{.Filename}} Page {{.Page}} of {{.Pages}}" "path/to/voucher_083.pdf" --position=10,-10
```

_Note: This (placeholder) feature will be available in upcoming release. If you want to use it right now, please build from the `master` branch._

#### Allowed font identifiers 

Currently the following font names are supported:
- **Courier**:	`courier`, `courier_bold`, `courier_oblique`, `courier_bold_oblique`
- **Helvetica**:	`helvetica`, `helvetica_bold`, `helvetica_oblique`, `helvetica_bold_oblique`
- **Times Roman**:	`times`, `times_bold`, `times_italic`, `times_bold_italic`

### Additional notes

- **Specifying Colors**: write them as 6 or 3 digit hexadecilal as used in CSS, without the #

- `--color`, `--font` and `--font-size` flag has no impact for Image watermarking
- `--scale-*`, `--tiles` and `--opacity` flag has no impact for Text watermarking
- Negative offset will set content positioning from opposite side (right for offsetX and botom from offsetY)
- Text with opacity is not supported at this moment. Instead, you can [create a transperent background PNG image](http://www.picturetopeople.org/text_generator/others/transparent/transparent-text-generator.html) with your text and then use it for watermarking.

## Roadmap

✅ Draw image on every page of PDF  
✅ Configure Opacity option  
✅ Configure watermark position by X and Y offset  
✅ Allow negative values to for offset to adjust from opposite direction  
✅️ Easy option for positioning image at center  
✅ Configure image rotation angle  
✅ Options to Stretch watermark to page width or height, proportionately  
✅ Options to Stretch watermark to page width or height at the middle of page  
✅ Tile Image all over the page  
✅ Render text on every page  
✅ Configure text color, style and font  
◻️ Configure text opacity  
✅ Configure text rotation angle  
✅ Text placement by offset  
✅ Put text at page center  

### Contribute

If you fix a bug or want to add/improve a feature, 
and it's alligned with the focus of this tool - _watermarking PDF with ease_, 
I will be glad to accept your PR. :) 

### Thanks

This tool was made using the beautiful [Unidoc](https://unidoc.io/) library. Thanks and ❤️ to **Unidoc**.

---
> "This is the Book about which there is no doubt, a guidance for those conscious of Allah" - [Al-Quran](http://quran.com)


