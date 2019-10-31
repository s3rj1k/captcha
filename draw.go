package captcha

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func distort(obj draw.Image, amplude, period float64, backgroundColor color.Color) {
	objSize := obj.Bounds().Size()

	dx := 2.0 * math.Pi / period

	for x := 0; x < objSize.X; x++ {
		for y := 0; y < objSize.Y; y++ {
			if obj.At(x, y) == backgroundColor {
				continue
			}

			xo := amplude * math.Sin(float64(y)*dx)
			yo := amplude * math.Cos(float64(x)*dx)

			rgb := obj.At(x+int(xo), y+int(yo))

			obj.Set(x, y, rgb)
		}
	}
}

func drawFreeTypeString(obj draw.Image, text string, x, y int, fontDPI, fontSize float64, textFont *truetype.Font, textColor color.Color) error {
	ctx := freetype.NewContext()

	ctx.SetDPI(fontDPI)
	ctx.SetDst(obj)
	ctx.SetClip(obj.Bounds())
	ctx.SetSrc(image.NewUniform(textColor))
	ctx.SetFontSize(fontSize)
	ctx.SetFont(textFont)
	ctx.SetHinting(font.HintingFull)

	pt := freetype.Pt(x, y)

	if _, err := ctx.DrawString(text, pt); err != nil {
		return err
	}

	return nil
}

func (opts *Options) drawDotNoise(obj draw.Image) {
	objSize := obj.Bounds().Size()
	noiseCount := (objSize.X * objSize.Y) / (opts.dotNoise + 1)

	for i := 0; i < noiseCount; i++ {
		x := opts.randomInt(objSize.X)
		y := opts.randomInt(objSize.Y)

		if i%2 == 0 {
			x += opts.randomInt(3)
			y += opts.randomInt(3)
		}

		obj.Set(x, y, opts.randomColor())
	}
}

func (opts *Options) drawRectNoise(obj draw.Image) {
	objSize := obj.Bounds().Size()
	noiseCount := (objSize.X * objSize.Y) / (opts.rectNoise + 1)

	for i := 0; i < noiseCount/6; i++ {
		x := opts.randomInt(objSize.X)
		y := opts.randomInt(objSize.Y)

		rect := image.Rect(
			x, y,
			x+opts.randomInt(3)+2,
			y+opts.randomInt(3)+2,
		)

		draw.Draw(
			obj, rect,
			&image.Uniform{
				opts.randomMiddleColor(),
			},
			image.ZP, draw.Src,
		)
	}
}

func (opts *Options) drawBorder(obj draw.Image) {
	objSize := obj.Bounds().Size()

	for x := 0; x < objSize.X; x++ {
		obj.Set(x, 0, opts.borderColor)
		obj.Set(x, objSize.Y-1, opts.borderColor)
	}

	for y := 0; y < objSize.Y; y++ {
		obj.Set(0, y, opts.borderColor)
		obj.Set(objSize.X-1, y, opts.borderColor)
	}
}

func (opts *Options) drawTextNoise(obj draw.Image) error {
	objSize := obj.Bounds().Size()
	noiseCount := (objSize.X * objSize.Y) / (opts.textNoise + 1)
	maxFontSize := float64(objSize.Y) * opts.fontScale

	for i := 0; i < noiseCount; i++ {
		fontSize := maxFontSize/2 + float64(opts.randomInt(6)) + opts.randomFloat64()
		textFont := opts.fonts[opts.randomInt(len(opts.fonts))]
		textColor := opts.randomLightColor()
		char := opts.randomString(1, opts.characterList)

		x := opts.randomInt(objSize.X)
		y := opts.randomInt(objSize.Y)

		if err := drawFreeTypeString(
			obj, char,
			x, y, opts.fontDPI,
			fontSize, textFont,
			textColor,
		); err != nil {
			return err
		}
	}

	return nil
}

func (opts *Options) drawCaptcha(obj draw.Image, text string) error {
	objSize := obj.Bounds().Size()
	overlayImage := image.NewNRGBA(
		image.Rect(
			5, 5,
			objSize.X-5,
			objSize.Y-5,
		),
	)
	overlayImageSize := overlayImage.Bounds().Size()
	textWidth := overlayImageSize.X / len(text)
	maxFontSize := float64(overlayImageSize.Y) * opts.fontScale

	for i, char := range text {
		fontSize := maxFontSize - 2*opts.randomFloat64()
		textFont := opts.fonts[opts.randomInt(len(opts.fonts))]
		textColor := opts.randomDarkColor()

		x := int(fontSize)/4 + i*int(fontSize) + textWidth/int(maxFontSize)
		if i > 0 {
			x -= opts.randomInt(overlayImageSize.X / 64)
		}

		y := overlayImageSize.Y/2 + int(maxFontSize/3) + opts.randomInt(overlayImageSize.Y/8)

		if err := drawFreeTypeString(
			overlayImage, string(char),
			x, y, opts.fontDPI,
			fontSize, textFont,
			textColor,
		); err != nil {
			return err
		}
	}

	distort(
		overlayImage,
		10.0,
		200.0,
		opts.backgroundColor,
	)

	draw.Draw(
		obj, overlayImage.Bounds(),
		overlayImage,
		image.ZP, draw.Over,
	)

	return nil
}

func (opts *Options) drawHollowLine(obj draw.Image) {
	objSize := obj.Bounds().Size()
	begin := objSize.X / 20
	end := begin * 18

	x1 := float64(opts.randomInt(begin)) + opts.randomFloat64()
	x2 := float64(opts.randomInt(begin)+end) + opts.randomFloat64()

	multiple := (float64(opts.randomInt(4)+1) + opts.randomFloat64()) / float64(5)
	if int(multiple*10)%3 == 0 {
		multiple *= -1.0
	}

	w := objSize.Y / 20

	for ; x1 < x2; x1++ {
		y := float64(objSize.Y/2) * math.Sin(x1*math.Pi*multiple/(float64(objSize.X)+opts.randomFloat64()))

		if multiple < 0 {
			y = y + float64(objSize.Y/2) + opts.randomFloat64()
		}

		obj.Set(int(x1), int(y), opts.backgroundColor)

		for i := 0; i <= w; i++ {
			obj.Set(int(x1), int(y)+i, opts.backgroundColor)
		}
	}
}
