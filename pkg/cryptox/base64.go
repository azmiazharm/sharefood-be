package cryptox

import "encoding/base64"

func Base64ToPlain(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ToBase64(plain string) string {
	return base64.URLEncoding.EncodeToString([]byte(plain))
}
