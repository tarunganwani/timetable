package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func BakeHeader(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func BakeJsonContentTypeInHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func GetMessageMap(msg string) (m map[string]string) {
	m = make(map[string]string)
	m["message"] = msg
	return
}

func EncodeToJSON(v any) (jsonResp []byte, err error) {

	jsonResp, err = json.Marshal(v)
	if err != nil {
		errmsg := fmt.Sprintf("Error in JSON encoding. Err: %s", err.Error())
		err = errors.New(errmsg)
		return
	}
	return
}

func Bake400Response(w http.ResponseWriter, msg string) (err error) {

	BakeHeader(w, http.StatusBadRequest)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(msg))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}

func Bake404Response(w http.ResponseWriter, msg string) (err error) {

	BakeHeader(w, http.StatusNotFound)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(msg))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}

func Bake500Response(w http.ResponseWriter, msg string) (err error) {

	BakeHeader(w, http.StatusInternalServerError)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(msg))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}

func Bake501Response(w http.ResponseWriter, msg string) (err error) {

	BakeHeader(w, http.StatusNotImplemented)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(msg))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}

func Bake200Response(w http.ResponseWriter, msgBytes []byte) (err error) {

	BakeHeader(w, http.StatusOK)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(string(msgBytes)))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}

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
