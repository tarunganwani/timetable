package main

import (
	"log"
	"servicediscovery/sdhttphandlers"
	"servicediscovery/sdmodel"

	"github.com/tarunganwani/timetable/utility"
)

const (
	sd_svc_host = "127.0.0.1"
	sd_svc_port = "8000"
)

func main() {

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in discovery service", err)
		}
	}()
	//Initialize logger
	utility.InitializeLogger("service_discovery.log")

	//Initialize model
	sdmodel.ServiceDiscoveryMap = make(sdmodel.ServiceMapType)

	//Initialize handlers
	log.Println("Listening on ", (sd_svc_host + ":" + sd_svc_port))
	router, err := sdhttphandlers.InitializeRouter(sdhttphandlers.RouterConfig{
		Host: sd_svc_host,
		Port: sd_svc_port,
	})
	if err != nil {
		log.Fatalln("error initializing http router", err)
	}

	err = utility.FireHttpServer(sd_svc_host, sd_svc_port, router)
	if err != nil {
		log.Fatalf("server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
