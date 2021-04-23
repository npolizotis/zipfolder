package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

func RecursiveZip(pathToZip, zipFileName string) error {
	destinationFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	myZip := zip.NewWriter(destinationFile)
	defer myZip.Close()

	err = filepath.Walk(pathToZip, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip zipping directories and dotfiles. Files within directories will be included
		if fileInfo.IsDir() || strings.HasPrefix(fileInfo.Name(), ".") {
			return nil
		}

		// remove leading "/"
		relPath := strings.TrimPrefix(filePath, pathToZip)
		if len(relPath) > 0 && relPath[0] == '/' {
			relPath = relPath[1:]
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fsFile.Close()

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Method = zip.Deflate
		header.Name = relPath
		zipFile, err := myZip.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		return err
	})
	if err != nil {
		return err
	}
	return myZip.Close()
}
