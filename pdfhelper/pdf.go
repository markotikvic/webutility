package pdfhelper

import (
	"fmt"
	"strings"

	"git.to-net.rs/marko.tikvic/gofpdf"
)

// Block ...
const (
	CENTER   = "C"
	MIDDLE   = "M"
	LEFT     = "L"
	RIGHT    = "R"
	TOP      = "T"
	BOTTOM   = "B"
	FULL     = "1"
	NOBORDER = ""
)

const (
	CONTINUE = 0
	NEWLINE  = 1
	BELLOW   = 2
)

const (
	cyrillicEncoding = "cp1251"
	latinEncoding    = "cp1250"
)

const threeDots = "\u2056\u2056\u2056"

type TableCell struct {
	W, H            float64
	Text            string
	Font, FontStyle string
	FontSize        float64
	Border          string
	Alignment       string
}

// Helper ...
type Helper struct {
	*gofpdf.Fpdf
	translators map[string]func(string) string
}

// New ...
func New(ori, unit, size string) *Helper {
	helper := &Helper{
		Fpdf: gofpdf.New(ori, unit, size, ""),
	}

	return helper
}

func (pdf *Helper) LoadTranslators() {
	pdf.translators = make(map[string]func(string) string)
	pdf.translators[latinEncoding] = pdf.UnicodeTranslatorFromDescriptor(latinEncoding)
	pdf.translators[cyrillicEncoding] = pdf.UnicodeTranslatorFromDescriptor(cyrillicEncoding)
}

func (pdf *Helper) DrawCell(x, y float64, c TableCell) {
	pdf.SetXY(x, y)
	pdf.SetFont(c.Font, c.FontStyle, c.FontSize)
	pdf.CellFormat(c.W, c.H, pdf.toUTF8(c.Text), c.Border, BELLOW, c.Alignment, false, 0, "")
}

func (pdf *Helper) DrawColumn(x, y float64, cells []TableCell) {
	pdf.SetXY(x, y)
	for _, c := range cells {
		pdf.SetFont(c.Font, c.FontStyle, c.FontSize)
		pdf.CellFormat(c.W, c.H, pdf.toUTF8(c.Text), c.Border, BELLOW, c.Alignment, false, 0, "")
	}
}

func (pdf *Helper) DrawRow(x, y float64, cells []TableCell) {
	pdf.SetXY(x, y)
	for _, c := range cells {
		pdf.SetFont(c.Font, c.FontStyle, c.FontSize)
		pdf.CellFormat(c.W, c.H, pdf.toUTF8(c.Text), c.Border, CONTINUE, c.Alignment, false, 0, "")
	}
}

func (pdf *Helper) TextLength(txt, family, style string, size float64) float64 {
	family, _, _, _ = pdf.setCorrectFontFamily(textEncoding(txt))
	return pdf.Fpdf.TextLength(txt, family, style, size)
}

func (pdf *Helper) LimitText(text, limiter string, maxWidth float64) string {
	parts := pdf.Fpdf.SplitText(text, maxWidth)
	if len(parts) > 1 {
		return parts[0] + limiter
	}

	return text
}

// InsertImage ...
func (pdf *Helper) InsertImage(img string, x, y, w, h float64) {
	imgType := ""
	if parts := strings.Split(img, "."); len(parts) >= 2 {
		imgType = parts[len(parts)-1]
	}
	opt := gofpdf.ImageOptions{
		ImageType:             imgType,
		ReadDpi:               false,
		AllowNegativePosition: false,
	}
	autoBreak := false // if it's not false then you can't draw the image at an arbitrary height (y position)
	pdf.ImageOptions(img, x, y, w, h, autoBreak, opt, 0, "")
}

func (pdf *Helper) PageHasSpace(requiredHeight float64) bool {
	_, h := pdf.GetPageSize()
	_, _, _, bot := pdf.GetMargins()
	return (h - bot - pdf.GetY()) > requiredHeight
}

// DrawBox ...
func (pdf Helper) DrawBox(x0, y0, w, h float64) {
	pdf.Line(x0, y0, x0+w, y0)
	pdf.Line(x0+w, y0, x0+w, y0+h)
	pdf.Line(x0+w, y0+h, x0, y0+h)
	pdf.Line(x0, y0+h, x0, y0)
}

// Strana %d/{TotalPages}
func (pdf *Helper) InsertPageNumber(x, y float64, format string) {
	num := fmt.Sprintf(format, pdf.PageNo())
	pdf.DrawColumn(x, y, []TableCell{{10, 1, num, "DejaVuSans", "", 8, NOBORDER, LEFT}})
}

func (pdf *Helper) SuperscriptText(x, y, cw, ch float64, text, script string) {
	family, style, size, sizeU := pdf.setCorrectFontFamily(textEncoding(text))

	pdf.DrawCell(x, y, TableCell{cw, ch, text, family, style, size, NOBORDER, LEFT})

	sx := x + pdf.TextLength(text, family, style, size)
	sy := y - sizeU*0.2
	pdf.DrawCell(sx, sy, TableCell{cw, ch, script, family, style, size - 2, NOBORDER, LEFT})
}

// toUTF8 ...
func (pdf *Helper) toUTF8(s string) string {
	encoding := textEncoding(s)
	pdf.setCorrectFontFamily(encoding)
	translator, ok := pdf.translators[encoding]
	if !ok {
		return ""
	}
	return translator(s)
}

func textEncoding(s string) string {
	encoding := latinEncoding
	runes := []rune(s)
	for _, r := range runes {
		if uint64(r) >= 0x0402 && uint64(r) <= 0x044f {
			encoding = cyrillicEncoding
			break
		}
	}
	return encoding
}

func (pdf *Helper) setCorrectFontFamily(enc string) (family, style string, ptSize, unitSize float64) {
	family, style, ptSize, unitSize = pdf.GetFontInfo()
	if enc == cyrillicEncoding {
		if !strings.HasSuffix(family, "cyrillic") {
			family += "cyrillic"
		}
	} else {
		if strings.HasSuffix(family, "cyrillic") {
			family = strings.TrimSuffix(family, "cyrillic")
		}
	}
	pdf.SetFont(family, style, ptSize)
	return family, style, ptSize, unitSize
}
