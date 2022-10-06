package sdhttphandlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"servicediscovery/sdbusiness"

	"github.com/gorilla/mux"
)

type ServiceData struct {
	Svcname string `json:"svcname"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func register_service_discovery_handlers(router *mux.Router) error {

	log.Println("Registering service discovery handlers")
	if router == nil {
		msg := "failed to register service discovery handlers, router not instantiated"
		log.Println(msg)
		return errors.New(msg)
	}
	router.HandleFunc("/heartbeat", HeartbeatServiceHandler).Methods("POST")
	router.HandleFunc("/deregister", DeregisterServiceHandler).Methods("POST")
	router.HandleFunc("/fetchaddress/{servicename}", FetchServiceAddressHandler).Methods("GET")
	log.Println("Handlers registered")
	return nil
}

func getRequestData(r *http.Request) (ServiceData, error) {
	decoder := json.NewDecoder(r.Body)
	var data ServiceData
	err := decoder.Decode(&data)
	return data, err
}

func HeartbeatServiceHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()
	log.Println("Got heartbeat")
	data, err := getRequestData(r)
	msg := ""
	if err != nil {
		msg = "Invalid input request"
		log.Println(msg, err.Error())
		Bake400Response(w, msg)
		return
	}

	log.Printf("Service details [svc %s host %s port %s] ", data.Svcname, data.Host, data.Port)
	err = sdbusiness.RegisterService(data.Svcname, data.Host, data.Port)
	if err != nil {
		msg = "Error registering service"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	msg = "ok"
	log.Println(msg)
	Bake200Response(w, []byte(msg))
}

func DeregisterServiceHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()
	log.Println("Handle deregister request")
	data, err := getRequestData(r)
	msg := ""
	if err != nil {
		msg = "Invalid input request"
		log.Println(msg, err.Error())
		Bake400Response(w, msg)
		return
	}

	err = sdbusiness.DeregisterService(data.Svcname, data.Host, data.Port)
	if err != nil {
		msg = "Error deregistering service - probablyService not registered"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	msg = "ok"
	log.Println(msg)
	Bake200Response(w, []byte(msg))
}

func FetchServiceAddressHandler(w http.ResponseWriter, r *http.Request) {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in proxy handler", err)
			Bake500Response(w, "internal error")
		}
	}()
	log.Println("Handle fetch-service request")
	reqvars := mux.Vars(r)
	svcModelData, err := sdbusiness.FetchServiceAddress(reqvars["servicename"])
	svcdata := ServiceData{
		Svcname: svcModelData.Name,
		Host:    svcModelData.Address,
		Port:    svcModelData.Port,
	}
	msg := ""
	if err != nil {
		msg = "Service not found"
		log.Println(msg, err.Error())
		Bake404Response(w, msg)
		return
	}

	svcdataBytes, err := json.Marshal(svcdata)
	if err != nil {
		msg = "Error processing request"
		log.Println(msg, err.Error())
		Bake500Response(w, msg)
		return
	}

	log.Println(string(svcdataBytes))
	Bake200Response(w, []byte(svcdataBytes))
}
