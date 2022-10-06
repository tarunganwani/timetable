package main

import (
	"log"
	"timetable/communication"
	"timetable/httphandlers"

	"github.com/tarunganwani/timetable/utility"
)

const (
	timetable_svc_host = "127.0.0.1"
	timetable_svc_port = "8080"
)

//TODO IMPROVEMENT - use messaging queues instead of http communication within among the internal services

//TODO MAJOR - write concurrent tests
//TODO MAJOR - run go tools to static analyze race conditions
//TODO MAJOR - fix race condition in model map ion timetable

//TODO MAJOR - write tests for http handlers, models and other packages for all modules
//TODO MAJOR - write benchmarks

//TODO MAJOR - implement health check between gateway and service discovery

//TODO MAJOR - dockerize all apps(gateway, service discovery and timetable service)
//TODO MAJOR - container respawn logic with K8S or other relevant tech
//TODO MAJOR - automate build with makefile or a single go command
//TODO MAJOR - set a CI/CD pipeline
//TODO MAJOR - automate deployment
//TODO MAJOR - refactor gateway service as a whole - divide into packages
//TODO MAJOR - implementation to prevent ddos attack in gateway
//DONE -- TODO MAJOR - recover from panic in important functions(main, handlers, etc) for graceful exits
//DONE -- TODO MAJOR - error handling and retries when one of the communicating service goes down

func main() {

	utility.InitializeLogger("timetable_srv.log")

	defer func() {
		err := recover()
		if err != nil {
			log.Println("error encountered in timetable service", err)
		}
	}()
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

	err = utility.FireHttpServer(timetable_svc_host, timetable_svc_port, router)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	err = communication.DeregisterFromServiceDiscovery(timetable_svc_host, timetable_svc_port)
	if err != nil {
		log.Fatalln(err)
	}

	log.Print("Server Exited Properly")

}
