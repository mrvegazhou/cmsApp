package slideCaptcha

import "math"
import "encoding/json"

type SlideCaptchaPoint struct {
	X         int `json:"x,float64"`
	Y         int `json:"y,float64"`
	SecretKey string
}

func (p *SlideCaptchaPoint) SetSecretKey(secretKey string) {
	p.SecretKey = secretKey
}

func NewSlideCaptchaPoint(x int, y int) *SlideCaptchaPoint {
	return &SlideCaptchaPoint{X: x, Y: y}
}

func (p *SlideCaptchaPoint) UnmarshalJSON(data []byte) error {
	clientPoint := struct {
		X         float64 `json:"x,float64"`
		Y         float64 `json:"y,float64"`
		SecretKey string
	}{}

	if err := json.Unmarshal(data, &clientPoint); err != nil {
		return err
	}

	p.Y = int(math.Floor(clientPoint.Y))
	p.X = int(math.Floor(clientPoint.X))
	p.SecretKey = clientPoint.SecretKey
	return nil
}
