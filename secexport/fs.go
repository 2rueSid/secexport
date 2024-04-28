package secexport

import (
	"errors"
	"log"
	"os"
	"path"
)

func cacheDir() (*string, error) {
	p, e := os.UserCacheDir()
	if e != nil {
		return nil, e
	}

	cacheDir := path.Join(p, "secexport")
	if _, err := os.Stat(cacheDir); errors.Is(err, os.ErrNotExist) {
		log.Printf("Im here")
		err := os.MkdirAll(cacheDir, os.ModePerm)
		log.Printf("The err is: %v", err)
		if err != nil {
			return nil, err
		}
	}
	return &cacheDir, nil
}

func createFile(filePath string) (os.File, error) {
	f, err := os.Create(filePath)
	if err != nil {
		return os.File{}, err
	}

	return *f, nil
}

// Writes file to the given destination.
func WriteFile(hash string, d []byte) error {
	dirPath, err := cacheDir()
	if err != nil {
		return err
	}

	filePath := path.Join(*dirPath, hash)
	file, err := createFile(filePath)
	if err != nil {
		return err
	}

	_, err = file.Write(d)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}
