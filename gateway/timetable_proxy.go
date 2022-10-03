package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/tarunganwani/timetable/utility"
)

const (
	timetableSvcname    = "timetable"
	serviceDiscoveryUrl = "http://localhost:8000/fetchaddress/"
	gateway_svc_host    = "127.0.0.1"
	gateway_svc_port    = "7000"
)

type ServiceData struct {
	Svcname string `json:"svcname"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

type ServiceDiscoveryResponse struct {
	Message string `json:"message"`
}

func fetchTimetableServiceAddress() (host string, port string, err error) {
	err = nil
	host = ""
	port = ""

	log.Println("fetching timetable service address")
	resp, err := http.Get(serviceDiscoveryUrl + timetableSvcname)
	if err != nil {
		err = fmt.Errorf("error fetching address from service discovery [%s]", err.Error())
		log.Println(err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("error fetching address from service discovery [Response code recd %d]", resp.StatusCode)
		log.Println(err.Error())
		return
	}
	respPayload, _ := io.ReadAll(resp.Body)
	var serviceDiscoveryResponse ServiceDiscoveryResponse
	err = json.Unmarshal(respPayload, &serviceDiscoveryResponse)
	if err != nil {
		err = fmt.Errorf("error decoding service discovery response [%s]", err.Error())
		log.Println(err.Error())
		return
	}
	var serviceData ServiceData
	err = json.Unmarshal([]byte(serviceDiscoveryResponse.Message), &serviceData)
	if err != nil {
		err = fmt.Errorf("error decoding service data [%s]", err.Error())
		log.Println(err.Error())
		return
	}
	// log.Printf("Response object %#v", serviceDiscoveryResponse)
	host = serviceData.Host
	port = serviceData.Port
	return
}

func proxyhandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = mux.Vars(r)["endpoint"]
		log.Println("Forwarding request to timetable service - endpoint ", r.URL.Path)
		p.ServeHTTP(w, r)
		log.Println("request served")
	}
}

func main() {

	//Initialize logger
	utility.InitializeLogger("gateway_service.log")

	host, port, err := fetchTimetableServiceAddress()
	if err != nil {
		log.Println(err)
	}

	timetableTargetUrl := "http://" + host + ":" + port
	remoteUrl, err := url.Parse(timetableTargetUrl)
	if err != nil {
		log.Println("Error parsing timetable target url ", err.Error())
	}
	ttProxy := httputil.NewSingleHostReverseProxy(remoteUrl)

	router := mux.NewRouter()
	router.HandleFunc("/{endpoint:.*}", proxyhandler(ttProxy)).Methods("GET")
	log.Printf("Listening on port %s ...\n", gateway_svc_host+":"+gateway_svc_port)
	http.Handle("/", router)

	err = utility.FireHttpServer(gateway_svc_host, gateway_svc_port, router)
	if err != nil {
		log.Fatalf("server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
