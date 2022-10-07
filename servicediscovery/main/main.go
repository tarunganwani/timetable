package main

import (
	"log"
	"servicediscovery/sdcontroller"
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

	//Get model
	model := sdmodel.NewSDModel()

	//Initialize controller
	controller, err := sdcontroller.NewSDController(model, sdcontroller.RouterConfig{
		Host:        sd_svc_host,
		Port:        sd_svc_port,
		Certificate: "",
	})
	if err != nil {
		log.Fatalln(err)
	}

	controller.Serve()
	log.Println("Done")
}
