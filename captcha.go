package captcha

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"

	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/gobolditalic"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/gofont/gomedium"
	"golang.org/x/image/font/gofont/gomediumitalic"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/gofont/gomonobold"
	"golang.org/x/image/font/gofont/gomonobolditalic"
	"golang.org/x/image/font/gofont/gomonoitalic"
	"golang.org/x/image/font/gofont/goregular"
	//	"golang.org/x/image/font/gofont/gosmallcaps"
	//	"golang.org/x/image/font/gofont/gosmallcapsitalic"
)

// DefaultCharsList defines default list of chars for a captcha image.
const DefaultCharsList = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Options manage captcha generation details.
type Options struct {
	backgroundColor color.Color
	borderColor     color.Color

	characterList string

	width  int
	height int
	length int

	textNoise int
	dotNoise  int
	rectNoise int

	fontDPI   float64
	fontScale float64

	fonts []*truetype.Font

	rng *rand.Rand

	mu *sync.RWMutex
}

// SetBackgroundColor sets captcha image's background color.
func (opts *Options) SetBackgroundColor(backgroundColor color.Color) {
	opts.backgroundColor = backgroundColor
}

// SetBorderColor sets captcha image's border color.
func (opts *Options) SetBorderColor(borderColor color.Color) {
	opts.borderColor = borderColor
}

// SetCharacterList sets available characters for a captcha image.
func (opts *Options) SetCharacterList(chars string) error {
	if len(chars) == 0 {
		return errors.New("empty character list")
	}

	opts.characterList = chars

	return nil
}

// SetCaptchaTextLength sets amount of characters, text lentgh, for a captcha image.
func (opts *Options) SetCaptchaTextLength(length int) error {
	if length <= 0 {
		return errors.New("captcha length must be greater than zero")
	}

	opts.length = length

	return nil
}

// SetFontDPI sets DPI (dots per inch) for a font.
func (opts *Options) SetFontDPI(dpi float64) error {
	if dpi < 25.0 || dpi > 300.0 {
		return errors.New("font DPI must be between 25.0 and 300.0")
	}

	opts.fontDPI = dpi

	return nil
}

// SetFontScale sets scale of a font.
func (opts *Options) SetFontScale(scale float64) error {
	if scale < 0.1 || scale > 5.0 {
		return errors.New("font scale must be between 0.1 and 5.0")
	}

	opts.fontScale = scale

	return nil
}

// SetNoiseDensity sets density between noise elements.
func (opts *Options) SetNoiseDensity(dot, rect, text float64) {
	if dot <= 0 {
		opts.dotNoise = math.MaxInt32
	} else {
		opts.dotNoise = int(1 / dot)
	}

	if rect <= 0 {
		opts.rectNoise = math.MaxInt32
	} else {
		opts.rectNoise = int(1 / rect)
	}

	if text <= 0 {
		opts.textNoise = math.MaxInt32
	} else {
		opts.textNoise = int(30 / text)
	}
}

// SetDimensions sets captcha image widtg and height.
func (opts *Options) SetDimensions(width, height int) error {
	if width <= 1 || height <= 1 {
		return errors.New("captcha width and/or height must be greater than 1px")
	}

	opts.width = width
	opts.height = height

	return nil
}

// SetFontsFromData sets font from bytes slice.
func (opts *Options) SetFontsFromData(fontsData ...[]byte) error {
	opts.fonts = make([]*truetype.Font, 0, len(fontsData))

	for _, fontData := range fontsData {
		fontTTF, err := freetype.ParseFont(fontData)
		if err != nil {
			return err
		}

		opts.fonts = append(opts.fonts, fontTTF)
	}

	return nil
}

// SetFontsFromPath sets font from font file.
func (opts *Options) SetFontsFromPath(paths ...string) error {
	fonts := make([][]byte, 0, len(paths))

	for _, path := range paths {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		fonts = append(fonts, data)
	}

	return opts.SetFontsFromData(fonts...)
}

// NewOptions creates new CAPTCHA options object with default values.
func NewOptions() (*Options, error) {
	opts := new(Options)

	// populate RNG seed
	opts.rng = rand.New(
		rand.NewSource(
			time.Now().UnixNano(),
		),
	)

	opts.mu = new(sync.RWMutex)

	opts.SetBackgroundColor(color.White)
	opts.SetBorderColor(color.Black)
	opts.SetNoiseDensity(0.05, 0.05, 0.05)

	if err := opts.SetDimensions(430, 100); err != nil {
		return nil, err
	}

	if err := opts.SetFontScale(0.6); err != nil {
		return nil, err
	}

	if err := opts.SetFontDPI(72.0); err != nil {
		return nil, err
	}

	if err := opts.SetCaptchaTextLength(8); err != nil {
		return nil, err
	}

	if err := opts.SetCharacterList(DefaultCharsList); err != nil {
		return nil, err
	}

	if err := opts.SetFontsFromData(
		gobold.TTF,
		gobolditalic.TTF,
		goitalic.TTF,
		gomedium.TTF,
		gomediumitalic.TTF,
		gomono.TTF,
		gomonobold.TTF,
		gomonobolditalic.TTF,
		gomonoitalic.TTF,
		goregular.TTF,
		//		gosmallcaps.TTF,
		//		gosmallcapsitalic.TTF,
	); err != nil {
		return nil, err
	}

	return opts, nil
}

// Captcha represents captcha object.
type Captcha struct {
	Text  string
	Image draw.Image
}

// CreateImage creates new captcha image with specified text.
func (opts *Options) CreateImage() (*Captcha, error) {
	var err error

	out := new(Captcha)
	out.Text = opts.randomString(opts.length, opts.characterList)

	out.Image = image.NewNRGBA(
		image.Rect(
			0, 0,
			opts.width,
			opts.height,
		),
	)

	draw.Draw(
		out.Image, out.Image.Bounds(),
		&image.Uniform{opts.backgroundColor},
		image.ZP, draw.Src,
	)

	opts.drawDotNoise(out.Image)

	err = opts.drawTextNoise(out.Image)
	if err != nil {
		return nil, err
	}

	err = opts.drawCaptcha(out.Image, out.Text)
	if err != nil {
		return nil, err
	}

	opts.drawRectNoise(out.Image)
	opts.drawHollowLine(out.Image)
	opts.drawBorder(out.Image)

	return out, nil
}
