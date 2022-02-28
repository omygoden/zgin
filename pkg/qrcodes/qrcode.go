package qrcodes

import (
	"fmt"
	"github.com/omygoden/gotools/sfrand"
	"github.com/skip2/go-qrcode"
	"os"
	"time"
	"zgin/pkg/util"
)

func QrcodeGenerate(url string) (string, error) {
	t := time.Now().Format("20060102")
	fileName := sfrand.RandMd5Str() + ".jpg"
	path := fmt.Sprintf("%s/%s/%s/%s", os.Getenv("GOPATH"), util.GetProjectName(), "public/qrcode", t)

	if _, err := os.Stat(path); err != nil {
		_ = os.Mkdir(path, os.ModePerm)
	}

	imgFullPath := fmt.Sprintf("%s/%s", path, fileName)

	err := qrcode.WriteFile(url, qrcode.Medium, 256, imgFullPath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("qrcode/%s/%s", t, fileName), nil
}
