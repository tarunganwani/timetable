package util

import (
	"io"
	"os"
)

func ABC() {

}
func ReadEntireFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}
