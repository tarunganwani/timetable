package utility

import (
	"io"
	"os"
)

func fn3() {

}

func ReadEntireFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}
