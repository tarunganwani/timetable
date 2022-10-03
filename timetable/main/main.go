package main

import (
	"log"
	"timetable/communication"
	"timetable/httphandlers"
	"timetable/util"
)

const (
	timetable_svc_host = "127.0.0.1"
	timetable_svc_port = "8080"
)

func main() {

	util.InitializeLogger()

	//Do nothing with the errors, since these are already logged
	err := communication.RegisterAndKeepAliveWithServiceDiscovery(timetable_svc_host, timetable_svc_port)
	if err != nil {
		log.Fatalln("Error initializing communnication with service discovery")
	}

	router, err := httphandlers.InitializeRouter(httphandlers.RouterConfig{
		Host: timetable_svc_host,
		Port: timetable_svc_port,
	})
	if err != nil {
		log.Fatalln("error initializing http router")
	}

	err = util.FireHttpServer(timetable_svc_host, timetable_svc_port, router)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}
