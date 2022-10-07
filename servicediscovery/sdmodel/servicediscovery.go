package sdmodel

import (
	"log"
	"time"
)

type ServiceData struct {
	Name                  string
	Address               string
	Port                  string
	LastHeartbeatReceived time.Time
}

type ServiceMapType map[string][]ServiceData

type Model struct {
	ServiceDiscoveryMap ServiceMapType
	Initialized         bool
}

func (md *Model) initModel() {
	md.ServiceDiscoveryMap = make(ServiceMapType)
	md.Initialized = true
}

func NewSDModel() *Model {
	log.Println("Initializing model..")
	model := new(Model)
	model.initModel()
	return model
}
