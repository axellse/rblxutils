package common

import (
	"bytes"
	"encoding/json"
	"image"
	"image/draw"
	"io"
	"net/http"

	"github.com/aarzilli/nucular"
	"github.com/axellse/rblxutils/resources"
	"github.com/nfnt/resize"
)

type ProxyStats struct {
	AvgRewriteDelayNs                     int
	AvgModifyResponseAssetDeliveryDelayNs int
	AvgModifyResponseCdnDelayNs           int
}

type SocketMessage struct {
	Type  string
	Error string //an error of "" should be interpreted as success
	DataB []byte //generic byte array that may be used for data
	Stats ProxyStats
}

func QueryAndUnmarshal(url string, resultPtr any) {
	resp, err := http.Get(url)
	if err != nil {
		FatalError(err)
	}

	ba, err := io.ReadAll(resp.Body)
	if err != nil {
		FatalError(err)
	}

	err = json.Unmarshal(ba, &resultPtr)
	if err != nil {
		FatalError(err)
	}
}

func LoadImageUI(ba []byte, width int, height int) *image.RGBA {
	rawImg, _, err := image.Decode(bytes.NewReader(ba))
	if err != nil {
		Error(err)
	}

	img := rawImg
	if width != 0 || height != 0 {
		img = resize.Resize(uint(width), uint(height), rawImg, resize.NearestNeighbor)
	}

	bounds := img.Bounds()
	rgbImg := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(rgbImg, rgbImg.Bounds(), img, bounds.Min, draw.Src)

	return rgbImg
}

func CalcWidth(win *nucular.Window) int {
	return win.WidgetBounds().W
}

func ResizeImage(img *image.RGBA, width int, height int) *image.RGBA {
	return resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor).(*image.RGBA)
}

func AutoResize(img *image.RGBA, win *nucular.Window) (*image.RGBA, int) {
	res := resize.Resize(uint(CalcWidth(win)), 0, img, resize.NearestNeighbor).(*image.RGBA)
	return res, res.Bounds().Dy()
}

var CountryCodeMap map[string]string
func InitCountryCodeMap() {
	err := json.Unmarshal(resources.CountriesJson, &CountryCodeMap)
	if err != nil {
		FatalError(err)
	}
}

func GetCountry(code string) string {
	country, ok := CountryCodeMap[code]
	if ok {
		return country
	}
	return code
}

