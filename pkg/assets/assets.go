package assets

import (
	"bytes"
	"image"
	"os"
	"path/filepath"
)

func Root() string {
	return "assets"
}

func GetAsset(path string) (*[]byte, error) {
	// read file from assets
	// return file content
	assetFile, err := os.Open(filepath.Join(Root(), path))
	if err != nil {
		return nil, err
	}
	defer assetFile.Close()

	// get the file size
	fileInfo, _ := assetFile.Stat()

	// read the file content
	fileContent := make([]byte, fileInfo.Size())
	_, err = assetFile.Read(fileContent)

	if err != nil {
		return nil, err
	}

	return &fileContent, nil
}

func GetImage(path string) (*image.Image, error) {
	// read file from assets
	// decode image
	// return pointer to image
	fileContent, err := GetAsset(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(*fileContent))
	if err != nil {
		return nil, err
	}

	return &img, nil
}
