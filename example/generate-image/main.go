package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	captcha "github.com/s3rj1k/go-captcha"
)

const (
	outDir         = "captcha"
	uniqueCaptchas = 420

	defaultCharsList = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func main() {
	_ = os.Mkdir(outDir, 0755)

	captchaConfig, err := captcha.NewOptions()
	if err != nil {
		panic(err)
	}

	if err = captchaConfig.SetCharacterList(defaultCharsList); err != nil {
		panic(err)
	}

	if err = captchaConfig.SetCaptchaTextLength(6); err != nil {
		panic(err)
	}

	if err = captchaConfig.SetDimensions(320, 100); err != nil {
		panic(err)
	}

	f := func() {
		captchaObj, err := captchaConfig.CreateImage()
		if err != nil {
			panic(err)
		}

		f, err := os.Create(filepath.Join(outDir, captchaObj.Text+".jpg"))
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if err = jpeg.Encode(f, captchaObj.Image, nil); err != nil {
			panic(err)
		}
	}

	i := 1
	startTime := time.Now()

	for {
		fmt.Printf("\r* Generating CAPTCHA %d.", i)

		f()

		if i >= uniqueCaptchas {
			fmt.Printf("\n* Elapsed Time: %s.\n", time.Since(startTime))

			break
		}

		i++
	}
}
