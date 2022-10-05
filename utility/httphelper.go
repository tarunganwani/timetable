package utility

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func FireHttpServer(host, port string, router *mux.Router) error {

	addr := host + ":" + port
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error while firing http server", err)
		}
	}()

	log.Println("listening on ", (host + ":" + port))

	//register interrupt and kill signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server Started")

	<-done
	log.Println("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	return srv.Shutdown(ctx)
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

func BakeHeader(w http.ResponseWriter, status int) {
	w.WriteHeader(http.StatusNotFound)
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

func Bake200Response(w http.ResponseWriter, ttJsonData []byte) (err error) {

	BakeHeader(w, http.StatusOK)
	BakeJsonContentTypeInHeader(w)
	jsonResp, err := EncodeToJSON(GetMessageMap(string(ttJsonData)))
	if err != nil {
		return
	}
	_, err = w.Write(jsonResp)
	return
}
