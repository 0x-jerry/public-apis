package img2ascii

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"strings"

	"golang.org/x/image/draw"
)

const asciiChars = " .'`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"

const charAspect = 0.5

type Options struct {
	MaxWidth  int
	MaxHeight int
}

func Convert(r io.Reader, opts *Options) (string, error) {
	if opts == nil {
		opts = &Options{MaxWidth: 80, MaxHeight: 80}
	}
	maxW := opts.MaxWidth
	maxH := opts.MaxHeight

	img, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	iw := bounds.Dx()
	ih := bounds.Dy()

	outW := maxW
	outH := int(math.Round(float64(outW) / float64(iw) * float64(ih) * charAspect))
	if outH > maxH {
		outH = maxH
		outW = int(math.Round(float64(outH) / (float64(ih) * charAspect) * float64(iw)))
	}

	resized := image.NewRGBA(image.Rect(0, 0, outW, outH))
	draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, bounds, draw.Over, nil)

	var lines []string
	for y := 0; y < outH; y++ {
		var line strings.Builder
		for x := 0; x < outW; x++ {
			c := resized.At(x, y).(color.RGBA)
			brightness := float64(c.R+c.G+c.B) / 3.0
			idx := int(math.Floor(brightness / 255.0 * float64(len(asciiChars)-1)))
			if idx < 0 {
				idx = 0
			}
			if idx >= len(asciiChars) {
				idx = len(asciiChars) - 1
			}
			line.WriteByte(asciiChars[idx])
		}
		lines = append(lines, line.String())
	}
	return strings.Join(lines, "\n"), nil
}
