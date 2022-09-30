package communication

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"timetable/util"
)

const (
	service_discovery_url         = "http://localhost:8000"
	sd_heartbeat_endpoint         = "heartbeat"
	sd_deregister_endpoint        = "deregister"
	sd_keep_alive_interval_in_sec = 5
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

	// TODO implement exit criteria for the below go routine and the microservice as a whole
	// TODO Deregister from service discovery at service exit

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
			status, _, err := util.HttpPost(sd_heartbeat_url, heartbeat_req_payload, request_headers)
			if err != nil {
				log.Println("Error POST-ing register request to service discovery", err.Error())
			}
			if status != http.StatusOK {
				log.Println("Error POST-ing register request to service discovery. Response status received: ", status)
			}
			time.Sleep(sd_keep_alive_interval_in_sec * time.Second)
		}
	}(svc_host, svc_port, heartbeat_req_payload)

	return nil
}
