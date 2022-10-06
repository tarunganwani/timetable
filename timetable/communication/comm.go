package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tarunganwani/timetable/utility"
)

const (
	service_discovery_url         = "http://localhost:8000"
	sd_heartbeat_endpoint         = "heartbeat"
	sd_deregister_endpoint        = "deregister"
	sd_keep_alive_interval_in_sec = 60
	self_service_name             = "timetable"
)

type SDRegisterPayload struct {
	Svcname string `json:"svcname"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}

func BakeServiceDiscoveryRequest(svc_host, svc_port string) ([]byte, error) {

	data := SDRegisterPayload{
		Svcname: self_service_name,
		Host:    svc_host,
		Port:    svc_port,
	}
	return json.Marshal(data)
}

func RegisterAndKeepAliveWithServiceDiscovery(svc_host, svc_port string) error {

	heartbeat_req_payload, err := BakeServiceDiscoveryRequest(svc_host, svc_port)
	if err != nil {
		errmsg := fmt.Sprintf("Error creating register request for service discovery %s", err.Error())
		log.Println(errmsg)
		return errors.New(errmsg)
	}

	go func(svc_host, svc_port string, heartbeat_req_payload []byte) {

		sd_heartbeat_url := service_discovery_url + "/" + sd_heartbeat_endpoint
		request_headers := make(map[string]string)
		request_headers["Content-Type"] = "application/json"
		for {
			log.Println("Ping(heart-beat) service discovery @ ", sd_heartbeat_url)
			status, _, err := utility.HttpPost(sd_heartbeat_url, heartbeat_req_payload, request_headers)
			if err != nil {
				log.Println("Error POST-ing register request to service discovery", err.Error())
			}
			if err == nil && status != http.StatusOK {
				log.Println("Error POST-ing register request to service discovery. Response status received: ", status)
			}
			time.Sleep(sd_keep_alive_interval_in_sec * time.Second)
		}
	}(svc_host, svc_port, heartbeat_req_payload)

	return nil
}

func DeregisterFromServiceDiscovery(svc_host, svc_port string) error {

	deregister_req_payload, err := BakeServiceDiscoveryRequest(svc_host, svc_port)
	if err != nil {
		errmsg := fmt.Sprintf("Error creating deregister request for service discovery %s", err.Error())
		log.Println(errmsg)
		return errors.New(errmsg)
	}

	sd_deregister_url := service_discovery_url + "/" + sd_deregister_endpoint
	request_headers := make(map[string]string)
	request_headers["Content-Type"] = "application/json"
	log.Println("Deregister from service discovery @ ", sd_deregister_url)
	status, _, err := utility.HttpPost(sd_deregister_url, deregister_req_payload, request_headers)
	if err != nil {
		return fmt.Errorf("error POST-ing deregister request to service discovery %s", err.Error())
	}
	if err == nil && status != http.StatusOK {
		return fmt.Errorf("error POST-ing deregister request to service discovery. Response status received: %d", status)
	}
	return nil
}
