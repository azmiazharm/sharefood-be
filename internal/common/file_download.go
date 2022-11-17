package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"sharefood/internal/typex"
)

func HTTPDownload(uri string) ([]byte, error) {
	res, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return d, err
}

func DownloadInMemory(uri string) (*typex.File, error) {
	res, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed download, response status is %v", res.StatusCode)
	}

	var (
		r bytes.Buffer
	)

	defer res.Body.Close()

	_, err = io.Copy(&r, res.Body)
	if err != nil {
		return nil, err
	}

	mime := ExtractFileExtension(r.Bytes())
	result := &typex.File{
		ContentType: mime,
		Ext:         FileExtension(mime),
		Size:        int64(len(r.Bytes())),
	}

	return result, err
}

func WriteFile(dst string, d []byte) error {
	return ioutil.WriteFile(dst, d, 0444)
}

func DownloadToFile(uri string, dst string) error {
	b, err := HTTPDownload(uri)
	if err != nil {
		return err
	}

	return WriteFile(dst, b)
}
