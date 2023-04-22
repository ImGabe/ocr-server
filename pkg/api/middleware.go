package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/imgabe/ocr-server/pkg/types"
)

const MAX_UPLOAD_SIZE = 1024 << 13 // 8 MB

func CheckFileSize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			render.Render(w, r, types.ErrRequestEntityTooLarge(err))
			return
		}

		next.ServeHTTP(w, r)
	})
}
