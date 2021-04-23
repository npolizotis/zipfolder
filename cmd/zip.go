package cmd

import (
	//"archive/zip"
	"github.com/klauspost/compress/zip"
	"io"
	"strings"
	"path/filepath"
	"os"
)

func RecursiveZip(pathToZip, zipFileName string) error {
	destinationFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, fileInfo os.FileInfo, err error) error {
		if fileInfo.IsDir() {
			return nil
		}
		basename:=fileInfo.Name()
		if basename[0]=='.' {
			return nil
		}

		if err != nil {
			return err
		}
		// remove leading "/"
		relPath := strings.TrimPrefix(filePath, pathToZip)
		if len(relPath)>0 && relPath[0]=='/' {
			relPath=relPath[1:]
		}
		//zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		

		header, err := zip.FileInfoHeader(fileInfo)
		header.Method=zip.Deflate
		header.Name = relPath
		zipFile, err := myZip.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
