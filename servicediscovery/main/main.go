package main

import (
	"log"
	"servicediscovery/sdhttphandlers"
	"servicediscovery/sdmodel"
	"servicediscovery/util"
)

func main() {

	util.InitializeLogger("service_discovery.log")
	r := sdhttphandlers.RouterConfig{
		Host: "localhost",
		Port: "8000",
	}

	//Initialize model
	sdmodel.ServiceDiscoveryMap = make(sdmodel.ServiceMapType)

	log.Println("Listening on ", (r.Host + ":" + r.Port))
	err := sdhttphandlers.InitializeRouter(r)
	if err != nil {
		log.Fatalln(err)
	}
}
