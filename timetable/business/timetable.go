package business

import (
	"io"
	"log"
	"os"
)

func FetchResource(schoolcode, grade, division string) (jsonData []byte, found bool, err error) {

	dataDir := os.Getenv("DATA_DIR")
	jsonfile := dataDir + string(os.PathSeparator) + schoolcode + "_" + grade + "_" + division + ".json"
	log.Println("Finding resource ", jsonfile)
	fp, err := os.Open(jsonfile)
	if err != nil {
		log.Println(jsonfile, " not found")
		found = false
		return
	}

	found = true
	log.Println(jsonfile, " found, trying to read content.. ")
	jsonData, err = io.ReadAll(fp)
	return
}
