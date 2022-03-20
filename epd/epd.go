package epd

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"golang.org/x/image/draw"
)

// Returns a new Rectangle that is resized and centered in `dst`
func fitRectInto(src *image.Rectangle, dst *image.Rectangle) image.Rectangle {
	var targetWidth int
	var targetHeight int
	var scale float64

	srcRatio := float64(src.Max.X) / float64(src.Max.Y)
	dstRatio := float64(dst.Max.X) / float64(dst.Max.Y)

	if srcRatio < dstRatio {
		// center horizontally, scale vertically
		scale = float64(dst.Max.Y) / float64(src.Max.Y)
	} else {
		// center vertically, scale horizontally
		scale = float64(dst.Max.X) / float64(src.Max.X)
	}

	targetWidth = int(float64(src.Max.X) * scale)
	targetHeight = int(float64(src.Max.Y) * scale)

	targetX := (dst.Max.X - targetWidth) / 2
	targetY := (dst.Max.Y - targetHeight) / 2

	return image.Rect(targetX, targetY, targetWidth+targetX, targetHeight+targetY)
}

// Converts a GIF, JPEG or PNG into a 16-color grayscale image. A single pixel
// is represented by 4 bits and each byte holds two pixels. The final 4-bits of
// the last byte may be discarded if width * height is odd.
func ConvertImage(input *io.Reader, width int, height int) ([]byte, error) {
	src, _, err := image.Decode(*input)
	if err != nil {
		return nil, err
	}
	dst := image.NewGray(image.Rect(0, 0, width, height))

	srcBounds := src.Bounds()
	targetRect := fitRectInto(&srcBounds, &dst.Rect)

	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	draw.CatmullRom.Scale(dst, targetRect, src, src.Bounds(), draw.Over, nil)

	// we prepend information that is required to decode the image correctly in a
  // header of ascii key-value pairs, separated by a colon and a space, one per
  // line. the header ends with two consecutive newlines.
	headerString := fmt.Sprintf("width: %d\nheight: %d\n\n", width, height) // remember to strconv.QuoteToASCII(headerString) if necessary
	header := []byte(headerString)

	// the actual conversion works by packing two nibbles together in a byte
	var result = make([]byte, (width*height+1)/2)
	for i, p := range dst.Pix {
		res := uint8((uint16(p) + 8) / 16)
		if i%2 == 0 {
			res = res << 4
		}
		result[i/2] |= res
	}

	return append(header, result...), nil
}
