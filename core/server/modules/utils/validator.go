package utils

import (
	"regexp"
	"strings"
)

const (
	Regex_Email        string = "[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?"
	Regex_Numeric      string = "^[-+]?[0-9]+$"
	Regex_Alpha        string = "^[a-zA-Z]+$"
	Regex_AlphaNumeric string = "^[a-zA-Z0-9]+$"
	Regex_Int          string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Regex_Float        string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	Regex_DataURI      string = "^data:.+\\/(.+);base64$"
	Regex_URL          string = `^((ftp|http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`
	Regex_Base64       string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
)

// IsEmail check email address
func IsEmail(data string) bool {
	regPattern := regexp.MustCompile(Regex_Email)
	return regPattern.MatchString(data)
}

// IsNumber returns whether the given data is a number
func IsNumeric(data string) bool {
	regPattern := regexp.MustCompile(Regex_Numeric)
	return regPattern.MatchString(data)
}

// IsAlpha check if the string contains only letters (a-zA-Z). Empty string is valid.
func IsAlpha(str string) bool {
	if IsEmpty(str) {
		return true
	}
	regPattern := regexp.MustCompile(Regex_Alpha)
	return regPattern.MatchString(str)
}

// IsAlphanumeric check if the string contains only letters and numbers. Empty string is invalid.
func IsAlphaNumeric(str string) bool {
	if IsEmpty(str) {
		return false
	}
	regPattern := regexp.MustCompile(Regex_AlphaNumeric)
	return regPattern.MatchString(str)
}

// IsInt check if the string is an integer. Empty string is invalid.
func IsInt(str string) bool {
	if IsEmpty(str) {
		return false
	}
	regPattern := regexp.MustCompile(Regex_Int)
	return regPattern.MatchString(str)
}

// IsFloat check if the string is a float. Empty string is invalid.
func IsFloat(str string) bool {
	if IsEmpty(str) {
		return false
	}
	regPattern := regexp.MustCompile(Regex_Float)
	return regPattern.MatchString(str)
}

// IsDataURI checks if a string is base64 encoded data URI such as an image
func IsDataURI(str string) bool {
	dataURI := strings.Split(str, ",")
	regPattern := regexp.MustCompile(Regex_DataURI)
	if !regPattern.MatchString(dataURI[0]) {
		return false
	}
	return IsBase64(dataURI[1])
}

// IsURL check if the string is an URL.
func IsURL(data string) bool {
	if data == "" || len(data) >= 2083 {
		return false
	}
	regPattern := regexp.MustCompile(Regex_URL)
	return regPattern.MatchString(data)
}

// IsBase64 check if a string is base64 encoded.
func IsBase64(str string) bool {
	regPattern := regexp.MustCompile(Regex_Base64)
	return regPattern.MatchString(str)
}

// IsEmpty check if the string is empty.
func IsEmpty(data string) bool {
	return len(data) == 0
}
