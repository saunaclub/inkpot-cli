# `inkpot-cli`

This project contains the source code to the inkpot cli, a tool that can be used to convert JPEG, PNG and GIF images to 4-bit, 16-color grayscale images. The image format is designed to be simple to decode.

## Usage

```
# inkpot help
A command-line-tool that can be used to prepare images by
resizing them, rotating them and converting them to 16-color
grayscale so they can be comfortably displayed on e-ink
displays powered by the epdiy driver.

Usage:
  inkpot-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  convert     Convert a single file to a 4-bit, 16-color grayscale image
  help        Help about any command
  serve       Run a webserver to convert images via HTTP

Flags:
  -h, --help   help for inkpot-cli

Use "inkpot-cli [command] --help" for more information about a command.
```

## Development

The package is written in [go](https://go.dev/) and you can use all the standard go tooling to run and build the binary. The individual commands are defined in the `cmd/` folder, the encoding / decoding bits are in `epd/`.

If you are using the [nix package manager](https://nixos.org/), there is a description of the build and dev environment in `default.nix`, `shell.nix` and `flake.nix` / `flake.lock`. It is recommended to use [https://direnv.net/](direnv) with [flakes support](https://github.com/nix-community/nix-direnv#flakes-support).

## Image Format (in progress)

An image blob starts with a header containing ascii-encoded key-value pairs. It holds meta-information (such as an image's `width` `w` and `height` `h`), where case-insensitive keys are separated from values using a space and a colon character (`: `). A value ends with a newline character (`\n`). The end of the header data is marked by two consecutive newline characters (`\n\n`).

The actual image data is following the header as a sequence of nibbles, or groups of four bits. Each nibble determines the grayscale value of a single pixel, where the first `w` nibbles designate the grayscale values of the first row of pixels (0 is black and 16 is white from left to right, the group of nibbles from `w + 1` to `2w` the second row of pixels from left to right and so on.
