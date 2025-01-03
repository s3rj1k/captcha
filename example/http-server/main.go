package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"html/template"
	"image/jpeg"
	"net/http"
	"syscall"

	captcha "github.com/s3rj1k/go-captcha"
)

var (
	tmpl          *template.Template
	captchaConfig *captcha.Options
)

func main() {
	var err error

	captchaConfig, err = captcha.NewOptions()
	if err != nil {
		panic(err)
	}

	tmpl, err = template.New("captcha.html").ParseFiles("captcha.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", captchaHandle)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func captchaHandle(w http.ResponseWriter, _ *http.Request) {
	captchaObj, err := captchaConfig.CreateImage()
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer

	err = jpeg.Encode(&buff, captchaObj.Image, nil)
	if err != nil {
		panic(err)
	}

	data := struct {
		Base64 string
		Text   string
	}{
		Base64: base64.StdEncoding.EncodeToString(buff.Bytes()),
		Text:   captchaObj.Text,
	}

	if err = tmpl.Execute(w, data); err != nil {
		if errors.Is(err, syscall.EPIPE) {
			return
		}

		panic(err)
	}
}
