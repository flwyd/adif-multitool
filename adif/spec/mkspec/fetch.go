package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

func fetch(uri string) (filename string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	filename = path.Join("downloads", u.Hostname(), u.Path)
	if path.IsAbs(filename) {
		return "", fmt.Errorf("missing hostname from %q to %q", uri, filename)
	}
	if _, e := os.Stat(filename); e == nil || !errors.Is(e, fs.ErrNotExist) {
		log.Printf("Using cached %s", filename)
		return
	}
	dir := path.Dir(filename)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return
	}
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	log.Printf("Downloading %s to %s", uri, filename)
	res, err := http.Get(uri)
	err = res.Write(f)
	return
}

func xmlFromZip(filename string) ([]byte, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	for _, f := range r.File {
		if path.Base(f.Name) == "all.xml" && path.Base(path.Dir(f.Name)) == "xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()
			var xml bytes.Buffer
			_, err = io.Copy(&xml, rc)
			if err != nil {
				return nil, err
			}
			return xml.Bytes(), nil
		}
	}
	return nil, nil
}
