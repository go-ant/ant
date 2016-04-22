package utils

import (
	"crypto/rand"
	"os"
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

// IsFile returns true if given path is a file,
// or returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}
