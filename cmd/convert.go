/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io"
	"log"
	"os"

	"github.com/saunaclub/inkpot-cli/epd"
	"github.com/spf13/cobra"
)

var width int
var height int
var infile string
var outfile string

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
			input, err = os.Open(args[0])
			if err != nil {
				log.Fatalf("Error opening file: %v", err)
			}
		}

		// write to given outfile, default to stdout
		var output io.Writer
		if outfile == "" || outfile == "-" {
			output = os.Stdout
		} else {
			file, err := os.OpenFile(outfile, os.O_RDWR|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("Error writing output file: %v", err)
			}
			defer file.Close()
			output = file
		}

		result, err := epd.ConvertImage(input, width, height)
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
