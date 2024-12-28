package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"image/jpeg"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/blake2s"

	captcha "github.com/s3rj1k/go-captcha"
)

// https://socketloop.com/tutorials/golang-saving-and-reading-file-with-gob

const (
	outFile        = "captcha.gob"
	uniqueCaptchas = 420

	defaultCharsList = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getStringHash(text ...string) string {
	b := []byte(strings.Join(text, ""))
	h := blake2s.Sum256(b)

	return fmt.Sprintf("%x", h)
}

// Data contains pregenerated CAPTCHAs.
type Data struct {
	Map  map[string]string
	Keys []string
}

func main() {
	data := Data{
		Map:  make(map[string]string),
		Keys: []string{},
	}

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
		var captchaObj *captcha.Captcha

		captchaObj, err = captchaConfig.CreateImage()
		if err != nil {
			panic(err)
		}

		var buff bytes.Buffer

		if err = jpeg.Encode(&buff, captchaObj.Image, nil); err != nil {
			panic(err)
		}

		data.Map[getStringHash(captchaObj.Text)] = base64.StdEncoding.EncodeToString(buff.Bytes())
	}

	startTime := time.Now()

	for {
		fmt.Printf("\r* Unique CAPTCHAs Generated: %d.", len(data.Map))

		if len(data.Map) == uniqueCaptchas {
			fmt.Printf("\n* Elapsed Time: %s.\n", time.Since(startTime))

			break
		}

		f()
	}

	data.Keys = make([]string, 0, len(data.Map))

	fmt.Printf("* Processing Keys.\n")

	for k := range data.Map {
		data.Keys = append(data.Keys, k)
	}

	fmt.Printf("* Creating GOB File.\n")

	file, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	if err = gob.NewEncoder(file).Encode(data); err != nil {
		panic(err)
	}
}
