package captcha

import (
	"image/color"
	"math/rand"
	"strings"
)

func randomString(rng *rand.Rand, length int, chars string) string {
	sb := strings.Builder{}
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(
			chars[rng.Intn(
				len(chars),
			)],
		)
	}

	return sb.String()
}

func (opts *Options) randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(opts.rng.Intn(256)),
		G: uint8(opts.rng.Intn(256)),
		B: uint8(opts.rng.Intn(256)),
		A: uint8(255),
	}
}

func (opts *Options) randomLightColor() color.Color {
	return color.RGBA{
		R: uint8(opts.rng.Intn(55) + 200),
		G: uint8(opts.rng.Intn(55) + 200),
		B: uint8(opts.rng.Intn(55) + 200),
		A: uint8(255),
	}
}

func (opts *Options) randomMiddleColor() color.Color {
	return color.RGBA{
		R: uint8(opts.rng.Intn(155) + 100),
		G: uint8(opts.rng.Intn(155) + 100),
		B: uint8(opts.rng.Intn(155) + 100),
		A: uint8(255),
	}
}

func (opts *Options) randomDarkColor() color.Color {
	return color.RGBA{
		R: uint8(opts.rng.Intn(100)),
		G: uint8(opts.rng.Intn(100)),
		B: uint8(opts.rng.Intn(100)),
		A: uint8(255),
	}
}
