package helper

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//FileToBytes method gets the data from selected file and retrieve the byte value of it
func FileToBytes(filePath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//ObjectToFile method save object in selected path
func ObjectToFile(path string, object interface{}) error {
	if err := ioutil.WriteFile(path, convertToByte(object), 0666); err != nil {
		return err
	}

	return nil
}

func convertToByte(object interface{}) []byte {
	str, _ := json.Marshal(object)
	return str
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fPath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fPath)
		}

		filenames = append(filenames, fPath)

		if f.FileInfo().IsDir() {
			// Make Folder
			if err := os.MkdirAll(fPath, f.Mode()); err != nil {
				return filenames, err
			}

			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fPath), f.Mode()); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func Zip(src string, dest string) error {
	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	zipDestFile := zip.NewWriter(destinationFile)
	if err := filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filePath == dest {
			return nil
		}

		if err != nil {
			return err
		}

		relPath := strings.TrimPrefix(filePath, src)
		zipFile, err := zipDestFile.Create(relPath)

		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if err := zipDestFile.Close(); err != nil {
		return err
	}

	return nil
}