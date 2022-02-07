/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/image/draw"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		width, _ := cmd.Flags().GetInt("width")
		height, _ := cmd.Flags().GetInt("height")

		convertImage(width, height)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().Int("width", 540, "Target width")
	convertCmd.Flags().Int("height", 960, "Target height")
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

	return image.Rect(targetX, targetY, targetWidth + targetX, targetHeight + targetY)
}

func convertImage(width int, height int) {
	// for now, we're reading a predefined input file
	input, err := os.Open("cat.jpg")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer input.Close()

	// … and writing to a predefined output file
	output, err := os.OpenFile("cat_resized.png", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
	defer output.Close()

	src, _ := jpeg.Decode(input)
	dst := image.NewGray16(image.Rect(0, 0, width, height))

	srcBounds := src.Bounds()
	targetRect := fitRectInto(&srcBounds, &dst.Rect)

	fmt.Printf("targetRect: %v", targetRect)

	draw.CatmullRom.Scale(dst, targetRect, src, src.Bounds(), draw.Over, nil)

	// we're taking care of a more efficient output encoding later
	png.Encode(output, dst)
}
