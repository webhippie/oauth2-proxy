// generaTed by fileb0x at "2017-05-05 01:22:11.173267581 +0200 CEST" from config file "ab0x.yaml"

package templates

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
)

var (
	// CTX is a context for webdav vfs
	CTX = context.Background()

	// FS is a virtual memory file system
	FS = webdav.NewMemFS()

	// Handler is used to server files through a http handler
	Handler *webdav.Handler

	// HTTP is the http file system
	HTTP http.FileSystem = new(HTTPFS)
)

// HTTPFS implements http.FileSystem
type HTTPFS struct{}

// FileIndexHTML is "./index.html"
var FileIndexHTML = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x7c\xd0\x31\x4f\xc5\x20\x10\x07\xf0\xbd\x9f\xe2\xbc\xdd\x77\x79\xce\x94\xc4\xa8\xf3\xeb\xe0\xe2\x88\x40\xc3\x29\x16\xc2\x61\x63\xbf\xbd\xa1\xd4\xc5\xa1\x13\x47\xb8\xff\x3f\xbf\xa0\xee\x9e\x6f\x4f\xaf\x6f\xd3\x0b\x84\xfa\x15\xf5\x30\xa8\x7e\x02\xa8\xe0\x8d\x6b\x03\x80\x8a\xbc\x7c\x42\xf1\x71\x44\xb6\x69\x41\x08\xc5\xcf\x23\x52\x32\xdf\x35\x3c\xdc\xe7\x92\x7e\x36\x32\x22\xbe\x0a\xcd\x66\x6d\x3b\x17\xb6\x09\xf5\xf0\x3f\x2e\x75\x8b\x5e\x82\xf7\xf5\xb4\x64\xbf\x5c\xac\x08\x02\xed\x16\xea\x98\x36\xbe\x27\xb7\x1d\x2c\xc7\x2b\xb0\x1b\xd1\xe4\x1c\xd9\x9a\xca\x69\xc1\xfe\xd4\xf8\x57\x7d\x7b\x6c\xd5\x30\xb5\x36\x45\xe1\x7a\xc4\xc8\xf1\xfa\x47\x13\x5b\x38\x57\x90\x62\x4f\x29\x1f\x82\x5a\x51\x5f\xee\xa0\xce\x50\xd4\xbf\xeb\x37\x00\x00\xff\xff\x54\x04\xb5\x2d\x47\x01\x00\x00")

func init() {
	if CTX.Err() != nil {
		log.Fatal(CTX.Err())
	}

	var err error

	var f webdav.File

	var rb *bytes.Reader
	var r *gzip.Reader

	rb = bytes.NewReader(FileIndexHTML)
	r, err = gzip.NewReader(rb)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Close()
	if err != nil {
		log.Fatal(err)
	}

	f, err = FS.OpenFile(CTX, "./index.html", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	Handler = &webdav.Handler{
		FileSystem: FS,
		LockSystem: webdav.NewMemLS(),
	}
}

// Open a file
func (hfs *HTTPFS) Open(path string) (http.File, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	f, err := FS.OpenFile(CTX, path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0, bytes.MinRead))

	// If the buffer overflows, we will get bytes.ErrTooLarge.
	// Return that as an error. Any other panic remains.
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	_, err = buf.ReadFrom(f)
	return buf.Bytes(), err
}

// WriteFile is adapTed from ioutil
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	f, err := FS.OpenFile(CTX, filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

// FileNames is a list of files included in this filebox
var FileNames = []string{
	"./index.html",
}
