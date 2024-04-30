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

func getFilePath() (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}
	pwd, err := os.Getwd()

	filePath := path.Join(*cacheDir, GetSHA1(&pwd))

	return filePath, nil
}

func CreateFile() (*os.File, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

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

func ReadFile() ([]byte, error) {
	filePath, err := getFilePath()
	if err != nil {
		return nil, err
	}

	if !isExists(filePath) {
		return nil, errors.New("record not exists for current pwd.")
	}

	f, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func DeleteFile() error {
	filePath, err := getFilePath()
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}
