package common

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnzipAt(reader io.ReaderAt, size int64, dest string) error {
	rd, err := zip.NewReader(reader, size)
	if err != nil {
		return err
	}

	for _, file := range rd.File {
		fileRd, err := file.Open()
		if err != nil {
			return err
		}

		finalPath := filepath.Join(dest, file.Name)
		if !file.FileInfo().IsDir() {
			output, err := os.OpenFile(finalPath, os.O_WRONLY|os.O_CREATE, file.Mode())
			if err != nil {
				return err
			}
			defer output.Close()

			_, err = io.Copy(output, fileRd)
			if err != nil {
				return err
			}
		} else {
			err := os.MkdirAll(finalPath, file.Mode())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetFiles(reader io.ReaderAt, size int64) (map[string][]byte, error) { //map[path]content
	output := map[string][]byte{}
	rd, err := zip.NewReader(reader, size)
	if err != nil {
		return map[string][]byte{}, err
	}

	for _, file := range rd.File {
		fileRd, err := file.Open()
		if err != nil {
			return map[string][]byte{}, err
		}

		if !file.FileInfo().IsDir() {
			ba, err := io.ReadAll(fileRd)
			if err != nil {
				return map[string][]byte{}, err
			}

			output[file.Name] = ba
		}
	}
	return output, nil
}