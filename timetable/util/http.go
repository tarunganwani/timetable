package util

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HttpPost(url string, payload []byte, headers map[string]string) (int, []byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("error creating post request %s", err.Error())
	}

	for hdrkey, hdrval := range headers {
		req.Header.Set(hdrkey, hdrval)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("post request error %s", err.Error())
	}
	defer resp.Body.Close()
	respbytes, err := io.ReadAll(resp.Body)
	log.Println("Status ", resp.StatusCode)
	log.Println(string(respbytes))
	return resp.StatusCode, respbytes, err

}
