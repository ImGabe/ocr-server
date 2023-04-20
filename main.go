package main

import (
	"net/http"

	"github.com/imgabe/ocr-server/pkg/api"
)

func main() {
	r := api.GetRouter()

	http.ListenAndServe(":8080", r)
}
