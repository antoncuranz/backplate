package img

import (
	"errors"
	"github.com/makeworld-the-better-one/dither/v2"
	"github.com/nfnt/resize"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var UniformPalette = []color.Color{
	color.Gray16{Y: 0x0000},
	color.Gray16{Y: 0x2492},
	color.Gray16{Y: 0x4924},
	color.Gray16{Y: 0x6DB6},
	color.Gray16{Y: 0x9248},
	color.Gray16{Y: 0xB6DA},
	color.Gray16{Y: 0xDB6C},
	color.Gray16{Y: 0xFFFE},
}

var ActualPalette = []color.Color{
	color.Gray16{Y: 0x0000},
	color.Gray16{Y: 0x4b4b},
	color.Gray16{Y: 0x5555},
	color.Gray16{Y: 0x6b6b},
	color.Gray16{Y: 0x8282},
	color.Gray16{Y: 0x9696},
	color.Gray16{Y: 0xbaba},
	color.Gray16{Y: 0xd2d2},
}

func Convert(srcPath string, dstPath string) error {
	imageinput, err1 := os.Open(srcPath)
	defer imageinput.Close()

	outfile, err2 := os.Create(dstPath)
	defer outfile.Close()

	src, _, err3 := image.Decode(imageinput)

	if err := errors.Join(err1, err2, err3); err != nil {
		return err
	}

	resized := resize.Resize(820, 1200, src, resize.Lanczos3)

	gray := image.NewGray(resized.Bounds())
	draw.Draw(gray, gray.Bounds(), resized, resized.Bounds().Min, draw.Src)

	d := dither.NewDitherer(ActualPalette)
	d.Matrix = dither.FloydSteinberg
	dithered := d.DitherPaletted(gray)
	dithered.Palette = UniformPalette

	return bmp.Encode(outfile, dithered)
}
