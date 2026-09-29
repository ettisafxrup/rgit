// Command genlogo draws the rgit logo and writes every asset derived from it:
//
//	assets/logo.svg        vector logo for the README
//	assets/logo.png        512×512 PNG
//	assets/rgit.ico        multi-size Windows icon
//	installer/wizard*.bmp  artwork for the Inno Setup wizard
//	docs/*                 logo and favicon for the website
//
// The logo is a lowercase "r" drawn as a git branch: a line of commits with
// a branch curving away to the right, in amber on a deep navy tile.
//
// Run it from the repository root:  go run ./tools/genlogo
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

// Geometry, in a 256×256 design space.
const (
	size        = 256.0
	cornerR     = 60.0 // tile corner radius
	strokeW     = 22.0 // line thickness
	ringOuter   = 25.0 // commit node radius
	ringInner   = 11.0 // hole in the commit node
	stemX       = 88.0
	stemTop     = 80.0
	stemBottom  = 180.0
	arcCenterX  = 152.0
	arcCenterY  = 144.0
	arcRadius   = 64.0 // arc from (88,144) up and right to (152,80)
	branchEndX  = 170.0
	branchNodeY = 80.0
)

// Brand colors. The SVG in svg() uses the same values.
var (
	tileColor  = rgb(0x15, 0x30, 0x40) // deep navy  #153040
	glyphColor = rgb(0xFF, 0xC6, 0x19) // amber      #FFC619
)

type vec struct{ r, g, b float64 }

func rgb(r, g, b uint8) vec { return vec{float64(r), float64(g), float64(b)} }

// sample returns the color and opacity of the logo at a point.
func sample(x, y float64, withTile bool) (vec, float64) {
	onGlyph := glyphDistance(x, y) <= 0
	inHole := math.Hypot(x-stemX, y-stemBottom) < ringInner ||
		math.Hypot(x-branchEndX, y-branchNodeY) < ringInner

	if !withTile {
		// Glyph only (used on top of the installer artwork).
		if onGlyph && !inHole {
			return glyphColor, 1
		}
		return vec{}, 0
	}
	if roundedRectDistance(x, y) > 0 {
		return vec{}, 0
	}
	if onGlyph && !inHole {
		return glyphColor, 1
	}
	return tileColor, 1
}

// glyphDistance is negative inside the "r" glyph.
func glyphDistance(x, y float64) float64 {
	half := strokeW / 2
	stem := segmentDistance(x, y, stemX, stemTop, stemX, stemBottom) - half
	arc := arcDistance(x, y) - half
	tail := segmentDistance(x, y, arcCenterX, branchNodeY, branchEndX, branchNodeY) - half
	bottomNode := math.Hypot(x-stemX, y-stemBottom) - ringOuter
	branchNode := math.Hypot(x-branchEndX, y-branchNodeY) - ringOuter
	return math.Min(math.Min(math.Min(stem, arc), tail), math.Min(bottomNode, branchNode))
}

func segmentDistance(px, py, ax, ay, bx, by float64) float64 {
	dx, dy := bx-ax, by-ay
	t := ((px-ax)*dx + (py-ay)*dy) / (dx*dx + dy*dy)
	t = math.Max(0, math.Min(1, t))
	return math.Hypot(px-(ax+t*dx), py-(ay+t*dy))
}

// arcDistance measures the distance to the quarter circle joining the stem
// to the branch: from angle 180° (left) to 270° (top), y pointing down.
func arcDistance(px, py float64) float64 {
	dx, dy := px-arcCenterX, py-arcCenterY
	if dx <= 0 && dy <= 0 { // inside the upper-left quadrant
		return math.Abs(math.Hypot(dx, dy) - arcRadius)
	}
	a := math.Hypot(px-(arcCenterX-arcRadius), py-arcCenterY)
	b := math.Hypot(px-arcCenterX, py-(arcCenterY-arcRadius))
	return math.Min(a, b)
}

func roundedRectDistance(x, y float64) float64 {
	cx, cy := math.Abs(x-size/2)-(size/2-cornerR), math.Abs(y-size/2)-(size/2-cornerR)
	outside := math.Hypot(math.Max(cx, 0), math.Max(cy, 0))
	return outside + math.Min(math.Max(cx, cy), 0) - cornerR
}

// render draws the logo at the given pixel size with 4×4 supersampling.
func render(px int, withTile bool) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, px, px))
	const ss = 4
	scale := size / float64(px)
	for y := 0; y < px; y++ {
		for x := 0; x < px; x++ {
			var sum vec
			var alpha float64
			for sy := 0; sy < ss; sy++ {
				for sx := 0; sx < ss; sx++ {
					fx := (float64(x) + (float64(sx)+0.5)/ss) * scale
					fy := (float64(y) + (float64(sy)+0.5)/ss) * scale
					c, a := sample(fx, fy, withTile)
					sum = vec{sum.r + c.r*a, sum.g + c.g*a, sum.b + c.b*a}
					alpha += a
				}
			}
			if alpha == 0 {
				continue
			}
			img.SetNRGBA(x, y, color.NRGBA{
				R: uint8(math.Round(sum.r / alpha)),
				G: uint8(math.Round(sum.g / alpha)),
				B: uint8(math.Round(sum.b / alpha)),
				A: uint8(math.Round(alpha / (ss * ss) * 255)),
			})
		}
	}
	return img
}

func encodePNG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

// writeICO stores PNG-compressed images of several sizes in one .ico file.
func writeICO(path string, sizes []int) {
	var header, data bytes.Buffer
	binary.Write(&header, binary.LittleEndian, []uint16{0, 1, uint16(len(sizes))})
	offset := 6 + 16*len(sizes)
	for _, s := range sizes {
		img := encodePNG(render(s, true))
		dim := uint8(s % 256) // 256 is stored as 0
		header.Write([]byte{dim, dim, 0, 0})
		binary.Write(&header, binary.LittleEndian, []uint16{1, 32})
		binary.Write(&header, binary.LittleEndian, []uint32{uint32(len(img)), uint32(offset)})
		data.Write(img)
		offset += len(img)
	}
	write(path, append(header.Bytes(), data.Bytes()...))
}

// wizardImage draws installer artwork: the navy background with the glyph centered
// in the upper part (large image) or the full tile (small image).
func wizardImage(w, h int, glyphSize int, glyphTop int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := tileColor
			img.SetNRGBA(x, y, color.NRGBA{uint8(c.r), uint8(c.g), uint8(c.b), 255})
		}
	}
	glyph := render(glyphSize, false)
	left := (w - glyphSize) / 2
	for y := 0; y < glyphSize; y++ {
		for x := 0; x < glyphSize; x++ {
			g := glyph.NRGBAAt(x, y)
			if g.A == 0 {
				continue
			}
			bg := img.NRGBAAt(left+x, glyphTop+y)
			a := float64(g.A) / 255
			blend := func(fg, bg uint8) uint8 { return uint8(float64(fg)*a + float64(bg)*(1-a)) }
			img.SetNRGBA(left+x, glyphTop+y, color.NRGBA{blend(g.R, bg.R), blend(g.G, bg.G), blend(g.B, bg.B), 255})
		}
	}
	return img
}

// writeBMP saves an image as a 24-bit BMP, the format Inno Setup expects.
func writeBMP(path string, img *image.NRGBA) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	rowSize := (w*3 + 3) &^ 3
	var buf bytes.Buffer
	buf.WriteString("BM")
	binary.Write(&buf, binary.LittleEndian, []uint32{uint32(54 + rowSize*h), 0, 54})
	binary.Write(&buf, binary.LittleEndian, []uint32{40, uint32(w), uint32(h)})
	binary.Write(&buf, binary.LittleEndian, []uint16{1, 24})
	binary.Write(&buf, binary.LittleEndian, []uint32{0, uint32(rowSize * h), 2835, 2835, 0, 0})
	row := make([]byte, rowSize)
	for y := h - 1; y >= 0; y-- { // BMP rows go bottom-up
		for x := 0; x < w; x++ {
			c := img.NRGBAAt(x, y)
			row[x*3], row[x*3+1], row[x*3+2] = c.B, c.G, c.R
		}
		buf.Write(row)
	}
	write(path, buf.Bytes())
}

// svg returns the same logo as a vector image.
func svg() string {
	return fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" width="256" height="256" role="img" aria-label="rgit logo">
  <rect width="256" height="256" rx="%[1]g" fill="#153040"/>
  <g fill="none" stroke="#FFC619" stroke-width="%[2]g" stroke-linecap="round">
    <path d="M%[3]g %[4]g V%[5]g"/>
    <path d="M%[3]g %[7]g A%[8]g %[8]g 0 0 1 %[6]g %[4]g H%[9]g"/>
  </g>
  <g fill="#FFC619">
    <circle cx="%[3]g" cy="%[5]g" r="%[10]g"/>
    <circle cx="%[9]g" cy="%[4]g" r="%[10]g"/>
  </g>
  <g fill="#153040">
    <circle cx="%[3]g" cy="%[5]g" r="%[11]g"/>
    <circle cx="%[9]g" cy="%[4]g" r="%[11]g"/>
  </g>
</svg>
`, cornerR, strokeW, stemX, stemTop, stemBottom, arcCenterX, arcCenterY, arcRadius, branchEndX, ringOuter, ringInner)
}

func write(path string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote", path)
}

func main() {
	write("assets/logo.svg", []byte(svg()))
	write("assets/logo.png", encodePNG(render(512, true)))
	writeICO("assets/rgit.ico", []int{16, 24, 32, 48, 64, 128, 256})

	// The website in docs/ is published on its own, so it gets its own copies.
	write("docs/logo.svg", []byte(svg()))
	write("docs/favicon.png", encodePNG(render(64, true)))
	write("docs/og-image.png", encodePNG(render(512, true)))

	// Inno Setup picks the best match for the screen's DPI from these sizes.
	writeBMP("installer/wizard-large-100.bmp", wizardImage(164, 314, 150, 60))
	writeBMP("installer/wizard-large-200.bmp", wizardImage(328, 628, 300, 120))
	writeBMP("installer/wizard-small-100.bmp", render55(55))
	writeBMP("installer/wizard-small-200.bmp", render55(110))
}

// render55 draws the full tile on a white background for the small image.
func render55(px int) *image.NRGBA {
	tile := render(px, true)
	img := image.NewNRGBA(tile.Bounds())
	for y := 0; y < px; y++ {
		for x := 0; x < px; x++ {
			c := tile.NRGBAAt(x, y)
			a := float64(c.A) / 255
			blend := func(v uint8) uint8 { return uint8(float64(v)*a + 255*(1-a)) }
			img.SetNRGBA(x, y, color.NRGBA{blend(c.R), blend(c.G), blend(c.B), 255})
		}
	}
	return img
}
