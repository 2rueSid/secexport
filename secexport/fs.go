package secexport

import (
	"errors"
	"os"
	"path"
)

func isExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}

	return true
}

func cacheDir() (*string, error) {
	p, e := os.UserCacheDir()
	if e != nil {
		return nil, e
	}

	cacheDir := path.Join(p, "secexport")
	if !isExists(cacheDir) {
		err := os.MkdirAll(cacheDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	return &cacheDir, nil
}

func CreateFile() (*os.File, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return nil, err
	}
	pwd, err := os.Getwd()

	filePath := path.Join(*cacheDir, GetSHA1(&pwd))

	isExists := isExists(filePath)
	if isExists {
		return nil, errors.New("data for that directory aready exists")
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Writes file to the given destination.
func WriteFile(f *os.File, d []byte) error {
	_, err := f.Write(d)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	return nil
}
