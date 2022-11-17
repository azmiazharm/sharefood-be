package common

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"sharefood/internal/typex"
	"sharefood/pkg/cryptox"
)

func ImageBase64Replacer(input string) string {
	if len(input) < 1 {
		return input
	}
	return input[strings.IndexByte(input, ',')+1:]
}

func DecodeBaseImage64(fileBase64, fileName string) (*typex.File, error) {
	dBytes, err := cryptox.Base64ToPlain(ImageBase64Replacer(fileBase64))
	if err != nil {
		return nil, err
	}

	mType := http.DetectContentType(dBytes)
	ext := FileExtension(mType)

	fl := &typex.File{
		Filename:    fmt.Sprintf("%s.%s", FileNameWithoutExtension(fileName), ext),
		ContentType: mType,
		Size:        int64(len(dBytes)),
		Ext:         ext,
		Buffer:      bytes.NewBuffer(dBytes),
	}

	return fl, nil
}

func DetectImageBase6ContentType(fileBase64 string) (string, error) {
	dBytes, err := cryptox.Base64ToPlain(ImageBase64Replacer(fileBase64))
	if err != nil {
		return "", err
	}

	return http.DetectContentType(dBytes), nil
}

func FileNameWithoutExtension(nm string) string {
	if len(nm) < 1 {
		return nm
	}

	i := strings.IndexByte(nm, '.')
	if i > 0 {
		return nm[:i]
	}

	return nm
}
