package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gobuffalo/packr"
)

const (
	GOOS = runtime.GOOS
)

func setup(assetsBox packr.Box) string {

	gshellDir, err := ioutil.TempDir("", "gshell")
	if err != nil {
		log.Fatal(err)
	}

	goZipBoxPath := fmt.Sprintf("%s/go.zip", GOOS)
	goZip, err := assetsBox.MustBytes(goZipBoxPath)
	if err != nil {
		log.Fatalf("static asset not found: %s", goZipBoxPath)
	}

	goZipPath := fmt.Sprintf("%s/go.zip", gshellDir)
	ioutil.WriteFile(goZipPath, goZip, 0644)
	_, err = unzip(goZipPath, gshellDir)
	if err != nil {
		log.Fatalf("Failed to unzip file %s -> %s", goZipPath, gshellDir)
	}

	return gshellDir
}

func unzip(src string, dest string) ([]string, error) {

	var filenames []string

	reader, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer reader.Close()

	for _, file := range reader.File {

		rc, err := file.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		fpath := filepath.Join(dest, file.Name)
		filenames = append(filenames, fpath)

		if file.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return filenames, err
			}
			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return filenames, err
			}
			_, err = io.Copy(outFile, rc)

			outFile.Close()

			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}
