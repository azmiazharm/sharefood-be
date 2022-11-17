package typex

import "bytes"

// File struct for file information and storage location
type File struct {
	Filename    string        `json:"-"`
	Ext         string        `json:"-"`
	ContentType string        `json:"-"`
	Size        int64         `json:"-"`
	Buffer      *bytes.Buffer `json:"-"`
}
