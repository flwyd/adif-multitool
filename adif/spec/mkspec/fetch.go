// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func httpGet(uri string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	// As of 2025-09-16, adif.org.uk rejects connections from certain user agents
	req.Header.Set("User-Agent", "ADIF Multitool mkspec")
	return client.Do(req)
}

func fetch(uri string) (filename string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	filename = path.Join("downloads", u.Hostname(), u.Path)
	if path.IsAbs(filename) {
		return "", fmt.Errorf("missing hostname from %q to %q", uri, filename)
	}
	if s, e := os.Stat(filename); (e == nil && s.Size() > 0) || !errors.Is(e, fs.ErrNotExist) {
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
	res, err := httpGet(uri)
	if err != nil {
		return
	}
	defer res.Body.Close()
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
