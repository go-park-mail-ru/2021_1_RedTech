package fileutils

import (
	"Redioteka/internal/pkg/utils/randstring"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func createFile(root, dir, name string) (*os.File, error) {
	_, err := os.ReadDir(root + dir)
	if err != nil {
		err = os.MkdirAll(root+dir, 0777)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(root + dir + name)
	return file, err
}

func UploadFile(reader io.Reader, root, path, urlRoot, ext string) (string, error) {
	filename := randstring.RandString(32) + ext
	log.Print("Created file with name ", filename)
	file, err := createFile(root, path, filename)
	if err != nil {
		return "", fmt.Errorf("file createing error %s", err)
	}
	defer file.Close()

	filename = urlRoot + path + filename
	_, err = io.Copy(file, reader)
	if err != nil {
		return "", fmt.Errorf("copy error: %s", err)
	}
	return filename, nil
}

func UploadFromRequest(r *http.Request, root, path, urlRoot string) (string, error) {
	r.ParseMultipartForm(5 * 1024 * 1024)
	uploaded, header, err := r.FormFile("avatar")
	if err != nil {
		return "", fmt.Errorf("error while parsing file: %s", err)
	}
	defer uploaded.Close()
	return UploadFile(uploaded, root, path, urlRoot, filepath.Ext(header.Filename))
}
