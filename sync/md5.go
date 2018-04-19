package sync

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func GetFileMd5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", nil
	}

	buf := hash.Sum(nil)
	s := fmt.Sprintf("%032x", buf)

	return s, nil
}
