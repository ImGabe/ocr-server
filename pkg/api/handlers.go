package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/go-chi/render"
	"github.com/imgabe/ocr-server/pkg/types"
	"github.com/otiai10/gosseract/v2"
)

const MAX_UPLOAD_SIZE = 1024 << 13 // 8 MB

func Index(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, types.NewResponse("ping"))
}

func File(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 1MB in size", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	defer file.Close()

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	if _, err = io.Copy(tempfile, file); err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
	client.Languages = []string{"por", "eng"}

	text, err := client.Text()
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	if err := render.Render(w, r, types.NewResponse(text)); err != nil {
		render.Render(w, r, types.ErrRender(err))
		return
	}
}

func Base64(w http.ResponseWriter, r *http.Request) {
	var body = &types.Base64Request{}

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	if len(body.Base64) == 0 {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	body.Base64 = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(body.Base64, "")
	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	tempfile.Write(b)

	client := gosseract.NewClient()
	defer client.Close()

	client.Languages = []string{"por", "eng"}
	client.SetImage(tempfile.Name())

	text, err := client.Text()
	if err != nil {
		render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	if err := render.Render(w, r, types.NewResponse(text)); err != nil {
		render.Render(w, r, types.ErrRender(err))
		return
	}
}
