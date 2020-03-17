package utils

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"time"

	dry "github.com/ungerik/go-dry"
)

func FileGetBase64(filenameOrURL string, timeout ...time.Duration) (out string, err error) {
	bs, err1 := dry.FileGetBytes(filenameOrURL, timeout...)
	if err1 != nil {
		err = err1
		return
	}
	out = base64.StdEncoding.EncodeToString(bs)
	return
}

func BinPath() string {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return p
}
