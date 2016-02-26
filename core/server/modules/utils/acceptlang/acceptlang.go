// https://github.com/martini-contrib/acceptlang
package acceptlang

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const (
	acceptLanguageHeader = "Accept-Language"
)

// A single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
	Language string
	Quality  float32
}

// A slice of sortable AcceptLanguage instances.
type AcceptLanguages []AcceptLanguage

// Returns the language with the highest quality score
func (al AcceptLanguages) Best() AcceptLanguage {
	return al[0]
}

// Returns the total number of items in the slice. Implemented to satisfy
// sort.Interface.
func (al AcceptLanguages) Len() int { return len(al) }

// Swaps the items at position i and j. Implemented to satisfy sort.Interface.
func (al AcceptLanguages) Swap(i, j int) { al[i], al[j] = al[j], al[i] }

// Determines whether or not the item at position i is "less than" the item
// at position j. Implemented to satisfy sort.Interface.
func (al AcceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }

// Returns the parsed languages in a human readable fashion.
func (al AcceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range al {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(al)-1 {
			output.WriteString(", ")
		}
	}

	if output.Len() == 0 {
		output.WriteString("[]")
	}

	return output.String()
}

// The parsed structure is a slice of Accept-Language values stored in an
// AcceptLanguages instance, sorted based on the language qualifier.
func Read(lang string) AcceptLanguages {
	if lang == "" {
		return make(AcceptLanguages, 0)
	}
	acceptLanguageHeaderValues := strings.Split(lang, ",")
	acceptLanguages := make(AcceptLanguages, len(acceptLanguageHeaderValues))

	for i, languageRange := range acceptLanguageHeaderValues {
		// Check if a given range is qualified or not
		if qualifiedRange := strings.Split(languageRange, ";q="); len(qualifiedRange) == 2 {
			quality, error := strconv.ParseFloat(qualifiedRange[1], 32)
			if error != nil {
				// When the quality is unparseable, assume it's 1
				acceptLanguages[i] = AcceptLanguage{trimLanguage(qualifiedRange[0]), 1}
			} else {
				acceptLanguages[i] = AcceptLanguage{trimLanguage(qualifiedRange[0]), float32(quality)}
			}
		} else {
			acceptLanguages[i] = AcceptLanguage{trimLanguage(languageRange), 1}
		}
	}

	sort.Sort(acceptLanguages)
	return acceptLanguages
}

func ReadHeader(header http.Header) AcceptLanguages {
	lang := header.Get(acceptLanguageHeader)
	return Read(lang)
}

func trimLanguage(language string) string {
	return strings.Trim(language, " ")
}
