package utils

import "github.com/skip2/go-qrcode"

func GenerateQR(link string) ([]byte, error) {
	return qrcode.Encode(link, qrcode.High, 256)
}
