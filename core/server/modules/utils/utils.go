package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"github.com/go-ant/ant/core/server/modules/setting"
	"github.com/go-ant/ant/core/server/modules/utils/uuid"
	"image/png"
	"os"
	"path"
	"regexp"
	"time"
)

// GetRandomString generate random string by specify chars.
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	if len(alphabets) == 0 {
		for i, b := range bytes {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		}
	} else {
		for i, b := range bytes {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}

	return string(bytes)
}

// StringInSlice check if given string in list
func StringInSlice(str string, list []string) bool {
	for _, b := range list {
		if b == str {
			return true
		}
	}
	return false
}

func SubString(str string, length int) string {
	rs := []rune(str)
	lth := len(rs)
	begin := 0
	end := begin + length
	if end > lth {
		end = lth
	}
	return string(rs[begin:end])
}

// Base64ImgUpload return image path
func Base64ImgUpload(data string) (string, error) {

	regBase64, _ := regexp.Compile("^data:[\\w]+/[\\w]+;base64,")
	strBase64 := regBase64.ReplaceAllString(data, "")
	byteBase64, _ := base64.StdEncoding.DecodeString(strBase64)

	img, err := png.Decode(bytes.NewReader(byteBase64))
	if err != nil {
		return "", err
	}

	savePath := time.Now().Format("/2006/01/")
	fileName := uuid.NewV4().String() + ".jpg"
	fullPath := setting.API.UploadFolder + savePath

	err = os.MkdirAll(fullPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	out, err := os.Create(fullPath + fileName)

	if err != nil {
		return "", err
	}
	defer out.Close()

	if err = png.Encode(out, img); err != nil {
		return "", err
	}

	return path.Join(setting.Host.Path, setting.API.FilesPath, savePath, fileName), nil
}
