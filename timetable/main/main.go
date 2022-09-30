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

	r := httphandlers.RouterConfig{
		Host: timetable_svc_host,
		Port: timetable_svc_port,
	}

	log.Println("Listening on ", (r.Host + ":" + r.Port))
	err = httphandlers.InitializeRouter(r)
	if err != nil {
		log.Fatalln(err)
	}
}
