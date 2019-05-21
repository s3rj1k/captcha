package captcha

import (
	"image/color"
	"strings"
)

func (opts *Options) randomString(length int, chars string) string {
	sb := strings.Builder{}
	sb.Grow(length)

	for i := 0; i < length; i++ {
		sb.WriteByte(
			chars[opts.randomInt(
				len(chars),
			)],
		)
	}

	return sb.String()
}

func (opts *Options) randomInt(n int) int {
	opts.mu.Lock()
	out := opts.rng.Intn(n)
	opts.mu.Unlock()
	return out
}

func (opts *Options) randomFloat64() float64 {
	opts.mu.Lock()
	out := opts.rng.Float64()
	opts.mu.Unlock()
	return out
}

func (opts *Options) randomColor() color.RGBA {
	return color.RGBA{
		R: uint8(opts.randomInt(256)),
		G: uint8(opts.randomInt(256)),
		B: uint8(opts.randomInt(256)),
		A: uint8(255),
	}
}

func (opts *Options) randomLightColor() color.Color {
	return color.RGBA{
		R: uint8(opts.randomInt(55) + 200),
		G: uint8(opts.randomInt(55) + 200),
		B: uint8(opts.randomInt(55) + 200),
		A: uint8(255),
	}
}

func (opts *Options) randomMiddleColor() color.Color {
	return color.RGBA{
		R: uint8(opts.randomInt(155) + 100),
		G: uint8(opts.randomInt(155) + 100),
		B: uint8(opts.randomInt(155) + 100),
		A: uint8(255),
	}
}

func (opts *Options) randomDarkColor() color.Color {
	return color.RGBA{
		R: uint8(opts.randomInt(100)),
		G: uint8(opts.randomInt(100)),
		B: uint8(opts.randomInt(100)),
		A: uint8(255),
	}
}
