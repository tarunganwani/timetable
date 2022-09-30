package sdmodel

import "time"

type ServiceData struct {
	Name                  string
	Address               string
	Port                  string
	LastHeartbeatReceived time.Time
}

type ServiceMapType map[string][]ServiceData

var ServiceDiscoveryMap ServiceMapType
