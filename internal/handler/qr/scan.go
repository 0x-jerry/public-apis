package qr

import (
	"encoding/json"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type scanResult struct {
	Version  int           `json:"version"`
	Data     string        `json:"data"`
	Location []interface{} `json:"location"`
}

func scan(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Not found file"})
		return
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Failed to decode image"})
		return
	}

	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "Failed to process image"})
		return
	}

	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bitmap, nil)
	if err != nil {
		writeJSON(w, http.StatusOK, scanResult{})
		return
	}

	var locs []interface{}
	for _, p := range result.GetResultPoints() {
		locs = append(locs, map[string]float64{"x": p.GetX(), "y": p.GetY()})
	}

	res := scanResult{
		Data:     result.GetText(),
		Location: locs,
	}
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
