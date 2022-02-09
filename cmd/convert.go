/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

var width int
var height int
var outfile string

func Foobar () (io.Writer, error) {
	return os.Stdout, nil
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert <file|->",
	Short: "Convert a single file to a 4-bit, 16-color grayscale image",
	Long: `Convert a single file to a 4-bit, 16-color grayscale image.
Supports jpeg, png and gif files.

Pass "-" as the filename to read from stdin.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		// read from given file or stdin
		var input io.Reader
		if args[0] == "-" {
			input = os.Stdin
		} else {
			input, err := os.Open(args[0])
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
			defer input.Close()
		}

	    // write to given outfile, default to stdout
	    var output io.Writer
		if outfile == "" {
			output = os.Stdout
		} else {
			file, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("Error writing output file: %v", err)
			}
			defer file.Close()
			output = file
		}

		result, err := convertImage(input, width, height)
		if err != nil {
			log.Fatalf("Could not convert image: %v", err)
		}
		output.Write(result)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().IntVarP(&width, "width", "x", 540, "target width")
	convertCmd.Flags().IntVarP(&height, "height", "y", 960, "target height")
	convertCmd.Flags().StringVarP(&outfile, "output", "o", "", "file to write the result to (default stdout)")
}

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

func convertImage(input io.Reader, width int, height int) ([]byte, error) {
	src, _, err := image.Decode(input)
	if err != nil {
		return nil, err
	}
	dst := image.NewGray(image.Rect(0, 0, width, height))

	srcBounds := src.Bounds()
	targetRect := fitRectInto(&srcBounds, &dst.Rect)

	draw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	draw.CatmullRom.Scale(dst, targetRect, src, src.Bounds(), draw.Over, nil)

	// the actual conversion works by packing two nibbles together in a byte
	var result = make([]byte, (width*height+1)/2)
	for i, p := range dst.Pix {
		res := uint8((uint16(p) + 8) / 16)
		if i%2 == 0 {
			result[i/2] = (res << 4)
		} else {
			// note that integer division makes sure we're writing at the same
			// index for odd and even indices
			result[i/2] = result[i/2] | res
		}
	}

	return result, nil
}
