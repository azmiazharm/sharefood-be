// Package util
package util

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"
)

const (
	letterBytes   = "0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	dateFormat     = `2006-01-02`
	dateTimeFormat = `2006-01-02 15:04:05`
)

var (
	// seed random
	seedRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// GenerateRandomNumberString generate random string number
func GenerateRandomNumberString(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

// GenerateReferenceID generates reference ID
func GenerateReferenceID(prefix string) string {
	now := time.Now().Format("20060102030405")
	buff := bytes.NewBufferString(now)
	buff.WriteString(GenerateRandomNumberString(8))

	return fmt.Sprintf("%s%s", prefix, buff.String())
}

// GenerateAppID generates reference ID
func GenerateAppID(prefix string) string {
	now := time.Now().Format("20060102030405")
	buff := bytes.NewBufferString(now)
	buff.WriteString(GenerateRandomNumberString(6))

	return fmt.Sprintf("%s%s", prefix, buff.String())
}

// GenerateRandomString generate random string
func GenerateRandomString(letterBytes string, n int) string {
	if n <= 0 {
		return ""
	}

	var letterRunes = []rune(letterBytes)
	b := make([]rune, n)

	if len(letterBytes) == 0 {
		return ""
	}

	for i := range b {
		b[i] = letterRunes[seedRand.Intn(len(letterRunes))]
	}

	return string(b)
}

func GenerateRandomBytesMask(letterBytes string, n int) string {
	b := make([]byte, n)
	l := len(letterBytes)
	for i, cache, remain := n-1, seedRand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = seedRand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < l {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
