// Package util
package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
	toCamelRegs   = map[string]*regexp.Regexp{
		" ": regexp.MustCompile(" +[a-zA-Z]"),
		"-": regexp.MustCompile("-+[a-zA-Z]"),
		"_": regexp.MustCompile("_+[a-zA-Z]"),
	}
)

func StringJoin(elems []string, sep, lastSep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%s%s", elems[0], lastSep)
	}
	n := len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(elems[0])
	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}

	if lastSep != "" {
		b.WriteString(lastSep)
	}

	return b.String()
}

// StringContains check contain string
func StringContains(s string, contains []string) bool {
	for i := 0; i < len(contains); i++ {
		if strings.Contains(strings.ToLower(s), strings.ToLower(contains[i])) {
			return true
		}
	}
	return false
}

// SubString substitute string
func SubString(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

// Deduplicate returns a new string
// with duplicates values removed.
func Deduplicate(s []rune) string {
	if len(s) <= 1 {
		return string(s)
	}

	result := []rune{}
	seen := make(map[rune]struct{})
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = struct{}{}
		}
	}
	return string(result)
}

func Contains(s string, haystack []string) bool {
	for i := 0; i < len(haystack); i++ {
		if strings.Contains(s, haystack[i]) {
			return true
		}
	}

	return false
}

func Abbreviate(s string, n int, next bool) string {
	var (
		rgx   = regexp.MustCompile(`[^a-zA-Z]+`)
		rgx2   = regexp.MustCompile(`[^a-zA-Z\s]+`)
		words []string
	)

	s = rgx2.ReplaceAllString(s, "")

	if len(s) == n {
		return strings.ToUpper(s)
	}

	words = strings.Fields(s)

	nWords := []string{}

	var result string

	for i := 0; i < len(words); i++ {

		if len(result) == n {
			break
		}

		if Contains(words[i], []string{".", ",", "'"}) {
			continue
		}

		w := rgx.ReplaceAllString(words[i], "")
		if len(w) <= 1 {
			continue
		}

		result += strings.Title(string([]rune(w)[0]))
		nWords = append(nWords, w[1:])
	}

	r := 0

	if len(result) < n {
		r = n - len(result)
	}

	if r > 0 && !next && len(nWords) > 0 {
		return addCharPosition(result, nWords, r)
	}

	if next {
		front := SubStringLeft(result, 1)
		result = front + strings.ToUpper(GenerateRandomString(Deduplicate([]rune(strings.Join(nWords, ""))), n-1))

	}

	return result
}

func addCharPosition(rc string, hs []string, n int) string {
	rn := rand.New(rand.NewSource(time.Now().Unix())) // initialize local pseudorandom generator
	idx := rn.Intn(len(hs))
	if idx == (len(hs) - 1) {
		rc += strings.ToUpper(GenerateRandomString(Deduplicate([]rune(hs[idx])), n))
		return rc
	}

	rc = fmt.Sprintf(
		"%s%s%s",
		SubStringLeft(rc, 1),
		SubStringRight(rc, 1),
		strings.ToUpper(GenerateRandomString(Deduplicate([]rune(hs[idx])), n)),
	)

	return rc
}

// SubStringRight substitute string from right
func SubStringRight(input string, length int) string {
	r := []rune(input)

	if length <= 0 {
		return input
	}

	if len(r) <= length {
		return input
	}

	return string(r[len(r)-length:])
}

// SubStringLeft substitute string from left
func SubStringLeft(input string, length int) string {
	r := []rune(input)

	if length <= 0 {
		return input
	}

	if len(r) <= length {
		return input
	}
	return string(r[:length])
}


func EmailDomain(input string) string {
	if len(input) < 1 {
		return input
	}
	return input[strings.IndexByte(input, '@')+1:]
}

func ReplaceDoubleSpace( s string) string  {
	s = strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	re := regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")
	return s
}

// CamelToSnakeCase converts a CamelCase string to a snake_case string.
func CamelToSnakeCase(s string) string {
	return strings.ToLower(regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}"))
}

// ToSnakeCase converts a CamelCase string to a snake_case string.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase convert string to camel case.
//	"range_price" -> "rangePrice"
//	"range price" -> "rangePrice"
//	"range-price" -> "rangePrice"
func ToCamelCase(s string, sep ...string) string {
	sepChar := "_"
	if len(sep) > 0 {
		sepChar = sep[0]
	}

	// Not contains sep char
	if !strings.Contains(s, sepChar) {
		return s
	}

	// Get regexp instance
	rgx, ok := toCamelRegs[sepChar]
	if !ok {
		rgx = regexp.MustCompile(regexp.QuoteMeta(sepChar) + "+[a-zA-Z]")
	}

	return rgx.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.TrimLeft(s, sepChar)
		return UpperFirst(s)
	})
}

func UpperFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	rs := []rune(s)
	f := rs[0]

	if 'a' <= f && f <= 'z' {
		return string(unicode.ToUpper(f)) + string(rs[1:])
	}
	return s
}
