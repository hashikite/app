package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	indexPath := "index.html"

	templates, err := template.ParseFiles(indexPath)
	if err != nil {
		log.Fatal(err)
	}

	kitePath := "kite.gif"

	kiteStat, err := os.Stat(kitePath)
	if err != nil {
		log.Fatal(err)
	}

	kiteTime := kiteStat.ModTime()

	file, err := os.Open(kitePath)
	if err != nil {
		log.Fatal(err)
	}

	kiteBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		feelings := r.FormValue("feelings")
		color := ""

		if feelings != "" {
			hash := md5.New()
			io.WriteString(hash, feelings)
			color = fmt.Sprintf("%.3x", hash.Sum(nil))
		}

		err = templates.ExecuteTemplate(w, "index.html", struct {
			Color    string
			Feelings string
		}{
			color,
			feelings,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(500), 500)
		}
	})

	colorRe := regexp.MustCompile("^[0-9a-f]{6}$")

	http.HandleFunc("/kite.gif", func(w http.ResponseWriter, r *http.Request) {
		newColorStr := r.URL.RawQuery
		if newColorStr == "" {
			http.ServeContent(w, r, "kite.gif", kiteTime, bytes.NewReader(kiteBytes))
		} else if colorRe.MatchString(newColorStr) {
			newColorBytes, err := hex.DecodeString(newColorStr)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, http.StatusText(500), 500)
				return
			}

			newColor := color.RGBA{newColorBytes[0], newColorBytes[1], newColorBytes[2], 255}

			image, err := gif.DecodeAll(bytes.NewReader(kiteBytes))
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				return
			}

			palette := image.Image[0].Palette
			index := palette.Index(color.White)
			palette[index] = newColor

			w.Header().Set("Content-Type", "image/gif")
			w.Header().Set("Last-Modified", kiteTime.UTC().Format(http.TimeFormat))

			err = gif.EncodeAll(w, image)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, http.StatusText(500), 500)
				return
			}
		} else {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	})

	log.Println("Serving hashy kites on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
