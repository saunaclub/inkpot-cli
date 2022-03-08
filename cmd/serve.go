/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saunaclub/inkpot-cli/epd"
	"github.com/spf13/cobra"
)

var defaultWidth int = 540
var defaultHeight int = 960
var port int

type Params struct {
	Width  int    `form:"height"`
	Height int    `form:"width"`
	Url    string `form:"url"`
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run a webserver to convert images via HTTP",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		// Set a lower memory limit for multipart forms (default is 32 MiB)
		router.MaxMultipartMemory = 8 << 20 // 8 MiB
		router.GET("/", getIndex)
		router.GET("/convert", getConvert)
		router.POST("/convert", postConvert)
		router.Run(fmt.Sprintf(":%d", port))
	},
}

func getIndex(c *gin.Context) {
	usage := `# inkpot-convert

A webserver to convert GIFs, PNGs and JPEGs to 4-bit grayscale images.

## Routes

- POST /convert can be used to convert a file on your filesystem
- GET /convert/[url] can be used to convert a file publically accessible via URL

Both routes accept a "width" and a "height" parameter to configure the output size.

## Examples

Via curl:

curl -X POST http://localhost:8080/convert \
  -F "file=@my_cat.jpeg" \
  -H "Content-Type: multipart/form-data"

Via httpie:

http --form POST :8080/convert file@my_cat.jpg
`

	c.String(http.StatusOK, usage)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getConvert(c *gin.Context) {
	var params Params
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.Error(err)
	}

	if params.Url == "" {
		c.String(http.StatusBadRequest, "Please supply an image URL.")
		return
	}

	width := min(params.Width, 2000)
	height := min(params.Height, 2000)
	if width <= 0 {
		width = defaultWidth
	}
	if height <= 0 {
		height = defaultHeight
	}

	response, err := http.Get(params.Url)
	if err != nil || response.StatusCode != http.StatusOK {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	reader := response.Body
	defer reader.Close()

	converted, err := epd.ConvertImage(reader, width, height)
	if err != nil {
		c.Error(err)
	}

	c.Header("X-Image-Width", fmt.Sprintf("%d", width))
	c.Header("X-Image-Height", fmt.Sprintf("%d", height))
	c.Data(http.StatusOK, "x-image/inkpot-epd", converted)
}

func postConvert(c *gin.Context) {
	var params Params
	err := c.ShouldBindQuery(&params)
	if err != nil {
		c.Error(err)
	}

	width := min(params.Width, 2000)
	height := min(params.Height, 2000)
	if width <= 0 {
		width = defaultWidth
	}
	if height <= 0 {
		height = defaultHeight
	}

	// Single file
	file, err := c.FormFile("file")
	if err != nil {
		c.Error(err)
	}

	reader, err := file.Open()
	if err != nil {
		c.Error(err)
	}
	defer reader.Close()

	converted, err := epd.ConvertImage(reader, width, height)
	if err != nil {
		c.Error(err)
	}

	c.Header("X-Image-Width", fmt.Sprintf("%d", width))
	c.Header("X-Image-Height", fmt.Sprintf("%d", height))
	c.Data(http.StatusOK, "x-image/inkpot-epd", converted)
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "port to bind to")
	// gin.SetMode(gin.ReleaseMode)
}
