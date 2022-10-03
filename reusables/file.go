package reusables

import (
	"io"
	"os"
)

func XYZ() {

}
func ReadEntireFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(file)
}
