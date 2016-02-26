package utils

import (
	"regexp"
	"strings"
)

var tagChecker = regexp.MustCompile("<.*?>")
var whitespaceChecker = regexp.MustCompile("\\s{2,}")

func StripTagsFromHtml(input string) string {
	output := tagChecker.ReplaceAllString(input, "")
	output = strings.Replace(output, "\n", " ", -1)
	output = strings.Replace(output, "\t", " ", -1)
	output = whitespaceChecker.ReplaceAllString(output, " ")
	return output
}
