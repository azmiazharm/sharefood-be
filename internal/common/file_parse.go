package common

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"sharefood/internal/typex"
	"sharefood/pkg/util"
)

func FileToBuffer(file multipart.File) (*bytes.Buffer, error) {

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	return buf, nil
}

func FileFromMultipart(r *http.Request, field string, maxMemory int64) (*multipart.FileHeader, error) {
	mpr, err := r.MultipartReader()
	x, err := mpr.ReadForm(maxMemory)
	if err != nil {
		return nil, err
	}

	d := x.File[field]
	if len(d) == 0 {
		return nil, nil
	}

	return d[0], nil
}

// MultipartFormFile parse multipart file to own type file
func MultipartFormFile(r *http.Request, fieldName string, maxFileSize int64, extension []string) (*typex.File, error) {
	f, h, err := r.FormFile(fieldName)
	if err == http.ErrMissingFile {

		return nil, fmt.Errorf("the %s field required", fieldName)
	}

	if err != nil {
		return nil, fmt.Errorf("parse form file %s error %v", fieldName, err)

	}

	var buff bytes.Buffer
	size, err := buff.ReadFrom(f)
	if err != nil {
		return nil, fmt.Errorf("parse form file %s error %v", fieldName, err)
	}

	if size > maxFileSize {
		return nil, fmt.Errorf("the %s image file size %s is too large, max allow is %s", fieldName, HumanFileSize(float64(h.Size)), HumanFileSize(float64(maxFileSize)))
	}

	ct := ExtractFileExtension(buff.Bytes())
	ext := FileExtension(ct)

	if !ValidFileExtension(ext, extension) {
		return nil, fmt.Errorf("the %s image file extension .%s not allowed, request original file content type is %s, only allow: %s", fieldName, ext, ct, strings.Join(extension, ", "))
	}

	result := &typex.File{
		Filename:    fmt.Sprintf("%s.%s", util.GenerateReferenceID(fmt.Sprintf("%s-", fieldName)), ext),
		Buffer:      &buff,
		Size:        size,
		ContentType: ct,
		Ext:         ext,
	}

	return result, nil
}

func ValidFileExtension(ext string, extension []string) bool {
	return util.InArray(ext, extension)
}

func ExtractFileExtension(data []byte) string {
	return http.DetectContentType(data)
}

func FileExtension(input string) string {
	if len(input) < 1 {
		return input
	}

	return input[strings.IndexByte(input, '/')+1:]
}
